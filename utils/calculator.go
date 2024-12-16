package utils

import (
	"errors"
	"strconv"
	"unicode"

	"github.com/KozuGemer/calculator-web-service/models"
)

// precedence - возвращает приоритет операции
func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	case '(':
		return 0
	}
	return -1
}

// applyOperator - применяет оператор к двум числам
func applyOperator(a, b float64, op rune) float64 {
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		if b == 0 {
			panic("division by zero")
		}
		return a / b
	}
	return 0
}

// isOperator - проверяет, является ли символ оператором
func isOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/'
}

// isValidExpression - проверяет, является ли выражение допустимым
func isValidExpression(expression string) error {
	openBrackets := 0
	prevChar := ' '

	for _, c := range expression {
		if unicode.IsSpace(c) {
			continue
		}
		if c == '(' {
			openBrackets++
		} else if c == ')' {
			openBrackets--
		} else if isOperator(c) {
			if isOperator(prevChar) {
				return errors.New("invalid expression: consecutive operators")
			}
		} else if !unicode.IsDigit(c) && c != '.' {
			return errors.New("invalid expression: invalid character")
		}
		prevChar = c
	}

	if openBrackets != 0 {
		return errors.New("invalid expression: unbalanced parentheses")
	}
	if prevChar == '(' || isOperator(prevChar) {
		return errors.New("invalid expression: expression cannot end with an operator or '('")
	}

	return nil
}

// Calc - вычисляет результат математического выражения
func Calc(expression string) (float64, error) {
	// Проверка на корректность выражения
	if err := isValidExpression(expression); err != nil {
		return 0, err
	}

	var values models.Stack
	var ops models.Stack

	for i := 0; i < len(expression); i++ {
		c := rune(expression[i])

		// Пропуск пробелов
		if unicode.IsSpace(c) {
			continue
		}

		// Если текущий символ - число
		if unicode.IsDigit(c) {
			var sb string
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				sb += string(expression[i])
				i++
			}
			i-- // Корректируем индекс
			num, err := strconv.ParseFloat(sb, 64)
			if err != nil {
				return 0, err
			}
			values.Push(num)
		} else if c == '(' {
			ops.Push(float64(c)) // Используем float64 для хранения рун
		} else if c == ')' {
			for len(ops.Items) > 0 && ops.Items[len(ops.Items)-1] != float64('(') {
				val2, _ := values.Pop()
				val1, _ := values.Pop()
				op, _ := ops.Pop()
				values.Push(applyOperator(val1, val2, rune(op)))
			}
			if len(ops.Items) > 0 {
				ops.Pop() // Удаляем '('
			}
		} else { // Оператор
			for len(ops.Items) > 0 && precedence(c) <= precedence(rune(ops.Items[len(ops.Items)-1])) {
				val2, _ := values.Pop()
				val1, _ := values.Pop()
				op, _ := ops.Pop()
				values.Push(applyOperator(val1, val2, rune(op)))
			}
			ops.Push(float64(c))
		}
	}

	// Обработка оставшихся операций
	for len(ops.Items) > 0 {
		val2, _ := values.Pop()
		val1, _ := values.Pop()
		op, _ := ops.Pop()
		values.Push(applyOperator(val1, val2, rune(op)))
	}

	return values.Pop()
}
