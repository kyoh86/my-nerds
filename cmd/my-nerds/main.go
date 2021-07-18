package main

import (
	"context"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/kyoh86/my-nerds/app"
	"github.com/spf13/cobra"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

var facadeCommand = &cobra.Command{
	Use:     app.Name,
	Short:   "Download scanned book archives from NAS and convert them to epub",
	Version: fmt.Sprintf("%s-%s (%s)", version, commit, date),
}

var (
	user string
	pass string
)

const (
	host = "192.168.11.12:21"
	root = "/sataraid1/nerd"
)

func init() {
	facadeCommand.PersistentFlags().StringVarP(&user, "user", "", "", "A user name to connect for the FTP server")
	facadeCommand.PersistentFlags().StringVarP(&pass, "pass", "", "", "A password to connect for the FTP server")
}

func main() {
	ctx := log.NewContext(context.Background(), &log.Logger{
		Handler: cli.New(os.Stderr),
		Level:   log.InfoLevel,
	})
	if err := facadeCommand.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
