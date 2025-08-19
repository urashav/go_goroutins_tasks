// https://stepik.org/lesson/745849/step/3?auth=login&unit=871359

package main

import (
	"errors"
	"fmt"
)

var ErrFull = errors.New("Queue is full")
var ErrEmpty = errors.New("Queue is empty")

// начало решения

// Queue - FIFO-очередь на n элементов
type Queue struct {
	Items chan int
}

// Get возвращает очередной элемент.
// Если элементов нет и block = false -
// возвращает ошибку.
func (q Queue) Get(block bool) (int, error) {
	select {
	case item := <-q.Items:
		return item, nil
	default:
		if !block {
			return 0, ErrEmpty
		}
	}
	return <-q.Items, nil
}

// Put помещает элемент в очередь.
// Если очередь заполнена и block = false -
// возвращает ошибку.
func (q Queue) Put(val int, block bool) error {
	select {
	case q.Items <- val:
	default:
		if !block {
			return ErrFull
		}

		q.Items <- val

	}
	return nil
}

// MakeQueue создает новую очередь
func MakeQueue(n int) Queue {
	return Queue{
		make(chan int, n),
	}
}

// конец решения

func main() {
	// Проверка без блокировки

	//q := MakeQueue(2)
	//
	//err := q.Put(1, false)
	//fmt.Println("put 1:", err)
	//
	//err = q.Put(2, false)
	//fmt.Println("put 2:", err)
	//
	//err = q.Put(3, true)
	//fmt.Println("put 3:", err)
	//
	//res, err := q.Get(false)
	//fmt.Println("get:", res, err)
	//
	//res, err = q.Get(false)
	//fmt.Println("get:", res, err)
	//
	//res, err = q.Get(false)
	//fmt.Println("get:", res, err)

	// Проверка с блокировкой
	q := MakeQueue(1)

	go func() {
		err := q.Put(11, true)
		fmt.Println("put 11:", err)
		// put 11: <nil>

		err = q.Put(12, true)
		fmt.Println("put 12:", err)
		// put 12: <nil>
	}()

	res, err := q.Get(true)
	fmt.Println("get:", res, err)
	// get: 11 <nil>

	res, err = q.Get(true)
	fmt.Println("get:", res, err)
	// get: 12 <nil>
}
