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
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.front == nil {
		l.back = item
	} else {
		item.Next = l.front
		l.front.Prev = item
	}
	l.front = item

	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.back == nil {
		l.front = item
	} else {
		item.Prev = l.back
		l.back.Next = item
	}
	l.back = item

	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	if l.len > 0 {
		if i.Prev != nil {
			i.Prev.Next = i.Next
		} else {
			l.front = i.Next
		}

		if i.Next != nil {
			i.Next.Prev = i.Prev
		} else {
			l.back = i.Prev
		}

		i.Next = nil
		i.Prev = nil
		l.len--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.len > 0 && l.front != i {
		if l.back == i { // if it is the last element
			l.back = i.Prev
			l.back.Next = nil
		} else { // if it is in the middle
			i.Prev.Next = i.Next
			i.Next.Prev = i.Prev
		}

		// manage the front element
		l.front.Prev = i
		i.Next = l.front
		i.Prev = nil
		l.front = i
	}
}

func NewList() List {
	return new(list)
}
