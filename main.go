package main

import (
	"fmt"
	"go-githubhost/runner"
)

func main() {
	text, err := runner.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
