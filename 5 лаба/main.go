package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"
	"time"
)

// Задача 1: Работа с горутинами и каналами.
// count считывает числа из канала, возводит их в квадрат и выводит.
func count(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("Число: %d, Квадрат: %d\n", num, num*num)
	}
}

// Задача 2 и 3: Обработка изображений.

// filter обрабатывает изображение последовательно.
func filter(img *image.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := uint16((r + g + b) / 3)
			img.SetRGBA(x, y, color.RGBA{
				R: uint8(gray >> 8),
				G: uint8(gray >> 8),
				B: uint8(gray >> 8),
				A: 255,
			})
		}
	}
}

// filterParallel обрабатывает одну строку изображения параллельно.
func filterParallel(img *image.RGBA, y int, wg *sync.WaitGroup) {
	defer wg.Done()
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		r, g, b, _ := img.At(x, y).RGBA()
		gray := uint16((r + g + b) / 3)
		img.SetRGBA(x, y, color.RGBA{
			R: uint8(gray >> 8),
			G: uint8(gray >> 8),
			B: uint8(gray >> 8),
			A: 255,
		})
	}
}

// processSequential выполняет последовательную обработку изображения.
func processSequential(img image.Image, bounds image.Rectangle) *image.RGBA {
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)
	start := time.Now()
	filter(rgbaImg)
	fmt.Println("Время последовательной обработки:", time.Since(start))
	return rgbaImg
}

// processParallel выполняет параллельную обработку изображения.
func processParallel(img image.Image, bounds image.Rectangle) *image.RGBA {
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)
	var wg sync.WaitGroup
	start := time.Now()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go filterParallel(rgbaImg, y, &wg)
	}
	wg.Wait()
	fmt.Println("Время параллельной обработки:", time.Since(start))
	return rgbaImg
}

func main() {

	// Задача 1: Работа с горутинами и каналами.
	fmt.Println("Задача 1: Работа с горутинами и каналами")
	numChannel := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)
	go count(numChannel, &wg)

	// Отправка чисел в канал.
	for i := 1; i <= 5; i++ {
		numChannel <- i
	}
	close(numChannel) // Закрываем канал.
	wg.Wait()         // Ожидаем завершения горутины.

	// Задачи 2 и 3: Обработка изображений.
	fmt.Println("\nЗадачи 2 и 3: Обработка изображений")

	// Открытие входного файла.
	inputFile, err := os.Open("1.png")
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer inputFile.Close()

	// Декодирование изображения.
	img, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Println("Ошибка декодирования изображения:", err)
		return
	}

	bounds := img.Bounds()

	// Выполнение последовательной обработки.
	fmt.Println("Последовательная обработка:")
	sequentialImg := processSequential(img, bounds)

	// Сохранение результата последовательной обработки.
	outputFile1, err := os.Create("output_sequential.png")
	if err != nil {
		fmt.Println("Ошибка создания файла для последовательной обработки:", err)
		return
	}
	defer outputFile1.Close()
	err = png.Encode(outputFile1, sequentialImg)
	if err != nil {
		fmt.Println("Ошибка сохранения изображения после последовательной обработки:", err)
		return
	}

	// Выполнение параллельной обработки.
	fmt.Println("Параллельная обработка:")
	parallelImg := processParallel(img, bounds)

	// Сохранение результата параллельной обработки.
	outputFile2, err := os.Create("output_parallel.png")
	if err != nil {
		fmt.Println("Ошибка создания файла для параллельной обработки:", err)
		return
	}
	defer outputFile2.Close()
	err = png.Encode(outputFile2, parallelImg)
	if err != nil {
		fmt.Println("Ошибка сохранения изображения после параллельной обработки:", err)
		return
	}

	fmt.Println("Обработка завершена.")
	fmt.Println("Результаты сохранены в output_sequential.png и output_parallel.png")
}
