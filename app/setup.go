// Package app preforms basic app functions like setup, loading config and global definitions
package app

import (
	"fmt"
)

// Setup will preform app setup or panic of it fails
// Can only be called once
func Setup() {
	// Logger setup
	{
		err := devLoggerSetup()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			panic("faled to setup logger")
		}
	}
}
