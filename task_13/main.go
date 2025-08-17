package main

import (
	"fmt"
	"math/rand"
)

// начало решения
type Result struct {
	Source   string
	Reversed string
}

// генерит случайные слова из 5 букв
// с помощью randomWord(5)
func generate(cancel <-chan struct{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case out <- randomWord(5):
			case <-cancel:
				return
			}
		}
	}()
	return out
}

// выбирает слова, в которых не повторяются буквы,
// abcde - подходит
// abcda - не подходит
func takeUnique(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string)

	isValidWord := func(word string) bool {
		letters := make(map[rune]bool)

		for _, r := range word {
			if letters[r] {
				return false
			}
			letters[r] = true
		}
		return true
	}

	go func() {
		defer close(out)

		for word := range in {
			if isValidWord(word) {
				select {
				case out <- word:
				case <-cancel:
					return
				}
			}
		}
	}()
	return out
}

// переворачивает слова
// abcde -> edcba
func reverse(cancel <-chan struct{}, in <-chan string) <-chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		result := Result{}
		for word := range in {
			runes := []rune(word)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}

			result.Source = word
			result.Reversed = string(runes)
			select {
			case out <- result:
			case <-cancel:
			}
		}
	}()
	return out
}

// объединяет c1 и c2 в общий канал
func merge(cancel <-chan struct{}, c1, c2 <-chan Result) <-chan Result {
	out := make(chan Result)

	go func() {
		defer close(out)
		for val := range c1 {
			select {
			case out <- val:
			case <-cancel:
				return
			}
		}
		for val := range c2 {
			select {
			case out <- val:
			case <-cancel:
				return
			}
		}
	}()

	return out
}

// печатает первые n результатов
func print(cancel <-chan struct{}, in <-chan Result, n int) {
	for i := 0; i < n; i++ {
		word := <-in
		fmt.Printf("%s -> %s\n", word.Source, word.Reversed)
	}
}

// конец решения

// генерит случайное слово из n букв
func randomWord(n int) string {
	const letters = "aeiourtnsl"
	chars := make([]byte, n)
	for i := range chars {
		chars[i] = letters[rand.Intn(len(letters))]
	}
	return string(chars)
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	c1 := generate(cancel)
	c2 := takeUnique(cancel, c1)
	c3_1 := reverse(cancel, c2)
	c3_2 := reverse(cancel, c2)
	c4 := merge(cancel, c3_1, c3_2)
	print(cancel, c4, 10)
}
