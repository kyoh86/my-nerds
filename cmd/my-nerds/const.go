package main

import "os"

const (
	host       = "192.168.11.12:21"
	user       = "kyoh86"
	serverRoot = "/sataraid1/nerd"
	localRoot  = "/home/kyoh86/Downloads/nerds"
)

var (
	pass = os.Getenv("MY_NERDS_PASS")
)
