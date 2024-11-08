package db

import (
	"github.com/gomodule/redigo/redis"
)

type Redis struct {
	conn     redis.Conn
	username string
	password string
	address  string
}

func NewRedis(username, password, address string) (*Redis, error) {
	conn, err := redis.Dial("tcp", address, redis.DialUsername(username), redis.DialPassword(password))
	if err != nil {
		return nil, err
	}

	return &Redis{
		conn:     conn,
		username: username,
		password: password,
		address:  address,
	}, nil
}

func (r *Redis) Close() error {
	return r.conn.Close()
}

func (r *Redis) Do(command string, args ...any) (any, error) {
	return r.conn.Do(command, args...)
}
