package tsearch

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Storage interface {
	SaveTokens(id uint32, tokens []string) (err error)
	GetTokens(id uint32) (tokens []string, err error)
}

type redisStorage struct {
	conn redis.Conn
}

func (s *redisStorage) getTokenKey(id uint32) string {
	return fmt.Sprintf("tsearch_token:%d", id)
}

func (s *redisStorage) SaveTokens(id uint32, tokens []string) (err error) {
	return nil
}

func (s *redisStorage) GetTokens(id uint32) (tokens []string, err error) {
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
