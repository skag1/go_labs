package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/greet", handleGreet)
	http.HandleFunc("/add", handleAddition)
	http.HandleFunc("/sub", handleSubtraction)
	http.HandleFunc("/mul", handleMultiplication)
	http.HandleFunc("/div", handleDivision)
	http.HandleFunc("/charcount", handleCharCount)

	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}

// Задание 1
func handleGreet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	age := r.URL.Query().Get("age")

	if name == "" || age == "" {
		http.Error(w, "Параметры 'name' и 'age' обязательны", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Меня зовут %s, мне %s лет", name, age)
}

// Задание 2
func handleAddition(w http.ResponseWriter, r *http.Request) {
	handleMathOperation(w, r, "+")
}

func handleSubtraction(w http.ResponseWriter, r *http.Request) {
	handleMathOperation(w, r, "-")
}

func handleMultiplication(w http.ResponseWriter, r *http.Request) {
	handleMathOperation(w, r, "*")
}

func handleDivision(w http.ResponseWriter, r *http.Request) {
	handleMathOperation(w, r, "/")
}

func handleMathOperation(w http.ResponseWriter, r *http.Request, operation string) {
	query := r.URL.Query()
	aStr := query.Get("a")
	bStr := query.Get("b")

	if aStr == "" || bStr == "" {
		http.Error(w, "Параметры 'a' и 'b' обязательны", http.StatusBadRequest)
		return
	}

	a, errA := strconv.ParseFloat(aStr, 64)
	b, errB := strconv.ParseFloat(bStr, 64)

	if errA != nil || errB != nil {
		http.Error(w, "Параметры 'a' и 'b' должны быть числами", http.StatusBadRequest)
		return
	}

	var result float64
	switch operation {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			http.Error(w, "Деление на ноль невозможно", http.StatusBadRequest)
			return
		}
		result = a / b
	}

	fmt.Fprintf(w, "Результат: %f", result)
}

// Задание 3
func handleCharCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Text string `json:"text"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.Text == "" {
		http.Error(w, "Некорректный JSON или поле 'text' отсутствует", http.StatusBadRequest)
		return
	}

	charCount := make(map[string]int)
	for _, char := range strings.Split(requestBody.Text, "") {
		charCount[char]++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(charCount)
}
