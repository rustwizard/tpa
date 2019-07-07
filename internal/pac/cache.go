package pac

import (
	"crypto/md5"
	"fmt"

	"github.com/gomodule/redigo/redis"

	"github.com/pkg/errors"

	"github.com/vmihailenco/msgpack"
)

var ErrNotFound = errors.New("cache: item not found")

type Cache struct {
	client redis.Conn
}

func NewCache(client redis.Conn) *Cache {
	return &Cache{client: client}
}

func (s *Cache) SetResponse(k string, r *Response) error {
	buf, err := msgpack.Marshal(r)
	if err != nil {
		return err
	}

	if _, err := s.client.Do("SET", key(k), string(buf)); err != nil {
		return err
	}

	return nil
}

func (s *Cache) GetResponse(k string) (*Response, error) {
	var resp Response
	reply, err := redis.String(s.client.Do("GET", key(k)))
	if err != nil && err != redis.ErrNil {
		return &resp, err
	}

	buf := []byte(reply)
	if len(buf) > 0 {
		if err := msgpack.Unmarshal(buf, &resp); err != nil {
			return &resp, err
		}
		return &resp, nil
	}

	return &resp, ErrNotFound
}

func key(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
