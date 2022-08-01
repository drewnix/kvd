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
		Run: func(cmd *cobra.Command, args []string) {
			res := kvcli.GetMetrics()
			fmt.Println("Keys Stored:", res.KeysStored)
			fmt.Println("Set Operations:", res.SetOps)
			fmt.Println("Get Operations:", res.GetOps)
			fmt.Println("Delete Operations:", res.DelOps)
		},
	}

	rootCmd.AddCommand(metricsCmd)
}
