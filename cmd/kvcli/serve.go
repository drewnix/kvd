package kvcli

import (
	"context"
	"fmt"
	"github.com/drewnix/kvd/pkg/kvd"
	"github.com/spf13/cobra"
	"os"
)

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

//
//import (
//"context"
//"flag"
//"fmt"
//"os"
//"wineapi/src/wineapi"
//)
//
//const defaultPort = 4000
//
//func readFlags(args []string, config *wineapi.Config) error {
//	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
//	var wineDataCSV = flags.String("csv", "", "CSV file containing wine data to load")
//	var port = flags.Int("port", defaultPort, "Port to listen on")
//	var maxRecords = flags.Int("max-records", defaultRecordsAllowed, "Max number of wines to allow a client to pull at once")
//	var wineDBSQLite = flags.String("db", defaultWineDBFile, "DB file to write data to (if using SQLite)")
//
//	if err := flags.Parse(args[1:]); err != nil {
//		return err
//	}
//
//	config.MetricsTick = *tick
//	config.WineDataFileCSV = *wineDataCSV
//	config.WineDBFileSQLite = *wineDBSQLite
//	config.MaxRecords = *maxRecords
//	config.Port = *port
//
//	return nil
//}
//
//func main() {
//	var config = wineapi.Config{}
//	err := readFlags(os.Args, &config)
//	var wa = wineapi.WineAPI{}
//	wa.Init(&config)
//
//	if err != nil {
//		fmt.Println("Error reading flags: ", err)
//		os.Exit(1)
//	}
//
//	ctx, err := wa.StartService(context.Background())
//	if err != nil {
//		fmt.Println("Error starting service: ", err)
//		os.Exit(1)
//	}
//
//	<-ctx.Done()
//	fmt.Println("Shutting down Wine API service")
//}
