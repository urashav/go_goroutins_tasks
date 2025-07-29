package main

import (
	"fmt"
	"time"
)

// gather выполняет переданные функции одновременно
// и возвращает срез с результатами, когда они готовы
func gather(funcs []func() any) []any {
	// начало решения
	type data struct {
		index int
		value int
	}

	result := make([]any, len(funcs))
	item := make(chan data, len(funcs))

	for i, f := range funcs {
		go func() {
			item <- data{i, f().(int)}
		}()
	}
	// выполните все переданные функции,
	for i := 0; i < len(funcs); i++ {
		num := <-item
		result[num.index] = num.value
	}
	// соберите результаты в срез
	// и верните его
	return result
	// конец решения
}

// squared возвращает функцию,
// которая считает квадрат n
func squared(n int) func() any {
	return func() any {
		time.Sleep(time.Duration(n) * 100 * time.Millisecond)
		return n * n
	}
}

func main() {
	funcs := []func() any{squared(2), squared(3), squared(4)}

	start := time.Now()
	nums := gather(funcs)
	elapsed := float64(time.Since(start)) / 1_000_000

	fmt.Println(nums)
	fmt.Printf("Took %.0f ms\n", elapsed)
}
