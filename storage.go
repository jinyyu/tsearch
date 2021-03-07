package tsearch

import "github.com/gomodule/redigo/redis"

type Storage interface {
	Speak() string
}

type redisStorage struct {
	conn redis.Conn
}

func (s *redisStorage) Speak() string {
	return ""
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
