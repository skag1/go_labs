package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

var db *sql.DB
var secretKey = "your-256-bit-secret"

func main() {

	var err error
	db, err = sql.Open("sqlite", "./my_database.db")
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создание таблицы
	query := `CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		surname TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/register", registerCustomerHandler)
	mux.HandleFunc("/login", loginCustomerHandler)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/customer", getCustomerHandler)

	mux.Handle("/customer", authRequired(protectedMux))

	fmt.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Неверный или просроченный токен", http.StatusForbidden)
			return
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			http.Error(w, "Ошибка извлечения данных из токена", http.StatusUnauthorized)
			return
		}
		r.Header.Set("UserID", claims.Subject)

		next.ServeHTTP(w, r)
	})
}

func registerCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST-запросы поддерживаются", http.StatusMethodNotAllowed)
		return
	}

	var customer Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	var existingEmail string
	err = db.QueryRow(`SELECT email FROM customers WHERE email = ?`, customer.Email).Scan(&existingEmail)
	if err == nil {
		http.Error(w, "Покупатель с таким адресом электронной почты уже существует", http.StatusConflict)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(customer.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка хеширования пароля: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`INSERT INTO customers (name, surname, email, password_hash) VALUES (?, ?, ?, ?)`,
		customer.Name, customer.Surname, customer.Email, string(passwordHash))
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка добавления пользователя: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Покупатель %s %s успешно зарегистрирован!", customer.Name, customer.Surname)
}

func loginCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST-запросы поддерживаются", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	var customer Customer
	err = db.QueryRow(`SELECT id, email, password_hash FROM customers WHERE email = ?`, request.Email).
		Scan(&customer.ID, &customer.Email, &customer.PasswordHash)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(request.Password))
	if err != nil {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("%d", customer.ID),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка создания токена: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func getCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только GET-запросы поддерживаются", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		http.Error(w, "Ошибка авторизации", http.StatusForbidden)
		return
	}

	var customer Customer
	err := db.QueryRow(`SELECT id, name, surname, email FROM customers WHERE id = ?`, userID).
		Scan(&customer.ID, &customer.Name, &customer.Surname, &customer.Email)
	if err != nil {
		http.Error(w, "Покупатель не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}
