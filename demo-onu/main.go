package main

import (
	"flag"
	"fmt"
	"time"

	"iskratel/gopon"
)

var (
	helpFlag = flag.Bool("h", false, "Show this help")
	//sleepFlag = flag.Duration("s", time.Millisecond*0, "Specify breaks between actions")
)

const usage = "`gopon` [options] <olt_ip>"

func main() {
	flag.Parse()

	if *helpFlag || flag.NArg() < 1 {
		fmt.Println(usage)
		flag.PrintDefaults()
		return
	}
	var err error
	start := time.Now()
	host := flag.Args()[0]

	//err = onuRegistrationDemo(host)
	//fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
	err = onuProfileUsageDemo(host)
	fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
}

func onuRegistrationDemo(host string) error {
	var err error
	olt := gopon.NewLumiaOlt(host)
	if !olt.HostIsReachable() {
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}
	// populate the Authorized Onu list with local file
	err = olt.LoadOnuAuthList("authList.txt")
	if err != nil {
		return err
	}
	err = olt.UpdateRegisteredOnuList()
	if err != nil {
		return err
	}

	var obll *gopon.OnuBlacklistList
	obll, err = olt.GetOnuBlacklist()
	if err != nil {
		return err
	}
	obll.Tabwrite()

	fmt.Println("\nCreating an Onu Config for each Onu on the Blacklist to attempt to Register")

	var ocfg *gopon.OnuConfig
	for _, e := range obll.Entry {
		intf := olt.NextAvailableOnuInterface(e.IfName)
		//uintf := gopon.UrlEncodeInterface(intf)
		ocfg = gopon.NewOnuConfig(e.SerialNumber, intf)
		err = olt.AuthorizeOnu(ocfg)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			break // only do one at a time for the first try
		}
	}

	nobll, err := olt.GetOnuBlacklist()
	if err != nil {
		return err
	}
	nobll.Tabwrite()

	return nil
}

func onuProfileUsageDemo(host string) error {
	var err error
	olt := gopon.NewLumiaOlt(host)
	if !olt.HostIsReachable() {
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	var opl *gopon.OnuProfileList
	opl, err = olt.GetOnuProfileUsage()
	if err != nil {
		return err
	}
	opl.Tabwrite()

	var spl *gopon.ServiceProfileList
	spl, err = olt.GetServiceProfiles()
	if err != nil {
		return err
	}
	spl.Tabwrite()

	var newProfile *gopon.OnuProfile
	// as an example, take the 2nd entry from the ServiceProfileList and Apply it to the 2nd item in the OnuProfileList
	onuIntf := opl.Entry[2].IfName
	onuSpName := spl.Entry[2].Name
	// this will not error as long as both are valid strings
	newProfile = gopon.NewOnuProfile(onuIntf, onuSpName)
	// error will be checked if the server accepts the patch
	err = olt.PostOnuProfile(newProfile)
	if err != nil {
		return err
	}
	// we will get a new onu profile list to check that this profile was added
	nopl, err := olt.GetOnuProfileUsage()
	if err != nil {
		return err
	}
	nopl.Tabwrite()

	// we can also remove profiles from devices by intf and service profile name
	err = olt.RemoveOnuProfileUsage(onuIntf, onuSpName)
	//err = olt.RemoveOnuProfileUsage("0/5/3", "102_DATA_Acc")
	if err != nil {
		return err
	}

	// and get the list to make sure it is gone
	nnopl, err := olt.GetOnuProfileUsage()
	if err != nil {
		return err
	}
	nnopl.Tabwrite()
	//fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(opl.Entry), len(nopl.Entry), len(nnopl.Entry))

	return nil
}
