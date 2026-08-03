package main

import "fmt"

func describe(n int, p *int) {
	fmt.Printf("value=%d pointer=%p deref=%d\n", n, p, *p)
}

func sum(nums ...int) int {
	total := 0
	for _, v := range nums {
		total += v
	}
	return total
}

func deferDemo() {
	defer fmt.Println("defer runs when deferDemo returns (LIFO order)")
	defer fmt.Println("second defer")
	fmt.Println("inside deferDemo")
}

func main() {
	// Types and zero values
	var count int
	var label string
	fmt.Println("zero values:", count, label)

	// Pointers
	x := 42
	ptr := &x
	*ptr = 7
	fmt.Println("after pointer write, x =", x)
	describe(x, ptr)

	// Variadic functions
	fmt.Println("variadic sum:", sum(1, 2, 3, 4))

	deferDemo()
}