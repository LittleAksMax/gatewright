package cmd

import (
	"fmt"
	"os"

	"github.com/LittleAksMax/gatewright/internal/config"
)

func Execute() {
	if err := createCommand(config.CreateConfig()).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
