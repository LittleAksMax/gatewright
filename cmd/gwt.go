package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/LittleAksMax/gatewright/internal/config"
)

func Execute() {
	cmd, err := createCommand(config.CreateConfig()).ExecuteC()
	if err != nil {
		if strings.HasPrefix(err.Error(), "unknown command ") {
			fmt.Fprintln(os.Stderr, cmd.UsageString())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
