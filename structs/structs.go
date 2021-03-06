package main

import (
	"fmt"
)
type position struct { 
	x	float32
	y	float32
}

type badGuy struct {
	name	string
	health	int
	pos		position
}

func whereIsBadGuy(b badGuy) {
	x := b.pos.x
	y := b.pos.y
	fmt.Printf("Bad Guy is at: ( %v , %v )", x, y)
}

func main() {
	p := position{4, 2}

	fmt.Println(p.x)
	fmt.Println(p.y)

	b := badGuy{"Dr. Evil", 100, p}
	fmt.Println(b)
	whereIsBadGuy(b)
}