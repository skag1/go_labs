package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")

	//3.2
	fmt.Println(hello("Андрей"))

	//3.3
	printEven(10, 7)

	//3.4
	fmt.Println(apply(5, 5, "+"))

}

func hello(name string) string {
	return ("Привет, " + name)
}

func printEven(a, b int64) error {
	if a > b {
		return errors.New("a > b")
	}
	for i := a; i <= b; i++ {
		if i%2 == 0 {
			fmt.Println(i)
		}
	}
	return nil
}

func apply(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, errors.New("b = 0")
		}
		return a / b, nil
	default:
		return 0, errors.New("not an operator")
	}
}
