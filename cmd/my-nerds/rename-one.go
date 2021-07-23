package main

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/kyoh86/my-nerds/usecase"
	"github.com/spf13/cobra"
)

var renameOneCommand = &cobra.Command{
	Use:  "rename-one",
	Args: cobra.ExactArgs(1),
	RunE: renameOne,
}

var (
	renameOneFlag struct {
		skipInvalid bool
		forceNoDiff bool
		force       bool
	}
)

func init() {
	renameOneCommand.Flags().BoolVarP(&renameOneFlag.forceNoDiff, "force-no-diff", "", false, "Rename forcely if there's no diff(without confirmation)")
	renameOneCommand.Flags().BoolVarP(&renameOneFlag.skipInvalid, "skip-invalid", "", false, "Skip invalid files")
	renameOneCommand.Flags().BoolVarP(&renameOneFlag.force, "force", "", false, "Rename forcely (without confirmation)")
	facadeCommand.AddCommand(renameOneCommand)
}

func renameOne(cmd *cobra.Command, args []string) error {
	server, err := openServer(cmd.Context())
	if err != nil {
		return err
	}
	err = usecase.RenameComic(server, args[0], usecase.RenameComicOption{
		SkipInvalid: renameOneFlag.skipInvalid,
		ForceNoDiff: renameOneFlag.forceNoDiff,
		Force:       renameOneFlag.force,
	})
	switch {
	case errors.Is(err, terminal.InterruptErr):
		fmt.Println("interrupted")
		return nil
	default:
		return err
	}
}
