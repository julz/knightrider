package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "knife",
	Long:  "Knife (pronounced kaynife) is a super simple program for working with knative yml",
	Short: `Knife is a super simple program for working with knative yml`,
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
