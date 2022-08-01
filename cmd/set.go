package kvcli

import (
	"fmt"
	"strings"

	"github.com/drewnix/kvd/pkg/kvcli"
	"github.com/drewnix/kvd/pkg/kvd"
	"github.com/spf13/cobra"
)

func SetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set",
		Aliases: []string{"s"},
		Short:   "Puts a set of keys from the KVD service",
		RunE: func(cmd *cobra.Command, args []string) error {
			var argsLen int = len(args)
			sets := make([]kvd.Record, argsLen)

			for i, s := range args {
				a := strings.Split(s, "=")

				var rec = kvd.Record{
					Key:   a[0],
					Value: a[1],
				}

				sets[i] = rec
			}
			kvcli.SetKeys(sets)
			fmt.Println("Keys set")
			return nil
		},
	}
}

func init() {
	var setCmd = SetCmd()
	rootCmd.AddCommand(setCmd)
}
