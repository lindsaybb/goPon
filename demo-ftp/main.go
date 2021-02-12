package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"iskratel/gopon"
)

var (
	helpFlag       = flag.Bool("h", false, "Show this help")
	pathFlag       = flag.String("o", ".", "Path to output files to")
	uploadScrFlag  = flag.String("us", "", "Path of Olt Config (.scr) file to Upload")
	uploadConfFlag = flag.String("uc", "", "Path of Innbox Config (.conf) file to Upload")
	deleteFlag     = flag.Bool("nd", true, "Delete uploads by name")
	//sleepFlag = flag.Duration("s", time.Millisecond*0, "Specify breaks between actions")
)

const usage = "`gopon` [options] <olt_ip>"

var newConfigScript = "../ref/testDemo.scr"
var newInnboxConfig = "../ref/testDemo.conf"

// tail doesn't work
// need to apply this to the config and onu-default
// need to show Uploading files
// need to be able to get a list of files by name only
// added bonus would be to give some idea of what the logs contain other than just their name

func main() {
	flag.Parse()

	if *helpFlag || flag.NArg() < 1 {
		fmt.Println(usage)
		flag.PrintDefaults()
		return
	}
	outPath, err := filepath.Abs(*pathFlag)
	if err != nil {
		fmt.Println("Invalid -o flag supplied", err)
		return
	}
	if *uploadScrFlag != "" {
		newConfigScript = *uploadScrFlag
	}
	if *uploadConfFlag != "" {
		newInnboxConfig = *uploadConfFlag
	}

	start := time.Now()
	host := flag.Args()[0]

	olt := gopon.NewLumiaOlt(host)
	if !olt.HostIsReachable() {
		fmt.Println(gopon.ErrNotReachable)
		return
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	err = olt.GetCurrentLogs(outPath)
	//err = olt.GetCurrentLogsTail(outPath)
	fmt.Printf("Error getting current logs: %v\n", err)

	err = olt.UploadConfig(newConfigScript)
	fmt.Printf("Error uploading config script to OLT: %v\n", err)
	err = olt.UploadConfig(newInnboxConfig)
	fmt.Printf("Error uploading Innbox config to OLT: %v\n", err)
	if *deleteFlag {
		err = olt.DeleteConfig(newConfigScript)
		fmt.Printf("Error removing config script from OLT: %v\n", err)
		err = olt.DeleteConfig(newInnboxConfig)
		fmt.Printf("Error removing Innbox config from OLT: %v\n", err)
	}
	fmt.Printf("Demo complete in %v\n", time.Since(start))
}
