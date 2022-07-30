package kvcli

import (
	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/spf13/cobra"
)

func init() {
	var deleteCmd = &cobra.Command{
		Use:     "delete",
		Aliases: []string{"d"},
		Short:   "Deletes a set of keys from the KVD service",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			kvcli.DeleteKey(args[0])
		},
	}

	rootCmd.AddCommand(deleteCmd)
}
