package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/kyoh86/my-nerds/model"
	"github.com/kyoh86/my-nerds/usecase"
	"github.com/spf13/cobra"
)

var renameCommand = &cobra.Command{
	Use:  "rename",
	RunE: rename,
}

var (
	renameFlag struct {
		replaceMapFile string
		skipInvalid    bool
		forceNoDiff    bool
		dryrun         bool
	}
)

func init() {
	renameCommand.Flags().StringVarP(&renameFlag.replaceMapFile, "replace-map-file", "", "", "A file holds replace regexp map")
	renameCommand.Flags().BoolVarP(&renameFlag.forceNoDiff, "force-no-diff", "", false, "Rename forcely if there's no diff(without confirmation)")
	renameCommand.Flags().BoolVarP(&renameFlag.skipInvalid, "skip-invalid", "", false, "Skip invalid files")
	renameCommand.Flags().BoolVarP(&renameFlag.dryrun, "dryrun", "", false, "Dry run")
	facadeCommand.AddCommand(renameCommand)
}

func rename(cmd *cobra.Command, _ []string) error {
	server, err := openServer(cmd.Context())
	if err != nil {
		return err
	}
	m, err := loadReplaceMap(renameFlag.replaceMapFile)
	if err != nil {
		return err
	}
	return usecase.WalkComics(server, serverRoot, func(path string, _ model.ComicType) error {
		err := usecase.RenameComic(server, path, usecase.RenameComicOption{
			ReplaceMap:  m,
			ForceNoDiff: renameFlag.forceNoDiff,
			SkipInvalid: renameFlag.skipInvalid,
			Dryrun:      renameFlag.dryrun,
		})
		switch {
		case errors.Is(err, terminal.InterruptErr):
			fmt.Println("interrupted")
			return usecase.ErrStopWalking
		case err == nil:
			return nil
		default:
			fmt.Println(err)
			return nil
		}
	})
}

func loadReplaceMap(filename string) (map[*regexp.Regexp]string, error) {
	if filename == "" {
		return nil, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open replace map file: %w", err)
	}
	defer file.Close()

	m := map[*regexp.Regexp]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		pair := strings.Split(line, "@")
		pat, rep := pair[0], pair[1]
		reg, err := regexp.Compile(pat)
		if err != nil {
			return nil, err
		}
		m[reg] = rep
	}
	return m, nil
}
