package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "kr",
	Long:  "kr  is a super simple program for working with knative yml",
	Short: `kr is a super simple program for working with knative yml`,
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
