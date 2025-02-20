package models

import (
	"errors"
	"fmt"
)

// Stack - структура для стека
type Stack struct {
	Items []float64
}

// Push - добавляет элемент в стек
func (s *Stack) Push(item float64) {
	fmt.Println("Stack after push:", append(s.Items, item))
	s.Items = append(s.Items, item)
}

// Pop - извлекает элемент из стека
func (s *Stack) Pop() (float64, error) {
	if len(s.Items) == 0 {
		fmt.Println("Pop called on empty stack!") // <- Лог для проверки
		return 0, errors.New("stack is empty")
	}
	item := s.Items[len(s.Items)-1]
	s.Items = s.Items[:len(s.Items)-1]
	fmt.Println("Popped:", item) // <- Лог, который покажет, что извлекли
	return item, nil
}
