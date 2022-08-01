package kvcli

import (
	"fmt"
	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/spf13/cobra"
	"os"
)

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "Delete a set of keys in the KVD service",
		RunE: func(cmd *cobra.Command, args []string) error {
			var argsLen int = len(args)
			dels := make([]string, argsLen)

			for i, s := range args {
				dels[i] = s
			}
			err := kvcli.DeleteKeys(dels)
			if err != nil {
				fmt.Print("Could not delete keys: ", err)
				os.Exit(1)
			}
			return nil
		},
	}
}

func init() {
	var deleteCmd = DeleteCmd()

	rootCmd.AddCommand(deleteCmd)
}
