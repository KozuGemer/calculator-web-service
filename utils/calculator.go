package utils

import (
	"errors"
	"fmt"
	"math"
	"regexp"
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
	case '^': // Возведение в степень выше умножения
		return 3
	case '~': // Угарный минус выше всех
		return 4
	case '(':
		return 0
	}
	return -1
}

func applyOperator(a, b float64, op rune) float64 {
	switch op {
	case '+':
		fmt.Println("Applying", string(op), "to", a, "and", b)
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
	case '^': // Возведение a в степень b
		return math.Pow(a, b)

	case '~': // Угарный минус
		return -b
	}

	panic(fmt.Sprintf("Unknown operator: %c", op))
}

// isOperator - проверяет, является ли символ оператором
func isOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '^' || c == '~'
}

// Tokenize - разбивает выражение на операнды и операторы
func Tokenize(expression string) ([]string, error) {
	// Регулярное выражение для поиска чисел, операторов и скобок
	re := regexp.MustCompile(`\d+(\.\d+)?|\+|\-|\*|\/|\^|\(|\)`)
	tokens := re.FindAllString(expression, -1)

	if len(tokens) == 0 {
		return nil, errors.New("invalid expression: empty or incorrect format")
	}

	return tokens, nil
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
			if c != '~' && isOperator(prevChar) { // Разрешаем `~` после оператора
				return errors.New("invalid expression: consecutive operators")
			}
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

	// Разбираем выражение на токены
	tokens, err := Tokenize(expression)
	if err != nil {
		return 0, err
	}

	var values models.Stack
	var ops models.Stack

	for _, token := range tokens {
		// Если число
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			values.Push(num)
		} else if token == "(" {
			ops.Push(float64('('))
		} else if token == ")" {
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
			op := rune(token[0])
			for len(ops.Items) > 0 && precedence(op) < precedence(rune(ops.Items[len(ops.Items)-1])) {
				val2, _ := values.Pop()
				val1, _ := values.Pop()
				op, _ := ops.Pop()
				values.Push(applyOperator(val1, val2, rune(op)))
			}
			ops.Push(float64(op))
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
