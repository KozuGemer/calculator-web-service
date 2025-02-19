package utils_test

import (
	"testing"

	"github.com/KozuGemer/calculator-web-service/utils"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedResult float64
		expectedError  string
	}{
		{
			name:           "Valid Expression",
			expression:     "2+2^2",
			expectedResult: 6,
			expectedError:  "",
		},
		{
			name:           "Invalid Expression - Consecutive Operators",
			expression:     "2+2*-",
			expectedResult: 0,
			expectedError:  "invalid expression: consecutive operators",
		},
		{
			name:           "Invalid Expression - Non-Mathematical Input",
			expression:     "Hello World",
			expectedResult: 0,
			expectedError:  "invalid expression: invalid character",
		},
		{
			name:           "Server Error - Division by Zero",
			expression:     "1/0",
			expectedResult: 0,
			expectedError:  "panic: division by zero",
		},
		{
			name:           "Invalid Expression - Unbalanced Parentheses",
			expression:     "(2+3",
			expectedResult: 0,
			expectedError:  "invalid expression: unbalanced parentheses",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.expectedError != "panic: division by zero" {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()

			result, err := utils.Calc(tt.expression)

			if err != nil && tt.expectedError == "" {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && tt.expectedError != "" {
				t.Errorf("expected error: %v, got none", tt.expectedError)
			}

			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err.Error())
			}

			if result != tt.expectedResult {
				t.Errorf("expected result: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}
