package tsearch

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

type Storage interface {
	SaveTokens(id uint32, tokens []string) (err error)
	GetTokens(id uint32) (tokens []string, err error)
	DeleteTokens(id uint32) (err error)

	UpdateIndex(id uint32, oldTokens []string, newTokens []string) (err error)
	SearchIndex(tokens []string) (hits map[uint32]*HitCounter, err error)
}

type redisStorage struct {
	conn redis.Conn
}

func (s *redisStorage) getTokenKey(id uint32) string {
	return fmt.Sprintf("tsearch_token:%d", id)
}

func (s *redisStorage) SaveTokens(id uint32, tokens []string) (err error) {
	data, err := json.Marshal(tokens)
	if err != nil {
		return err
	}
	key := s.getTokenKey(id)
	_, err = redis.String(s.conn.Do("SET", key, data))
	return err
}

func (s *redisStorage) DeleteTokens(id uint32) (err error) {
	key := s.getTokenKey(id)
	_, err = s.conn.Do("DEL", key)
	return err
}

func (s *redisStorage) GetTokens(id uint32) (tokens []string, err error) {
	key := s.getTokenKey(id)
	data, err := redis.String(s.conn.Do("GET", key))
	if err == redis.ErrNil {
		return tokens, nil
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &tokens)
	return
}

func (s *redisStorage) indexKey(key string) string {
	return fmt.Sprintf("index:%s", key)
}

func (s *redisStorage) UpdateIndex(id uint32, oldTokens []string, newTokens []string) (err error) {
	if len(oldTokens) == 0 && len(newTokens) == 0 {
		return nil
	}

	_, err = s.conn.Do("MULTI")
	if err != nil {
		return err
	}
	for _, token := range oldTokens {
		key := s.indexKey(token)
		_ = s.conn.Send("SREM", key, id)
	}

	for _, token := range newTokens {
		key := s.indexKey(token)
		_ = s.conn.Send("SADD", key, id)
	}
	_, err = s.conn.Do("EXEC")
	return err
}

type HitCounter struct {
	Count int
}

func (s *redisStorage) SearchIndex(tokens []string) (hits map[uint32]*HitCounter, err error) {
	hits = make(map[uint32]*HitCounter)
	if len(tokens) == 0 {
		return
	}
	_, err = s.conn.Do("MULTI")
	if err != nil {
		return
	}

	for _, token := range tokens {
		key := s.indexKey(token)
		_ = s.conn.Send("SMEMBERS", key)
	}

	values, err := redis.Values(s.conn.Do("EXEC"))
	if err != nil {
		return
	}

	for _, value := range values {
		idList, err := redis.Strings(value, nil)
		if err != nil {
			return nil, err
		}

		for _, idStr := range idList {
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return nil, err
			}

			counter, ok := hits[uint32(id)]
			if ok {
				counter.Count += 1
			} else {
				counter = &HitCounter{
					Count: 1,
				}
				hits[uint32(id)] = counter
			}
		}
	}
	return
}

func NewRedisStorage(addr string) (Storage, error) {
	conn, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &redisStorage{
		conn: conn,
	}, nil

}
