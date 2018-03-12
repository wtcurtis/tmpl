package main

import (
	"fmt"
	"os"
	"github.com/wtcurtis/tmpl/cmd"
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
