package main

import (
	"go-githubhost/runner"
)

func main() {
	err := runner.Run()
	if err != nil {
		panic(err)
	}
}
