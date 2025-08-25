// https://stepik.org/lesson/1363480/step/9?auth=login&unit=1379378

// Ограничитель вызовов
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrBusy = errors.New("busy")
var ErrCanceled = errors.New("canceled")

// начало решения
type Handler struct {
	count    int
	ticker   *time.Ticker
	canceled bool
	lock     sync.Mutex
}

func (h *Handler) Close() {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.canceled {
		return
	}

	h.ticker.Stop()
	h.canceled = true
}

func NewHandler() *Handler {
	ticker := time.NewTicker(time.Second)

	return &Handler{
		ticker: ticker,
	}
}

// throttle следит, чтобы функция fn выполнялась не более limit раз в секунду.
// Возвращает функции handle (выполняет fn с учетом лимита) и cancel (останавливает ограничитель).
func throttle(limit int, fn func()) (handle func() error, cancel func()) {
	handler := NewHandler()

	cancelFunc := func() {
		handler.Close()
	}

	handlerFunc := func() error {
		handler.lock.Lock()
		defer handler.lock.Unlock()
		if handler.canceled {
			return ErrCanceled
		}

		select {
		case <-handler.ticker.C:
			handler.count = 0
		default:
			if handler.count >= limit {
				return ErrBusy
			}
		}

		handler.count++
		fn()
		return nil
	}
	return handlerFunc, cancelFunc
}

// конец решения

func main() {
	work := func() {
		fmt.Print(".")
	}

	handle, cancel := throttle(5, work)
	defer cancel()

	const n = 8
	var nOK, nErr int
	for i := 0; i < n; i++ {
		err := handle()
		if err == nil {
			nOK += 1
		} else {
			nErr += 1
		}
	}
	fmt.Println()
	fmt.Printf("%d calls: %d OK, %d busy\n", n, nOK, nErr)
}
