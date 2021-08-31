package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	message := stringutil.Reverse("Hello, OTUS!")

	fmt.Println(message)
}
