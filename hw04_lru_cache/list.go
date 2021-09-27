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
	size  int
	// Place your code here.
}

func NewList() *list {
	return new(list)
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	if l.size == 0 {
		return nil
	}
	return l.front
}

func (l *list) Back() *ListItem {
	if l.size == 0 {
		return nil
	}
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	if l.size == 0 {
		l.front = &ListItem{
			v,
			nil,
			nil,
		}
		l.back = l.front
	} else {
		newFront := &ListItem{
			v,
			l.front,
			nil,
		}
		l.front.Prev = newFront
		l.front = newFront
	}
	l.size++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	if l.size == 0 {
		l.back = &ListItem{
			v,
			nil,
			nil,
		}
		l.front = l.back
	} else {
		newBack := &ListItem{
			v,
			nil,
			l.back,
		}
		l.back.Next = newBack
		l.back = newBack
	}
	l.size++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.Len() == 1:
		i.Value = nil
	case i == l.front:
		i.Next.Prev = nil
		l.front = i.Next
	case i == l.back:
		i.Prev.Next = nil
		l.back = i.Prev
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Value = nil
	}
	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	switch {
	case l.front == i:
	case l.back == i:
		l.back = l.back.Prev
		l.back.Next = nil

		i.Next = l.front
		i.Prev = nil
		l.front = i

		i.Next.Prev = i
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev

		i.Next = l.front
		i.Prev = nil
		l.front = i

		i.Next.Prev = i
	}
}
