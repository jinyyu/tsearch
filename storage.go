package tsearch

import (
	"github.com/gomodule/redigo/redis"
)

type KeyValue struct {
	Key   string
	Value string
}

type Storage interface {
	MultiSet(kvs ...*KeyValue) (err error)
	MultiGet(keys ...string) (values []string, err error)
	MultiDel(keys ...string) (err error)

	MultiSetAdd(kvs ...*KeyValue) (err error)
	MultiSetDel(kvs ...*KeyValue) (err error)
	MultiGetMembers(keys ...string) (reply []interface{}, err error)
}

type redisStorage struct {
	conn redis.Conn
}

func (s *redisStorage) MultiSet(kvs ...*KeyValue) (err error) {
	if len(kvs) == 0 {
		return
	}

	param := make([]interface{}, 0, 2*len(kvs))
	for _, kv := range kvs {
		param = append(param, kv.Key, kv.Value)
	}
	_, err = s.conn.Do("MSET", param...)
	return err
}

func (s *redisStorage) MultiGet(keys ...string) (values []string, err error) {
	if len(keys) == 0 {
		return
	}

	param := make([]interface{}, len(keys))
	for i := range keys {
		param[i] = keys[i]
	}
	values, err = redis.Strings(s.conn.Do("MGET", param...))
	if err != nil && err != redis.ErrNil {
		return values, err
	}
	return values, nil
}

func (s *redisStorage) MultiDel(keys ...string) (err error) {
	if len(keys) == 0 {
		return
	}
	param := make([]interface{}, len(keys))
	for i := range keys {
		param[i] = keys[i]
	}

	_, err = s.conn.Do("DEL", param...)
	return err
}

func (s *redisStorage) MultiSetAdd(kvs ...*KeyValue) (err error) {
	if len(kvs) == 0 {
		return
	}
	_, err = s.conn.Do("MULTI")
	for _, kv := range kvs {
		_ = s.conn.Send("SADD", kv.Key, kv.Value)
	}
	_, err = s.conn.Do("EXEC")
	return err
}

func (s *redisStorage) MultiSetDel(kvs ...*KeyValue) (err error) {
	if len(kvs) == 0 {
		return
	}
	_, err = s.conn.Do("MULTI")
	for _, kv := range kvs {
		_ = s.conn.Send("SREM", kv.Key, kv.Value)
	}
	_, err = s.conn.Do("EXEC")
	return err
}

func (s *redisStorage) MultiGetMembers(keys ...string) (reply []interface{}, err error) {
	if len(keys) == 0 {
		return nil, err
	}

	_, err = s.conn.Do("MULTI")
	for _, key := range keys {
		_ = s.conn.Send("SMEMBERS", key)
	}

	reply, err = redis.Values(s.conn.Do("EXEC"))
	if err != nil && err != redis.ErrNil {
		return
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
