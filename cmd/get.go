package kvcli

import (
	"fmt"
	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/spf13/cobra"
	"os"
)

func GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Puts a get of keys from the KVD service",
		RunE: func(cmd *cobra.Command, args []string) error {
			var argsLen int = len(args)
			gets := make([]string, argsLen)

			for i, s := range args {
				gets[i] = s
			}
			fmt.Println(gets)
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
