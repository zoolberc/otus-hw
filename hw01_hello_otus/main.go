package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	reversWord := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(reversWord)
}
