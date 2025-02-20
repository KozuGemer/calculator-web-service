package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/KozuGemer/calculator-web-service/utils"
)

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	expression := r.Form.Get("expression")
	result, err := utils.Calc(expression)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%v", result)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("site/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/calculate", calculateHandler)
	http.Handle("/style.css", http.FileServer(http.Dir("site")))

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
