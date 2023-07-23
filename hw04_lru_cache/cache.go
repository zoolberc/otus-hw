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

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	element := &ListItem{Value: value, Key: key}
	if _, ok := lc.items[key]; !ok {
		if lc.queue.Len() == lc.capacity {
			lastElement := lc.queue.Back()
			lc.queue.Remove(lastElement)
			delete(lc.items, lastElement.Key)
		}
		lc.items[key] = element
		lc.queue.PushFront(element)
		return false
	}
	lc.items[key] = element
	lc.queue.MoveToFront(element)
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
