package main

import (
	"fmt"
	"os"

	"github.com/okeyaki/dddoc/lib/adapter/console/commands"
	"github.com/okeyaki/dddoc/lib/platform"
	"github.com/spf13/cobra"
)

func main() {
	cobra.OnInitialize(func() {
		if err := platform.LoadConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "%+v\n", err)

			os.Exit(1)
		}
	})

	cmd := &cobra.Command{
		Use: "dddoc",
	}
	cmd.AddCommand(commands.NewVisualizeCommand())

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)

		os.Exit(1)
	}
}
