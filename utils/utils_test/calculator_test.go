package utils_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KozuGemer/calculator-web-service/utils"
)

// Функция-обработчик для тестирования API
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
		expectedResult string
		expectedStatus int
	}{
		{
			name:           "Exponentiation test",
			expression:     "3^3+3",
			expectedResult: `{"result":30}`,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём JSON-запрос
			body, _ := json.Marshal(map[string]string{"expression": tt.expression})
			resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer resp.Body.Close()

			responseBody, _ := io.ReadAll(resp.Body)
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if string(responseBody) != tt.expectedResult {
				t.Errorf("expected result: %v, got: %v", tt.expectedResult, string(responseBody))
			}
		})
	}
}
