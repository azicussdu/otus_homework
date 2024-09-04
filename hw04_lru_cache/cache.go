package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       sync.Mutex
}

type CacheItem struct {
	key   Key
	value interface{}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if foundListItem, ok := l.items[key]; ok {
		foundListItem.Value.(*CacheItem).value = value
		l.queue.MoveToFront(foundListItem)
		return true
	}

	cacheItem := &CacheItem{key: key, value: value}
	listItem := l.queue.PushFront(cacheItem)
	l.items[key] = listItem

	if l.queue.Len() > l.capacity {
		lastListItem := l.queue.Back()
		delete(l.items, lastListItem.Value.(*CacheItem).key)
		l.queue.Remove(lastListItem)
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if foundListItem, ok := l.items[key]; ok {
		l.queue.MoveToFront(foundListItem)
		return foundListItem.Value.(*CacheItem).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.items = map[Key]*ListItem{}
	l.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
