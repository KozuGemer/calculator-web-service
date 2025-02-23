package utils_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/KozuGemer/calculator-web-service/utils"
)

// calcHandler обрабатывает HTTP-запросы для теста API.
func calcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, err := utils.Calc(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := map[string]float64{"result": result}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func TestCalcHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(calcHandler))
	defer server.Close()

	tests := []struct {
		name           string
		expression     string
		expectedResult map[string]float64
		expectedStatus int
	}{
		{
			name:           "valid expression with unary minus",
			expression:     "~2+2",
			expectedResult: map[string]float64{"result": 0}, // Ожидаемый результат
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid expression with stepping-stone",
			expression:     "2^2",
			expectedResult: map[string]float64{"result": 4}, // Ожидаемый результат
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid expression with hard expression",
			expression:     "~2^2+2*(~12)^4-8",
			expectedResult: map[string]float64{"result": 41468}, // Ожидаемый результат
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid expression with simple expression",
			expression:     "2+2-4",
			expectedResult: map[string]float64{"result": 0}, // Ожидаемый результат
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid expression with / expression",
			expression:     "45/9",
			expectedResult: map[string]float64{"result": 5}, // Ожидаемый результат
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid expression with * expression",
			expression:     "5*9",
			expectedResult: map[string]float64{"result": 45}, // Ожидаемый результат
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём JSON-запрос
			body, _ := json.Marshal(map[string]string{"expression": tt.expression})
			req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Отправляем запрос
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Проверяем статус-код ответа
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			// Читаем JSON-ответ
			var actualResult map[string]float64
			bodyBytes, _ := io.ReadAll(resp.Body)
			json.Unmarshal(bodyBytes, &actualResult)

			// Проверяем, что ответ совпадает с ожидаемым
			if !reflect.DeepEqual(actualResult, tt.expectedResult) {
				t.Errorf("expected result: %v, got: %v", tt.expectedResult, actualResult)
			}
		})
	}
}
