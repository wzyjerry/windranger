package main

import (
	"github.com/spf13/cobra"
	"github.com/wzyjerry/windranger/internal/command/gogo"
)

func main() {
	cmd := &cobra.Command{
		Use: "windranger",
	}
	cmd.AddCommand(
		gogo.Gogo(),
	)
	_ = cmd.Execute()
}
