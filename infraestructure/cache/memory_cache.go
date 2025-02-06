package cache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type memoryCacheItem struct {
	key         string
	value       interface{}
	expireAfter time.Time
}

type MemoryCache struct {
	mostRecentlyRead *list.List
	elementsByKey    map[string]*list.Element
	mu               sync.Mutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		mostRecentlyRead: list.New(),
		elementsByKey:    make(map[string]*list.Element),
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elmt, ok := c.elementsByKey[key]
	if !ok {
		return nil, false, nil
	}

	item := elmt.Value.(memoryCacheItem)

	if time.Now().After(item.expireAfter) {
		c.Clear(ctx, key)
		return nil, false, nil
	}

	c.mostRecentlyRead.MoveToBack(elmt)
	return item.value, true, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, val interface{}, expireAfter time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := memoryCacheItem{
		key:         key,
		value:       val,
		expireAfter: time.Now().Add(expireAfter),
	}
	if len(c.elementsByKey) > 100 {
		oldest := c.mostRecentlyRead.Front()
		c.mostRecentlyRead.Remove(oldest)
		delete(c.elementsByKey, oldest.Value.(memoryCacheItem).key)
	}

	elmt := c.mostRecentlyRead.PushBack(item)
	c.elementsByKey[key] = elmt
	return nil
}

func (c *MemoryCache) Clear(ctx context.Context, key string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elmt, ok := c.elementsByKey[key]
	if !ok {
		return false, nil
	}

	c.mostRecentlyRead.Remove(elmt)
	delete(c.elementsByKey, key)
	return true, nil
}
