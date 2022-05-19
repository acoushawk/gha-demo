package main

import "os"

func main() {
	d1 := []byte("hello\nworld\n")
	_ = os.WriteFile("output.text", d1, 0644)

}
