package hw04lrucache

import (
	"fmt"
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	sync.Mutex
	items map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.Lock()
	defer lc.Unlock()
	if _, ok := lc.items[key]; ok {
		lc.items[key].Value = value
		lc.queue.MoveToFront(lc.items[key])
		return true
	}
	lc.items[key] = lc.queue.PushFront(value)
	if lc.capacity < len(lc.items) {
		backListItem := lc.queue.Back()
		lc.queue.Remove(backListItem)
		delete(lc.items, Key(fmt.Sprintf("%v", backListItem.Value)))
	}
	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.Lock()
	defer lc.Unlock()
	if _, ok := lc.items[key]; !ok {
		return nil, false
	}
	lc.queue.MoveToFront(lc.items[key])
	return lc.items[key].Value, true
}

func (lc *lruCache) Clear() {
	*lc = lruCache{
		capacity: lc.capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, lc.capacity),
	}
}
