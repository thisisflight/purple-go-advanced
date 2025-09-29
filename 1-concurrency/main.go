package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	sliceSize := 10
	numCh := make(chan int)
	squareNumCh := make(chan int)

	wg.Add(1)
	go func() {
		generateSeries(sliceSize, numCh)
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		generateSquareNums(numCh, squareNumCh)
		defer wg.Done()
	}()

	go func() {
		wg.Wait()

	}()

	squareSlice := make([]string, 0, 10)

	for num := range squareNumCh {
		squareSlice = append(squareSlice, fmt.Sprintf("%d", num))
	}
	fmt.Println(strings.Join(squareSlice, " "))
}

func generateSeries(sliceSize int, numCh chan int) {
	for range sliceSize {
		numCh <- rand.Intn(100)
	}
	close(numCh)
}

func generateSquareNums(numCh chan int, squareNumCh chan int) {
	for num := range numCh {
		squareNumCh <- num * num
	}
	close(squareNumCh)
}
