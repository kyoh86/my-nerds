package main

import (
	"fmt"

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

func list(cmd *cobra.Command, args []string) error {
	count := 0
	return usecase.WalkComics(cmd.Context(), user, pass, func(path string, _ usecase.ComicType) error {
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
