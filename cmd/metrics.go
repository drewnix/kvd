package kvcli

import (
	"fmt"

	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/spf13/cobra"
)

func init() {
	var metricsCmd = &cobra.Command{
		Use:     "metrics",
		Aliases: []string{"m"},
		Short:   "Gets metrics from the KVD service",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := kvcli.NewClient(ServerAddress)
			
			metrics, err := client.GetMetrics()
			if err != nil {
				return fmt.Errorf("could not get metrics: %w", err)
			}
			
			fmt.Println("Keys Stored:", metrics.KeysStored)
			fmt.Println("Bytes Stored (Values):", metrics.ValueBytesStored)
			fmt.Println("Set Operations:", metrics.SetOps)
			fmt.Println("Get Operations:", metrics.GetOps)
			fmt.Println("Delete Operations:", metrics.DelOps)
			
			return nil
		},
	}

	rootCmd.AddCommand(metricsCmd)
}
