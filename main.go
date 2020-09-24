package main

import (
	"fmt"
	"github.com/go-rest-api/cmd"
)

const (
	appVersion = "v.0.0.1"
	banner     = "### Golang RESTful API %s ###"
)

func main() {
	fmt.Printf(banner, appVersion)
	cmd.Execute()
}
