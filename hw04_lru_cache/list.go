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
	listLen int
	head    *ListItem
	tail    *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.listLen
}

func (l *list) Front() *ListItem {
	if l.listLen == 0 {
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
	if l.listLen == 0 {
		l.head = element
		l.tail = element
		l.listLen++
		return element
	}
	element.Next = l.head
	l.head.Prev = element
	l.head = element
	l.listLen++
	return element
}

func (l *list) PushBack(v interface{}) *ListItem {
	element := &ListItem{Value: v}
	if l.listLen == 0 {
		l.head = element
		l.tail = element
		l.listLen++
		return element
	}
	element.Prev = l.tail
	l.tail.Next = element
	l.tail = element
	l.listLen++
	return element
}

func (l *list) Remove(i *ListItem) {
	switch i {
	case l.tail:
		l.tail = i.Prev
		l.tail.Next = nil
	case l.head:
		l.head = i.Next
		l.head.Prev = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.listLen--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.head {
		return
	}
	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	i.Next = l.head
	i.Prev, l.head.Prev = l.head.Prev, i.Prev
	if i == l.tail {
		l.tail = l.head
	}
	l.head = i
}
