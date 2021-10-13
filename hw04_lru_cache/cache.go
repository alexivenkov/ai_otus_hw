package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item := cacheItem{
		key:   key,
		value: value,
	}
	_, exists := c.items[key]

	if !exists {
		c.queue.PushFront(item)

		if len(c.items) > c.capacity {
			c.queue.Remove(c.queue.Back())
		}
	} else {
		c.queue.Remove(c.queue.Front())
		c.queue.PushFront(item)
	}

	c.items[key] = c.queue.Front()

	return exists
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	val, exists := c.items[key]

	if exists {
		c.queue.MoveToFront(val)

		return val.Value.(cacheItem).value, exists
	}

	return nil, false
}

func (c *lruCache) Clear() {
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
