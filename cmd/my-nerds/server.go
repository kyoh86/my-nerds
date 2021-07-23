package main

import (
	"context"

	"github.com/kyoh86/my-nerds/driver/source"
)

func openServer(ctx context.Context) (*source.FTPServer, error) {
	return source.OpenFTPServer(ctx, host, user, pass)
}
