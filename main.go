package main

import (
	"fmt"
	"os"

	"github.com/voidwyrm-2/jeffm/cmd"
)

func main() {
	if err := cmd.Execute("1.0"); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
