package main

import (
	"fmt"
	"sort"

	"github.com/guptarohit/asciigraph"
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
	comic, err := usecase.AnalyzeComic(src, args[0])
	if err != nil {
		return err
	}
	sort.Slice(comic.Pages, func(i, j int) bool {
		return comic.Pages[i].Rel < comic.Pages[j].Rel
	})
	for _, p := range comic.Pages {
		fmt.Printf("%s (%s, %dx%d (%f), %s)\n", p.Rel, p.Format, p.Config.Width, p.Config.Height, p.AspectRatio, model.ColorModelToString(p.Config.ColorModel))
	}
	fmt.Println("aspectRatio:")
	fmt.Printf("      avg: %f\n", comic.AspectRatioAverage)
	fmt.Printf("      var: %f\n", comic.AspectRatioVariance)
	fmt.Printf("    stdev: %f\n", comic.AspectRatioStandardDeviation)
	{
		fmt.Println("hue histogram:")
		graph := asciigraph.Plot(comic.HueHistogram50[:], asciigraph.Width(200))
		fmt.Println(graph)
		fmt.Print("     ")
		for _, v := range comic.HueHistogram50 {
			fmt.Printf("%.1f ", v)
		}
		fmt.Println()
	}
	{
		fmt.Println("chroma histogram:")
		graph := asciigraph.Plot(comic.ChromaHistogram50[:], asciigraph.Width(200))
		fmt.Println(graph)
		fmt.Print("     ")
		for _, v := range comic.ChromaHistogram50 {
			fmt.Printf("%.1f ", v)
		}
		fmt.Println()
	}
	{
		fmt.Println("luminance histogram:")
		graph := asciigraph.Plot(comic.LuminanceHistogram50[:], asciigraph.Width(200))
		fmt.Println(graph)
		fmt.Print("     ")
		for _, v := range comic.LuminanceHistogram50 {
			fmt.Printf("%.1f ", v)
		}
		fmt.Println()
	}
	return nil
}
