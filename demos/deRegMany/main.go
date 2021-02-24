package main

import (
	"flag"
	"fmt"

	"github.com/lindsaybb/gopon"
)

var (
	helpFlag   = flag.Bool("h", false, "Show this help")
	deAuthFile = flag.String("df", "deAuthList.txt", "Path to file that contains list of ONU Serial Numbers to Deauthorize")
)

// this demo is for removing a customer from service permanently, moving or changing their ONU
// for temporarily removing service to the ONU this is not the most efficient method

func main() {
	flag.Parse()

	if *helpFlag || flag.NArg() < 1 {
		fmt.Println("Supply OLT IP Address as Arg")
		flag.PrintDefaults()
		return
	}
	host := flag.Args()[0]
	olt := gopon.NewLumiaOlt(host)
	if !olt.HostIsReachable() {
		fmt.Printf("Host %s is not reachable\n", host)
		return
	}
	err := deRegisterOnuFromFile(olt)
	fmt.Println(err)
}

func deRegisterOnuFromFile(olt *gopon.LumiaOlt) error {
	var err error
	// this function updates the list of currently connected devices with two GET requests
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	if len(olt.Registration) < 1 {
		return gopon.ErrNotExists
	}
	olt.TabwriteRegistry()
	// Deauthorized any Onu Serial Numbers that appear in the file
	err = olt.DeauthOnuBySnList(*deAuthFile)
	if err != nil {
		return err
	}

	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	olt.TabwriteRegistry()
	// Blacklist might not show the devices yet, takes about 1 min to update

	return nil
}
