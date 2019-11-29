package main

import (
	"fmt"
	"os"

	"github.com/p1ass/mikku"
)

func main() {
	if err := mikku.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
