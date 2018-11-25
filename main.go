package main

import (
	"bytes"
	"fmt"
)

func main() {
	hello := "hello"
	helloBytes := []byte(hello)
	fmt.Println(helloBytes)

	world := "world"
	worldBytes := []byte(world)
	fmt.Println(worldBytes)

	helloWord := [][]byte{helloBytes, worldBytes}
	fmt.Println(helloWord)

	helloWords := bytes.Join(helloWord, []byte{})
	fmt.Println(helloWords)
}
