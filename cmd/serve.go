package kvcli

import (
	"context"
	"fmt"
	"os"

	"github.com/drewnix/kvd/pkg/kvd"
	"github.com/spf13/cobra"
)

func ServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "serve",
		Aliases: []string{"srv"},
		Short:   "Run the KVD Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			startService()
			return nil
		},
	}
}

func startService() {
	//res := kvd.StartService
	var svc = kvd.Kvd{}
	const defaultPort = 4000
	var config = kvd.Config{}

	config.Port = defaultPort

	svc.Init(&config)

	ctx, err := svc.StartService(context.Background())
	if err != nil {
		fmt.Println("Error starting service: ", err)
		os.Exit(1)
	}

	<-ctx.Done()
	fmt.Println("Shutting down KVD service")
}

func init() {
	var daemon bool
	var serveCmd = &cobra.Command{
		Use:     "serve",
		Aliases: []string{"srv"},
		Short:   "Runs the KVD service",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if daemon {
				fmt.Println("gonne start", daemon)
			}
			startService()
		},
	}
	serveCmd.Flags().BoolVarP(&daemon, "deamon", "d", false, "is daemon?")

	rootCmd.AddCommand(serveCmd)
}
