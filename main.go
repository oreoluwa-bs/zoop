package main

import (
	"github.com/oreoluwa-bs/zoop/cmd"
)

var version = "dev"

func main() {
	cmd.Version = version
	cmd.Execute()
}
