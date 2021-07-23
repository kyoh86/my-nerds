package main

import (
	"fmt"

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
	return usecase.WalkComics(server, serverRoot, func(path string, _ model.ComicType) error {
		fmt.Println(path)
		if listFlag.limit > 0 {
			count++
			if count >= listFlag.limit {
				return usecase.ErrStopWalking
			}
		}
		return nil
	})
}
