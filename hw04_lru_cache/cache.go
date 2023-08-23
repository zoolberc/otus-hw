package hw04lrucache

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
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	if _, ok := lc.items[key]; !ok {
		if lc.queue.Len() == lc.capacity {
			lastElement := lc.queue.Back()
			delete(lc.items, lc.queue.Back().Key)
			lc.queue.Remove(lastElement)
		}
		lc.items[key] = lc.queue.PushFront(value)
		lc.items[key].Key = key
		return false
	}
	lc.items[key].Value = value
	lc.queue.MoveToFront(lc.items[key])
	return true
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	element, ok := lc.items[key]
	if !ok {
		return nil, false
	}
	lc.queue.MoveToFront(element)
	return element.Value, true
}

func (lc *lruCache) Clear() {
	lc.items = make(map[Key]*ListItem, lc.capacity)
	lc.queue = NewList()
}
