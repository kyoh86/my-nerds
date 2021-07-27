package usecase

import (
	"context"
	"os"
)

func RemoveAll(ctx context.Context, dir string) error {
	return os.RemoveAll(dir)
}
