package main

import (
	"fmt"
	"os"

	"github.com/voidwyrm-2/jeffm/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
