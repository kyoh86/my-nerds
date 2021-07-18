package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:  "list",
	RunE: list,
}

var (
	listFlag struct {
		once int
	}
)

func init() {
	listCommand.Flags().IntVarP(&listFlag.once, "once", "", 100, "Limit to traverse")
	facadeCommand.AddCommand(listCommand)
}

func list(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	conn, err := connect(ctx)
	if err != nil {
		return err
	}

	done := map[string]struct{}{}
	walker := conn.Walk(root)
	for i := 0; (i < listFlag.once || listFlag.once < 1) && walker.Next(); i++ {
		if walker.Stat().Type != ftp.EntryTypeFile {
			continue
		}
		path := walker.Path()
		if strings.HasSuffix(path, ".AppleDouble") || strings.HasSuffix(filepath.Dir(path), ".AppleDouble") {
			walker.SkipDir()
			continue
		}
		lower := strings.ToLower(path)
		if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") || strings.HasSuffix(lower, "png") {
			path := filepath.Dir(path)
			if _, ok := done[path]; ok {
				continue
			}
			done[path] = struct{}{}
			fmt.Println(path)
			walker.SkipDir()
		} else if strings.HasSuffix(lower, ".rar") {
			if _, ok := done[path]; ok {
				continue
			}
			done[path] = struct{}{}
			fmt.Println(path)
		}
	}
	if err := walker.Err(); err != nil {
		return fmt.Errorf("walk on the server: %w", err)
	}
	return nil
}
