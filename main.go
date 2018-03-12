package main

import (
	"fmt"
	"os"
	"templater/cmd"
)

/**
 * Initialize Cobra and boot the app.
 */
func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
