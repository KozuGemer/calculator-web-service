package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KozuGemer/calculator-web-service/models"
	"github.com/KozuGemer/calculator-web-service/utils"
)

// CalculateHandler - обработчик запросов для вычислений
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	result, err := utils.Calc(req.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(models.Response{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.Response{Result: fmt.Sprintf("%f", result)})
}
