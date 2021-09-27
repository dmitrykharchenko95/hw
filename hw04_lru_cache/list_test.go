package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
	t.Run("additions", func(t *testing.T) {
		// перезапись элементов l1 в l2 в обратном порядке с последующим удалением l1
		l1 := NewList()
		l2 := NewList()

		for i := 0; i < 10; i++ {
			l1.PushFront(i)
		}
		for i := l1.Back(); i != nil; i = i.Prev {
			l2.PushFront(i.Value)
			l1.Remove(i)
		}
		elems := make([]int, 0, l2.Len())
		for i := l2.front; i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Nil(t, l1.front.Value, l1.front.Prev, l1.front.Next, l1.back.Value, l1.back.Prev, l1.back.Next)
		require.Equal(t, []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, elems)
	})
}
