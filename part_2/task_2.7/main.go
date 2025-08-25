// https://stepik.org/lesson/1363480/step/3?auth=login&unit=1379378

// Конкурентно-безопасная карта.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ConcMap - безопасная в многозадачной среде карта.
type ConcMap[K comparable, V any] struct {
	items map[K]V
	lock  sync.Mutex
}

// NewConcMap создает новую карту.
func NewConcMap[K comparable, V any]() *ConcMap[K, V] {
	return &ConcMap[K, V]{items: map[K]V{}}
}

// Get возвращает значение по ключу.
func (cm *ConcMap[K, V]) Get(key K) V {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	return cm.items[key]
}

// Set устанавливает значение по ключу.
func (cm *ConcMap[K, V]) Set(key K, val V) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.items[key] = val
}

// начало решения

// SetIfAbsent устанавливает новое значение по ключу
// и возвращает его, но только если такого ключа нет в карте.
// Если ключ уже есть - возвращает старое значение по ключу.
func (cm *ConcMap[K, V]) SetIfAbsent(key K, val V) V {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	var v V
	_, exists := cm.items[key]

	if !exists {
		cm.items[key] = val
		v = val
	} else {
		v = cm.items[key]
	}
	return v
}

// Compute устанавливает значение по ключу, применяя к нему функцию.
// Возвращает новое значение. Функция выполняется атомарно.
func (cm *ConcMap[K, V]) Compute(key K, f func(V) V) V {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	var v V
	v = f(cm.items[key])
	cm.items[key] = v

	return v
}

// конец решения

func getSet() {
	m := NewConcMap[string, int]()

	var wg sync.WaitGroup
	wg.Go(func() {
		m.Set("hello", rand.Intn(100))
	})
	wg.Go(func() {
		m.Set("hello", rand.Intn(100))
	})
	wg.Wait()

	fmt.Println("hello =", m.Get("hello"))
	// hello = 71 (случайное)
}

func setIfAbsent() {
	m := NewConcMap[string, int]()

	var wg sync.WaitGroup
	wg.Go(func() {
		time.Sleep(5 * time.Millisecond)
		m.SetIfAbsent("hello", 42)
	})
	wg.Go(func() {
		time.Sleep(10 * time.Millisecond)
		m.SetIfAbsent("hello", 84)
	})
	wg.Wait()

	fmt.Println("hello =", m.Get("hello"))
	// hello = 42 (от первой горутины)
}

func compute() {
	m := NewConcMap[string, int]()
	var wg sync.WaitGroup

	wg.Go(func() {
		for range 100 {
			m.Compute("hello", func(v int) int {
				return v + 1
			})
		}
	})

	wg.Go(func() {
		for range 100 {
			m.Compute("hello", func(v int) int {
				return v + 1
			})
		}
	})

	wg.Wait()
	fmt.Println("hello =", m.Get("hello"))
	// hello = 200 (каждая горутина увеличила hello на 100)
}

func main() {
	getSet()
	fmt.Println("---")
	setIfAbsent()
	fmt.Println("---")
	compute()
}
