package main

import (
	"fmt"
	"strings"
	"sync"
	"unicode"
)

// counter хранит количество цифр в каждом слове.
// ключ карты - слово, а значение - количество цифр в слове.
type counter map[string]int

// countDigitsInWords считает количество цифр в словах фразы
func countDigitsInWords(phrase string) counter {
	words := strings.Fields(phrase)
	syncStats := sync.Map{}

	var wg sync.WaitGroup

	// начало решения
	wg.Add(len(words))
	for _, word := range words {
		go func() {
			defer wg.Done()
			n := countDigits(word)
			syncStats.Store(word, n)
		}()
	}
	wg.Wait()
	// Посчитайте количество цифр в словах,
	// используя отдельную горутину для каждого слова.

	// Чтобы записать результаты подсчета,
	// используйте syncStats.Store(word, count)

	// В результате syncStats должна содержать слова
	// и количество цифр в каждом.

	// конец решения

	return asStats(&syncStats)
}

// countDigits возвращает количество цифр в строке
func countDigits(str string) int {
	count := 0
	for _, char := range str {
		if unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

// asStats преобразует статистику из sync.Map в обычную карту
func asStats(m *sync.Map) counter {
	stats := counter{}
	m.Range(func(word, count any) bool {
		stats[word.(string)] = count.(int)
		return true
	})
	return stats
}

// printStats печатает слова и количество цифр в каждом
func printStats(stats counter) {
	for word, count := range stats {
		fmt.Printf("%s: %d\n", word, count)
	}
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	counts := countDigitsInWords(phrase)
	printStats(counts)
}
