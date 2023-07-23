package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Key   Key
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	items map[*ListItem]struct{}
	head  *ListItem
	tail  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return len(l.items)
}

func (l *list) Front() *ListItem {
	if l.items == nil {
		return nil
	}
	return l.head
}

func (l *list) Back() *ListItem {
	if l.tail == nil {
		return nil
	}
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	element := &ListItem{Value: v}
	if l.items == nil {
		l.head = element
		l.tail = element
		l.items = make(map[*ListItem]struct{})
		l.items[element] = struct{}{}
		return element
	}
	element.Next = l.head
	l.head.Prev = element
	l.head = element
	l.items[element] = struct{}{}
	return element
}

func (l *list) PushBack(v interface{}) *ListItem {
	element := &ListItem{Value: v}
	if l.items == nil {
		l.head = element
		l.tail = element
		l.items = make(map[*ListItem]struct{})
		l.items[element] = struct{}{}
		return element
	}
	element.Prev = l.tail
	l.tail.Next = element
	l.tail = element
	l.items[element] = struct{}{}
	return element
}
func (l *list) Remove(i *ListItem) {
	if i == l.tail {
		l.tail = i.Prev
		l.tail.Next = nil
	} else if i == l.head {
		l.head = i.Next
		l.head.Prev = nil
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	delete(l.items, i)
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head {
		return
	}
	if _, ok := l.items[i]; !ok {
		return
	}
	i.Prev.Next = i.Next
	if i != l.tail {
		i.Next.Prev = i.Prev
	}
	i.Next = l.head
	i.Prev = nil

	l.head.Prev = i
	l.head = i
}
