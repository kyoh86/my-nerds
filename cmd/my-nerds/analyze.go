package main

import (
	"fmt"
	"sort"

	"github.com/kyoh86/my-nerds/driver/source"
	"github.com/kyoh86/my-nerds/model"
	"github.com/kyoh86/my-nerds/usecase"
	"github.com/spf13/cobra"
)

var analyzeCommand = &cobra.Command{
	Use:  "analyze",
	Args: cobra.ExactArgs(1),
	RunE: analyze,
}

var analyzeFlag struct {
	local bool
}

func init() {
	analyzeCommand.Flags().BoolVarP(&analyzeFlag.local, "local", "l", false, "Analyze local dir")
	facadeCommand.AddCommand(analyzeCommand)
}

func analyze(cmd *cobra.Command, args []string) error {
	var src source.Source
	if analyzeFlag.local {
		local, err := source.OpenLocal()
		if err != nil {
			return err
		}
		src = local
	} else {
		server, err := openServer(cmd.Context())
		if err != nil {
			return err
		}
		src = server
	}
	pages, err := usecase.AnalyzeComic(src, args[0])
	if err != nil {
		return err
	}
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Rel < pages[j].Rel
	})
	for _, p := range pages {
		fmt.Printf("%s (%s, %dx%d, %s)\n", p.Rel, p.Format, p.Config.Width, p.Config.Height, model.ColorModelToString(p.Config.ColorModel))
	}
	return nil
}
