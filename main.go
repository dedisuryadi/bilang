package main

import (
	"os"

	"github.com/dedisuryadi/bilang/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
