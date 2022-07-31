package kvcli

import (
	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/spf13/cobra"
)

func init() {
	var setCmd = &cobra.Command{
		Use:     "set",
		Aliases: []string{"s"},
		Short:   "Puts a set of keys from the KVD service",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			kvcli.SetKey(args[0], args[1])
		},
	}

	rootCmd.AddCommand(setCmd)
}
