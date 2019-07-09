package main

import (
	"fmt"
	"os"

	"github.com/RobinBaeckman/ragnar/pkg/http/rest"
)

func main() {
	if err := rest.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
