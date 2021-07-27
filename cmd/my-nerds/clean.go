package main

import (
	"github.com/kyoh86/my-nerds/app"
	"github.com/kyoh86/my-nerds/usecase"
	"github.com/spf13/cobra"
)

var cleanCommand = &cobra.Command{
	Use:  "clean",
	RunE: clean,
}

func init() {
	facadeCommand.AddCommand(cleanCommand)
}

func clean(cmd *cobra.Command, args []string) error {
	return usecase.RemoveAll(cmd.Context(), app.LocalRoot)
}
