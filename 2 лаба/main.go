package main

import (
	"errors"
	"fmt"
	"math"
)

func main() {
	// Задание 1
	ip := [4]byte{192, 168, 1, 1}
	fmt.Println("IP-адрес:", formatIP(ip))

	evens, err := listEven(5, 15)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Чётные числа:", evens)
	}

	// Задание 2
	text := "hello world"
	charCounts := countChars(text)
	fmt.Println("Частота символов:")
	for char, count := range charCounts {
		fmt.Printf("%c: %d\n", char, count) // %c выводит символ, а не его код
	}

	// Задание 3
	a := Point{0, 0}
	b := Point{3, 0}
	c := Point{0, 4}

	segment := Segment{Start: a, End: b}
	length := segment.Length()
	fmt.Println("Расстояние между точками: ", length)

	triangle := Triangle{A: a, B: b, C: c}
	circle := Circle{Center: a, Radius: 5}

	printArea(triangle)
	printArea(circle)

	// Задание 4
	values := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	fmt.Println("Исходный срез:", values)

	square := func(x float64) float64 {
		return x * x
	}

	squaredValues := Map(values, square)

	fmt.Println("Срез после применения функции (копия):", squaredValues)
	fmt.Println("Исходный срез после Map (остался неизменным):", values)
}

// Задание 1. Массивы и срезы.
func formatIP(ip [4]byte) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}

func listEven(start, end int) ([]int, error) {
	if start > end {
		return nil, errors.New("левая граница больше правой")
	}
	var evens []int
	for i := start; i <= end; i++ {
		if i%2 == 0 {
			evens = append(evens, i)
		}
	}
	return evens, nil
}

// Задание 2. Карты.
func countChars(input string) map[rune]int {
	result := map[rune]int{}
	for i := range input {
		char := rune(input[i])
		result[char]++
	}
	return result
}

// Задание 3. Структуры, методы и интерфейсы.
type Point struct {
	X, Y float64
}

type Segment struct {
	Start, End Point
}

func distance(p1, p2 Point) float64 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (s Segment) Length() float64 {
	return distance(s.Start, s.End)
}

type Triangle struct {
	A, B, C Point
}

func (t Triangle) Area() float64 {
	a := distance(t.A, t.B)
	b := distance(t.B, t.C)
	c := distance(t.C, t.A)

	p := (a + b + c) / 2

	return math.Sqrt(p * (p - a) * (p - b) * (p - c))
}

type Circle struct {
	Center Point
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

type Shape interface {
	Area() float64
}

func printArea(s Shape) {
	result := s.Area()
	fmt.Printf("Площадь фигуры: %.2f\n", result)
}

// Задание 4

func Map(slice []float64, square func(float64) float64) []float64 {
	result := make([]float64, len(slice))
	copy(result, slice)

	for i := range result {
		result[i] = square(result[i])
	}

	return result
}
