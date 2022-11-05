package main

import (
	parse "GolangProjects/Parse"
	"fmt"
)

func main() {
	parse.Start()
	fmt.Println(len(parse.Groups))
}
