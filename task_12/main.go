// https://stepik.org/lesson/740281/step/8?auth=login&unit=871358
package main

import (
	"fmt"
	"sync"
	"time"
)

// rangeGen отправляет в канал числа от start до stop-1
func rangeGen(start, stop int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; i < stop; i++ {
			time.Sleep(50 * time.Millisecond)
			out <- i
		}

	}()
	return out
}

// начало решения

// merge выбирает числа из входных каналов и отправляет в выходной
func merge(channels ...<-chan int) <-chan int {
	// объедините все исходные каналы в один выходной
	// последовательное объединение НЕ подходит
	out := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, c := range channels {
		go func() {
			defer wg.Done()
			for val := range c {
				out <- val
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// конец решения

func main() {
	in1 := rangeGen(11, 15)
	in2 := rangeGen(21, 25)
	in3 := rangeGen(31, 35)

	start := time.Now()
	merged := merge(in1, in2, in3)
	for val := range merged {
		fmt.Print(val, " ")
	}
	fmt.Println()
	fmt.Println("Took", time.Since(start))
}
