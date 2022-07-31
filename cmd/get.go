package kvcli

import (
	"fmt"
	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/spf13/cobra"
)

func init() {
	var getCmd = &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Gets a set of keys from the KVD service",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			res := kvcli.GetKey(args[0])
			fmt.Println(res)
		},
	}

	rootCmd.AddCommand(getCmd)
}
