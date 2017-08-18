package memcache

import (
	"errors"
	"strings"
	"time"

	"github.com/adolphlxm/atc/cache"
	"github.com/bradfitz/gomemcache/memcache"
)

type Cache struct {
	conn     *memcache.Client
	conninfo []string
}

// NewMemCache create new memcache adapter.
func NewMemCache() cache.Cache {
	return &Cache{}
}

// Get get value from memcache.
func (c *Cache) Get(key string) ([]byte, error) {
	if c.conn == nil {
		if err := c.connectInit(); err != nil {
			return nil,err
		}
	}

	item, err := c.conn.Get(key)
	if err != nil  {
		return nil, err
	}

	return item.Value,nil
}

// Set set value to memcache. only support string.
func (c *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	if c.conn == nil {
		if err := c.connectInit(); err != nil {
			return err
		}
	}
	v, ok := val.(string)
	if !ok {
		return errors.New("val must string")
	}
	item := memcache.Item{Key: key, Value: []byte(v), Expiration: int32(timeout / time.Second)}
	return c.conn.Set(&item)
}

// Delete delete value in memcache.
func (c *Cache) Delete(key string) error {
	if c.conn == nil {
		if err := c.connectInit(); err != nil {
			return err
		}
	}
	return c.conn.Delete(key)
}

// FlushAll clear all cached in memcache.
func (c *Cache) FlushAll() error {
	if c.conn == nil {
		if err := c.connectInit(); err != nil {
			return err
		}
	}
	return c.conn.FlushAll()
}

// New start memcache adapter.
// if connecting error, return.
func (c *Cache) New(config string) error {
	c.conninfo = strings.Split(config, ";")
	if c.conn == nil {
		if err := c.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connectInit to memcache and keep the connection.
func (c *Cache) connectInit() error {
	c.conn = memcache.New(c.conninfo...)
	return nil
}


func init() {
	cache.Register("memcache", NewMemCache)
}