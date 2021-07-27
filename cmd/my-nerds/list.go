package main

import (
	"fmt"

	"github.com/kyoh86/my-nerds/app"
	"github.com/kyoh86/my-nerds/model"
	"github.com/kyoh86/my-nerds/usecase"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:  "list",
	RunE: list,
}

var (
	listFlag struct {
		limit int
	}
)

func init() {
	listCommand.Flags().IntVarP(&listFlag.limit, "limit", "", 0, "Limit count")
	facadeCommand.AddCommand(listCommand)
}

func list(cmd *cobra.Command, _ []string) error {
	count := 0
	server, err := openServer(cmd.Context())
	if err != nil {
		return err
	}
	return usecase.WalkComics(server, app.ServerRoot, func(comic model.Comic) error {
		fmt.Println(comic.Path)
		if listFlag.limit > 0 {
			count++
			if count >= listFlag.limit {
				return usecase.ErrStopWalking
			}
		}
		return nil
	})
}
