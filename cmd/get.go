package kvcli

import (
	"fmt"
	"os"

	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/spf13/cobra"
)

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Gets values for keys from the KVD service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no keys provided to get")
			}
			
			client := kvcli.NewClient(ServerAddress)
			
			if len(args) == 1 {
				// Single key get
				value, err := client.Get(args[0])
				if err != nil {
					return fmt.Errorf("could not get key %s: %w", args[0], err)
				}
				_, err = fmt.Fprintf(os.Stdout, "%s: %s\n", args[0], value)
				if err != nil {
					return err
				}
			} else {
				// Bulk get
				results, err := client.BulkGet(args)
				if err != nil {
					return fmt.Errorf("could not get keys: %w", err)
				}
				
				// Print in order of requested keys
				for _, key := range args {
					if value, ok := results[key]; ok {
						_, err := fmt.Fprintf(os.Stdout, "%s: %s\n", key, value)
						if err != nil {
							return err
						}
					}
				}
			}
			
			return nil
		},
	}
}

func init() {
	var getCmd = GetCmd()

	rootCmd.AddCommand(getCmd)
}
