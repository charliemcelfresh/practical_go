package main

import (
	"charliemcelfresh/practical_go/cmd"
	"fmt"
)

func init() {
	fmt.Println("Running main config")
}

func main() {
	cmd.Execute()
}
