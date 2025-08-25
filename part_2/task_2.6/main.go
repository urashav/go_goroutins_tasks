// https://stepik.org/lesson/867253/step/9?auth=login&unit=871362

// Карта на RWMutex
package main

import (
	"fmt"
	"sync"
)

// начало решения

// Counter представляет безопасную карту частот слов.
// Ключ - строка, значение - целое число.
type Counter struct {
	// не меняйте название и тип поля lock
	lock   sync.RWMutex
	values map[string]int
}

// Increment увеличивает значение по ключу на 1.
func (c *Counter) Increment(str string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.values[str]++
}

// Value возвращает значение по ключу.
func (c *Counter) Value(str string) int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.values[str]
}

// Range проходит по всем записям карты,
// и для каждой вызывает функцию fn, передавая в нее ключ и значение.
func (c *Counter) Range(fn func(key string, val int)) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for key, value := range c.values {
		fn(key, value)
	}
}

// NewCounter создает новую карту частот.
func NewCounter() *Counter {
	return &Counter{values: make(map[string]int)}
}

// конец решения

func main() {
	counter := NewCounter()

	var wg sync.WaitGroup
	wg.Add(3)

	increment := func(key string, val int) {
		defer wg.Done()
		for ; val > 0; val-- {
			counter.Increment(key)
		}
	}

	go increment("one", 100)
	go increment("two", 200)
	go increment("three", 300)

	wg.Wait()

	fmt.Println("two:", counter.Value("two"))

	fmt.Print("{ ")
	counter.Range(func(key string, val int) {
		fmt.Printf("%s:%d ", key, val)
	})
	fmt.Println("}")
}
