package main

import (
	"fmt"
	"sync"
	"time"
)

// pipeline: stage reads ints, squares them, sends to out.
func squareStage(in <-chan int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for n := range in {
		out <- n * n
	}
}

func fanOut(input []int, workers int) <-chan int {
	in := make(chan int)
	out := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go squareStage(in, out, &wg)
	}

	go func() {
		for _, v := range input {
			in <- v
		}
		close(in)
		wg.Wait()
		close(out)
	}()

	return out
}

func selectWithTimeout(results <-chan int) {
	timeout := time.After(200 * time.Millisecond)
	for {
		select {
		case v, ok := <-results:
			if !ok {
				fmt.Println("results closed")
				return
			}
			fmt.Println("got", v)
		case <-timeout:
			fmt.Println("select timeout — stopping consumer")
			return
		}
	}
}

func main() {
	nums := []int{1, 2, 3, 4, 5}
	results := fanOut(nums, 2)
	selectWithTimeout(results)
}