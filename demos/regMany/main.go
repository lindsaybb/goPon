package main

import (
	"flag"
	"fmt"

	//"time"

	"iskratel/gopon"
)

var (
	helpFlag = flag.Bool("h", false, "Show this help")
	authFile = flag.String("af", "authList.txt", "Path to file that contains list of Authorized ONU Serial Numbers")
)

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
	err := registerOnuFromFile(olt)
	fmt.Println(err)
}

func registerOnuFromFile(olt *gopon.LumiaOlt) error {
	var err error
	// this function updates the list of currently connected devices with two GET requests
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	olt.TabwriteRegistry()
	// this function creates new registration entries from file
	err = olt.LoadOnuAuthList(*authFile)
	if err != nil {
		return err
	}
	//
	//	There are two paths that can be taken, pre-register does not have OLT port
	//	But Blacklist register requires multiple devices to be sitting waiting to register
	//	If the wrong OLT Port is supplied as pre-config, OLT will alarm when device comes up
	//	How to get the auth file to include the OLT port as optional?
	//	Can this operation check my database when a blacklist entry appears?
	//	Focus on current need: multiple ONU sit on the Blacklist, want to auth them all and add SP according to file

	// this function performs a GET request to retrieve the ONU currently on the Blacklist
	var obll *gopon.OnuBlacklistList
	obll, err = olt.GetOnuBlacklist()
	if err != nil {
		return err
	}
	obll.Tabwrite()

	fmt.Println("\nCreating an Onu Config for each Onu on the Blacklist to attempt to Register")

	var ocfg *gopon.OnuConfig
	for _, e := range obll.Entry {
		// since the Olt Registry has the new SerialNumbers, this check will include them
		if olt.ValidateSn(e.SerialNumber) {
			onuReg, err := olt.GetOnuRegisterBySn(e.SerialNumber)
			if err != nil {
				fmt.Println(err)
				continue
			}
			var intf string
			if onuReg.Interface == "" {
				onuReg = olt.NextAvailableOnuInterfaceUpdateRegister(e.IfName, onuReg)
			}
			intf = onuReg.Interface
			ocfg = gopon.NewOnuConfig(e.SerialNumber, intf)
			err = olt.AuthorizeOnu(ocfg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if len(onuReg.Services) < 1 {
				continue
			}
			spl, err := olt.GetServiceProfiles()
			if err != nil {
				return err
			}
			for _, sp := range onuReg.Services {
				if spl.ProfileExists(sp) {
					// created new ONU Profile by combining registered interface with service profile name
					newProfile := gopon.NewOnuProfile(ocfg.IfName, sp)
					// error will be checked if the server accepts the patch
					err = olt.PostOnuProfile(newProfile)
					if err != nil {
						return err
					}
				}
			}
		} else {
			fmt.Printf("Serial Number [%s] is on Blacklist but not on Authorized List. Skipping Registration.\n", e.SerialNumber)
		}
	}

	nobll, err := olt.GetOnuBlacklist()
	if err != nil {
		return err
	}
	nobll.Tabwrite()
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	olt.TabwriteRegistry()

	return nil
}
