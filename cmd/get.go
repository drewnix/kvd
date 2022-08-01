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
		Short:   "Puts a get of keys from the KVD service",
		RunE: func(cmd *cobra.Command, args []string) error {
			var argsLen int = len(args)
			gets := make([]string, argsLen)
			copy(gets, args)

			recs := kvcli.GetKeys(gets)
			for _, rec := range recs {
				_, err := fmt.Fprintf(os.Stdout, "%s: %s\n", rec.Key, rec.Value)
				if err != nil {
					return err
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
