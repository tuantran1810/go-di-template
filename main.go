package main

import (
	"fmt"

	"github.com/tuantran1810/go-di-template/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(fmt.Sprintf("Failed to execute command: %v", err))
	}
}
