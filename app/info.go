package app

import "os"

const (
	Name       = "my-nerds"
	Host       = "192.168.11.12:21"
	User       = "kyoh86"
	ServerRoot = "/sataraid1/nerd"
	LocalRoot  = "/home/kyoh86/Downloads/nerds"
)

var (
	Pass = os.Getenv("MY_NERDS_PASS")
)
