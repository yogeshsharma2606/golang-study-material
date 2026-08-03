package main

import "fmt"

func sliceAliasing() {
	original := []int{1, 2, 3}
	alias := original
	alias[0] = 99
	fmt.Println("aliasing:", original, alias)
}

func appendGotcha() {
	a := []int{1, 2, 3}
	b := append(a, 4)
	b[0] = 100
	fmt.Println("after append alias write:", a, b)

	// Capacity sharing: sometimes append reuses backing array
	c := make([]int, 3, 6)
	copy(c, []int{10, 20, 30})
	d := append(c, 40)
	d[1] = 999
	fmt.Println("shared backing array?", c, d)
}

func mapIteration() {
	scores := map[string]int{"alice": 10, "bob": 8}
	for name, score := range scores {
		fmt.Printf("%s: %d\n", name, score)
	}
	// Map keys are not ordered; run multiple times to see order change.
}

func main() {
	fmt.Println("=== slice aliasing ===")
	sliceAliasing()
	fmt.Println("=== append gotchas ===")
	appendGotcha()
	fmt.Println("=== map iteration ===")
	mapIteration()
}