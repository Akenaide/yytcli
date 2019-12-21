package main

import (
	"github.com/Akenaide/yytcli/cmd"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start().Stop()

	cmd.Execute()
}
