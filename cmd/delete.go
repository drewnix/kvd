package kvcli

import (
	"fmt"
	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/spf13/cobra"
)

func DeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "Delete a set of keys in the KVD service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no keys provided to delete")
			}
			
			client := kvcli.NewClient(ServerAddress)
			
			if len(args) == 1 {
				// Single key delete
				err := client.Delete(args[0])
				if err != nil {
					return fmt.Errorf("could not delete key: %w", err)
				}
			} else {
				// Bulk delete
				err := client.BulkDelete(args)
				if err != nil {
					return fmt.Errorf("could not delete keys: %w", err)
				}
			}
			
			fmt.Println("Keys deleted")
			return nil
		},
	}
}

func init() {
	var deleteCmd = DeleteCmd()

	rootCmd.AddCommand(deleteCmd)
}
