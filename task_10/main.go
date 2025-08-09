// https://stepik.org/lesson/740355/step/7?auth=login&unit=871357

package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// say печатает фразу от имени обработчика
func say(id int, phrase string) {
	for _, word := range strings.Fields(phrase) {
		fmt.Printf("Worker #%d says: %s...\n", id, word)
		dur := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(dur)
	}
}

// начало решения

// makePool создает пул на n обработчиков
// возвращает функции handle и wait
func makePool(n int, handler func(int, string)) (func(string), func()) {
	// создайте пул на n обработчиков
	// используйте для канала имя pool и тип chan int
	// определите функции handle() и wait()

	// handle() выбирает токен из пула
	// и обрабатывает переданную фразу через handler()

	// wait() дожидается, пока все токены вернутся в пул

	/////////////////////////////////////
	// Создаем канал с буфером объемом n //
	pool := make(chan int, n)

	// Отправляем в канал n ID
	for i := range n {
		pool <- i
	}

	// В функции получаем ID из пула, обрабатываем фразу и возвращаем ID в пул
	handle := func(str string) {
		id := <-pool
		handler(id, str)
		pool <- id
	}

	// Ожидаем все ID из пулов
	wait := func() {
		for range n {
			<-pool
		}
		close(pool)
	}

	return handle, wait
}

// конец решения

func main() {
	phrases := []string{
		"go is awesome",
		"cats are cute",
		"rain is wet",
		"channels are hard",
		"floor is lava",
	}

	handle, wait := makePool(2, say)
	for _, phrase := range phrases {
		handle(phrase)
	}
	wait()
}
