package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jlaffaye/ftp"
)

func connect(ctx context.Context) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(host, ftp.DialWithTimeout(10*time.Second), ftp.DialWithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("dial to server: %w", err)
	}
	if err := conn.Login(user, pass); err != nil {
		return nil, fmt.Errorf("login on server: %w", err)
	}
	return conn, err
}
