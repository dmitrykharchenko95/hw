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
	//List// Remove me after realization.
	Data 		[]ListItem
	FrontItem 	ListItem
	BackItem 	ListItem
	// Place your code here.
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return len(l.Data)
}

func (l list) Front() *ListItem  {
	for _, f := range l.Data {
		if f.Prev == nil {
			return &f
		}
		f = *f.Prev
	}
	return nil
}

func (l list) Back() *ListItem {
	for _, b := range l.Data {
		if b.Next == nil {
			return &b
		}
		b = *b.Next
	}
	return nil
}

func (l list) PushFront (v interface{}) *ListItem  {
	f := ListItem {
		v,
		l.Front(),
		nil,
	}
	l.Front().Prev = &f
	return &f
}

func (l list) PushBack (v interface{}) *ListItem  {
	b := ListItem {
		v,
		nil,
		l.Back(),
	}
	l.Back().Next = &b
	return &b
}

func (l list) Remove(i *ListItem)  {
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	i.Next, i.Prev = nil, nil
}

func (l list) MoveToFront(i *ListItem)()  {
	i.Prev = nil
	i.Next = l.Front()
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
}

