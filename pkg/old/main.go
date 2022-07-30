package old

//
//import (
//	"context"
//	"flag"
//	"fmt"
//	"os"
//	"wineapi/src/wineapi"
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
