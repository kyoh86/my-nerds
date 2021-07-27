package main

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/kyoh86/my-nerds/app"
	"github.com/kyoh86/my-nerds/driver/archive"
	"github.com/kyoh86/my-nerds/usecase"
	"github.com/spf13/cobra"
)

var getCommand = &cobra.Command{
	Use:  "get",
	Args: cobra.ExactArgs(1),
	RunE: get,
}

func init() {
	facadeCommand.AddCommand(getCommand)
}

func get(cmd *cobra.Command, args []string) error {
	server, err := openServer(cmd.Context())
	if err != nil {
		return err
	}

	pathFrom := args[0]
	name := path.Base(pathFrom)
	ext := filepath.Ext(name)
	if strings.ToLower(ext) == ".rar" {
		pathTo := filepath.Join(app.LocalRoot, strings.TrimSuffix(name, ext))
		return usecase.ExtractComic(server, archive.ExtractRAR, pathFrom, pathTo)
	} else {
		pathTo := filepath.Join(app.LocalRoot, name)
		return usecase.DownloadComicDir(server, pathFrom, pathTo)
	}
}
