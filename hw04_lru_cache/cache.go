package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	defer c.Unlock()

	item := cacheItem{
		key:   key,
		value: value,
	}
	c.Lock()
	_, exists := c.items[key]

	if exists {
		c.queue.Remove(c.queue.Front())
	}

	if len(c.items) > c.capacity {
		c.queue.Remove(c.queue.Back())
	}

	c.queue.PushFront(item)

	c.items[key] = c.queue.Front()

	return exists
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	defer c.Unlock()
	c.Lock()
	val, exists := c.items[key]

	if exists {
		c.queue.MoveToFront(val)

		return val.Value.(cacheItem).value, exists
	}

	return nil, false
}

func (c *lruCache) Clear() {
	defer c.Unlock()
	c.Lock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
