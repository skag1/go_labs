package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Создаем приложение
	a := app.New()
	w := a.NewWindow("Калькулятор")

	// Переменная для хранения текущего ввода
	var currentInput string

	// Функция обновления текста в поле
	updateDisplay := func() {
		display.SetText(currentInput)
	}

	// Функция для обработки нажатия кнопки с числом
	onNumberButtonTapped := func(number string) {
		currentInput += number
		updateDisplay()
	}

	// Функция для обработки арифметических операций
	onOperatorButtonTapped := func(operator string) {
		currentInput += " " + operator + " "
		updateDisplay()
	}

	// Функция для вычисления выражения
	onEqualsButtonTapped := func() {
		result, err := calculate(currentInput)
		if err != nil {
			currentInput = "Ошибка"
		} else {
			currentInput = fmt.Sprintf("%v", result)
		}
		updateDisplay()
	}

	// Функция для очистки экрана
	onClearButtonTapped := func() {
		currentInput = ""
		updateDisplay()
	}

	// Виджет для отображения ввода
	display := widget.NewLabel("0")
	display.Alignment = fyne.TextAlignCenter

	// Кнопки калькулятора
	buttons := []struct {
		label  string
		action func()
	}{
		{"7", func() { onNumberButtonTapped("7") }},
		{"8", func() { onNumberButtonTapped("8") }},
		{"9", func() { onNumberButtonTapped("9") }},
		{"/", func() { onOperatorButtonTapped("/") }},
		{"4", func() { onNumberButtonTapped("4") }},
		{"5", func() { onNumberButtonTapped("5") }},
		{"6", func() { onNumberButtonTapped("6") }},
		{"*", func() { onOperatorButtonTapped("*") }},
		{"1", func() { onNumberButtonTapped("1") }},
		{"2", func() { onNumberButtonTapped("2") }},
		{"3", func() { onNumberButtonTapped("3") }},
		{"-", func() { onOperatorButtonTapped("-") }},
		{"0", func() { onNumberButtonTapped("0") }},
		{".", func() { onNumberButtonTapped(".") }},
		{"=", onEqualsButtonTapped},
		{"+", func() { onOperatorButtonTapped("+") }},
		{"C", onClearButtonTapped},
	}

	// Создание контейнера для кнопок
	var buttonWidgets []fyne.CanvasObject
	for _, btn := range buttons {
		buttonWidgets = append(buttonWidgets, widget.NewButton(btn.label, btn.action))
	}

	// Создание сетки для кнопок
	buttonGrid := container.NewGridWithColumns(4, buttonWidgets...)

	// Компоновка окна
	w.SetContent(container.NewVBox(display, buttonGrid))

	// Открытие окна
	w.Resize(fyne.NewSize(300, 400))
	w.ShowAndRun()
}

// Функция для вычисления выражения
func calculate(input string) (float64, error) {
	// Преобразуем строку в арифметическое выражение
	input = input
	// Парсим и вычисляем результат
	result, err := evalExpression(input)
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Простая функция для вычисления выражений
func evalExpression(expression string) (float64, error) {
	// Для простоты будем использовать только базовые операции и парсинг
	// Разделяем по пробелам и выполняем операции
	parts := parseParts(expression)
	num1, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}

	operator := parts[1]
	num2, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, err
	}

	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0, fmt.Errorf("деление на ноль невозможно")
		}
		return num1 / num2, nil
	default:
		return 0, fmt.Errorf("неизвестная операция: %s", operator)
	}
}

// Функция для парсинга части выражения
func parseParts(input string) []string {
	// Разбиваем строку на части
	return []string{"10", "+", "5"} // Заглушка, нужно будет реализовать более сложное парсинг
}
