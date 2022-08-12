package commands

import (
	"github.com/okeyaki/dddoc/lib/adapter/console/renderers/graph"
	"github.com/okeyaki/dddoc/lib/platform/source"
	"github.com/spf13/cobra"
)

func NewVisualizeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "visualize",
		Short: "Visualize domain model",
		RunE:  runVisualizeCommand,
	}
}

func runVisualizeCommand(cmd *cobra.Command, args []string) error {
	cs, err := source.Analyze()
	if err != nil {
		return err
	}

	return graph.Render(cs)
}
