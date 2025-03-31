package kvcli

import (
	"fmt"
	"strings"

	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/spf13/cobra"
)

func SetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set",
		Aliases: []string{"s"},
		Short:   "Sets key-value pairs in the KVD service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no key-value pairs provided to set")
			}
			
			client := kvcli.NewClient(ServerAddress)
			
			// Parse key=value pairs
			if len(args) == 1 && strings.Contains(args[0], "=") {
				// Single key-value pair
				parts := strings.SplitN(args[0], "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid key-value format, use key=value")
				}
				
				key, value := parts[0], parts[1]
				if key == "" {
					return fmt.Errorf("key cannot be empty")
				}
				
				err := client.Set(key, value)
				if err != nil {
					return fmt.Errorf("could not set key %s: %w", key, err)
				}
			} else {
				// Multiple key-value pairs
				kvPairs := make(map[string]string)
				
				for _, arg := range args {
					parts := strings.SplitN(arg, "=", 2)
					if len(parts) != 2 {
						return fmt.Errorf("invalid key-value format for %q, use key=value", arg)
					}
					
					key, value := parts[0], parts[1]
					if key == "" {
						return fmt.Errorf("key cannot be empty")
					}
					
					kvPairs[key] = value
				}
				
				err := client.BulkSet(kvPairs)
				if err != nil {
					return fmt.Errorf("could not set keys: %w", err)
				}
			}
			
			fmt.Println("Keys set")
			return nil
		},
	}
}

func init() {
	var setCmd = SetCmd()
	rootCmd.AddCommand(setCmd)
}
