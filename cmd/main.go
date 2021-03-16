package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
	"strings"

	"github.com/lindsaybb/gopon"
)

var (
	helpFlag       = flag.Bool("h", false, "Show this help")
	getBlacklist   = flag.Bool("gb", false, "Show the current ONU Blacklist on the OLT")
	getWhitelist   = flag.Bool("gw", false, "Show the current ONU Whitelist on the OLT")
	registerOne    = flag.Bool("ro", false, "Manually register one ONU to the OLT")
	registerMany   = flag.Bool("rm", false, "Automatically Register many ONU to the system using a file [af]; SN must exist on Blacklist")
	authFile       = flag.String("af", "authList.txt", "Path to file that contains list of Authorized ONU Serial Numbers and their Service Profiles")
	deregisterOne  = flag.Bool("do", false, "Manually Deregister one ONU from the OLT")
	deregisterMany = flag.Bool("dm", false, "Automatically Deregister multiple ONU from the system using a file [df]")
	deAuthFile     = flag.String("df", "deAuthList.txt", "Path to file that contains list of ONU Serial Numbers to Deauthorize")
	addOneSp       = flag.Bool("ap", false, "Manually Add one service profile to a registered ONU from a list of created Service Profiles")
	remOneSp       = flag.Bool("rp", false, "Manually Remove one service profile from a registered ONU")
	showSpDetails  = flag.Bool("sp", false, "View Detailed Information about Service Profiles")
)

// to add: deregister all from a specified port, authorize all from blacklist, indirect add/rem of service profiles in bulk

const usage = "`gopon_cmd` [options] <olt_ip>"

func main() {
	flag.Parse()

	if *helpFlag || flag.NArg() < 1 {
		fmt.Println(usage)
		flag.PrintDefaults()
		return
	}
	var err error
	host := flag.Args()[0]
	olt := gopon.NewLumiaOlt(host)
	if !olt.HostIsReachable() {
		fmt.Printf("Host %s is not reachable\n", host)
		return
	}
	if *getBlacklist {
		fmt.Println(">> Get Blacklist Called [-gb]")
		var obll *gopon.OnuBlacklistList
		obll, err = olt.GetOnuBlacklist()
		if err != nil {
			fmt.Println(err)
			return
		}
		obll.Tabwrite()
		promptContinue()
	}
	if *getWhitelist {
		fmt.Println(">> Get Whitelist Called [-gw]")
		err = olt.UpdateOnuRegistry()
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(olt.Registration) < 1 {
			fmt.Println("No Registered ONU on OLT!")
		} else {
			olt.TabwriteRegistry()
		}
		promptContinue()
	}
	if *registerOne {
		fmt.Println(">> Register One Called [-ro]")
		err = manuallyRegisterOnu(olt)
		if err != nil {
			fmt.Printf("Error running demo: %v\n", err)
		}
		promptContinue()
	}
	if *deregisterOne {
		fmt.Println(">> Deregister One Called [-do]")
		err = manuallyDeregisterOnu(olt)
		if err != nil {
			fmt.Printf("Error running demo: %v\n", err)
		}
		promptContinue()
	}
	if *registerMany {
		fmt.Println(">> Register Many Called [-rm]")
		err = registerOnuFromFile(olt)
		if err != nil {
			fmt.Printf("Error running demo: %v\n", err)
		}
		promptContinue()
	}
	if *deregisterMany {
		fmt.Println(">> Deregister Many Called [-dm]")
		err = deRegisterOnuFromFile(olt)
		if err != nil {
			fmt.Printf("Error running demo: %v\n", err)
		}
		promptContinue()
	}
	if *addOneSp {
		fmt.Println(">> Add One Service Profile Called [-ap]")
		err = addServiceToOnu(olt)
		if err != nil {
			fmt.Printf("Error running demo: %v\n", err)
		}
		promptContinue()
	}
	if *remOneSp {
		fmt.Println(">> Remove One Service Profile Called [-rp]")
		err = removeServiceFromOnu(olt)
		if err != nil {
			fmt.Printf("Error running demo: %v\n", err)
		}
		promptContinue()
	}
	if *showSpDetails {
		fmt.Println(">> Show Service Profile Details called [-sp]")
		err = displayServiceProfiles(olt)
		if err != nil {
			fmt.Printf("Error running demo: %v\n", err)
		}
		promptContinue()
	}
}

func manuallyRegisterOnu(olt *gopon.LumiaOlt) error {
	var err error
	var obll *gopon.OnuBlacklistList
	// perform GET request on OLT Blacklist and hold this data in a local variable
	obll, err = olt.GetOnuBlacklist()
	if err != nil {
		return err
	}
	// perform GET request on OLT WhiteList and update app's db of currently provisioned ONU
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	// display the Blacklist entries only
	obll.Tabwrite()
	// interactive: pause and ask for input, first show blacklist
	fmt.Println(">> Provide OLT Port to Register ONU to:")
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)
	oltIntf := sanitizeIntfInput(readFromStdin(reader))
	if oltIntf == "" {
		return gopon.ErrNotInput
	}
	fmt.Println(">> Provide ONU Serial Number to Register:")
	sn := sanitizeSnInput(readFromStdin(reader))
	if sn == "" {
		return gopon.ErrNotInput
	}
	fmt.Printf(">> Input: %s:%s\n", oltIntf, sn)
	// External Auth List held by application, not sent to OLT
	err = olt.AddSnToAuthList(sn)
	if err != nil {
		return err
	}
	// using the supplied OLT Port, find the next available ONU interface 'y' as in 0/x/y
	intf := olt.NextAvailableOnuInterface(oltIntf)
	// generate a new ONU Config with the supplied SN and ONU interface
	ocfg := gopon.NewOnuConfig(sn, intf)
	// Authorize ONU checks if SN exists in AuthorizedOnuSn list
	// If it does, it executes a PATCH request updating the existing Provisioned ONUs with the new ONU Config
	err = olt.AuthorizeOnu(ocfg)
	if err != nil {
		return err
	}
	// perform GET request on OLT WhiteList and update app's db of currently provisioned ONU
	//	err = olt.UpdateOnuRegistry()
	//	if err != nil {
	//		return err
	//	}
	// perform GET request on OLT Service Profiles and display them
	var spl *gopon.ServiceProfileList
	spl, err = olt.GetServiceProfiles()
	if err != nil {
		return err
	}
	spl.Tabwrite()
	fmt.Println(">> Provide Service Profile to apply to ONU:")
	sps := strings.Fields(readFromStdin(reader))
	for _, sp := range sps {
		sp = sanitizeInput(strings.TrimSpace(sp))
		if spl.ProfileExists(sp) {
			// created new ONU Profile by combining registered interface with service profile name
			newProfile := gopon.NewOnuProfile(ocfg.IfName, sp)
			// error will be checked if the server accepts the patch
			err = olt.PostOnuProfile(newProfile)
			if err != nil {
				return err
			}
		} else {
			return gopon.ErrNotInput
		}
	}
	// perform GET request on OLT WhiteList and update app's db of currently provisioned ONU
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	olt.TabwriteRegistry()
	return nil
}

func registerOnuFromFile(olt *gopon.LumiaOlt) error {
	now := time.Now()
	fmt.Println("Starting the Timer")
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
	//      There are two paths that can be taken, pre-register does not have OLT port
	//      But Blacklist register requires multiple devices to be sitting waiting to register
	//      If the wrong OLT Port is supplied as pre-config, OLT will alarm when device comes up
	//      How to get the auth file to include the OLT port as optional?
	//      Can this operation check my database when a blacklist entry appears?
	//      Focus on current need: multiple ONU sit on the Blacklist, want to auth them all and add SP according to file

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

	fmt.Printf("Elapsed operation time: %v\n", time.Since(now))

	return nil
}

func manuallyDeregisterOnu(olt *gopon.LumiaOlt) error {
	var err error
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	if len(olt.Registration) < 1 {
		return gopon.ErrNotExists
	}
	// The ONU Registry is a stand-in for a real database
	olt.TabwriteRegistry()

	fmt.Println(">> Provide Serial Number of ONU to Deregister:")
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)
	sn := sanitizeSnInput(readFromStdin(reader))
	if sn == "" {
		return gopon.ErrNotInput
	}
	err = olt.DeauthOnuBySn(sn)
	if err != nil {
		return err
	}
	fmt.Printf(">> Deregistered ONU Serial Number: [%s]\n", sn)
	// the deauth function took care of updating the database
	// but as a "stateful" check we will perform the update registry
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	olt.TabwriteRegistry()

	return nil
}

func deRegisterOnuFromFile(olt *gopon.LumiaOlt) error {
	now := time.Now()
	fmt.Println("Starting the Timer")
	
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

//	err = olt.UpdateOnuRegistry()
//	if err != nil {
//		return err
//	}
//	olt.TabwriteRegistry()
	// Blacklist might not show the devices yet, takes about 1 min to update
	var obll *gopon.OnuBlacklistList
	obll, err = olt.GetOnuBlacklist()
	if err != nil {
		return err
	}
	if len(obll.Entry) < 1 {
		fmt.Println("No entries on Blacklist yet, waiting 30s")
		for i := 1; i < 31; i++ {
			time.Sleep(1 * time.Second)
			fmt.Print("|")
		}
		fmt.Println("")
		obll, err = olt.GetOnuBlacklist()
		if err != nil {
			return err
		}
	}
	obll.Tabwrite()

	fmt.Printf("Elapsed operation time: %v\n", time.Since(now))
	
	return nil
}

func addServiceToOnu(olt *gopon.LumiaOlt) error {
	var err error
	var spl *gopon.ServiceProfileList
	// first show the available Service Profiles
	spl, err = olt.GetServiceProfiles()
	if err != nil {
		return err
	}
	spl.Tabwrite()
	// next show the current device provisioning
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	olt.TabwriteRegistry()
	// interactive: pause and ask for input, first show blacklist
	fmt.Println(">> Provide ONU Serial Number to Add Service to:")
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)
	sn := sanitizeSnInput(readFromStdin(reader))
	if sn == "" {
		return gopon.ErrNotInput
	}
	fmt.Println(">> Provide Service Profile Name to add to ONU:")
	sp := sanitizeInput(readFromStdin(reader))
	if sp == "" || !spl.ProfileExists(sp) {
		return gopon.ErrNotInput
	}

	onuReg, err := olt.GetOnuRegisterBySn(sn)
	if err != nil {
		return err
	}
	err = olt.AddServiceToOnu(onuReg, sp)
	if err != nil {
		return err
	}
	// perform GET request on OLT WhiteList and update app's db of currently provisioned ONU
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	olt.TabwriteRegistry()

	return nil
}

func removeServiceFromOnu(olt *gopon.LumiaOlt) error {
	// first show the current device provisioning
	err := olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	olt.TabwriteRegistry()
	// interactive: pause and ask for input, first show blacklist
	fmt.Println(">> Provide ONU Serial Number to Remove Service from:")
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)
	sn := sanitizeSnInput(readFromStdin(reader))
	if sn == "" {
		return gopon.ErrNotInput
	}
	fmt.Println(">> Provide Service Profile Name to remove from ONU:")
	sp := sanitizeInput(readFromStdin(reader))
	// don't need to check if profile exists, it will not be removed if it is the case
	if sp == "" {
		return gopon.ErrNotInput
	}

	onuReg, err := olt.GetOnuRegisterBySn(sn)
	if err != nil {
		return err
	}
	err = olt.RemoveOnuProfileUsage(onuReg.Interface, sp)
	if err != nil {
		return err
	}
	// perform GET request on OLT WhiteList and update app's db of currently provisioned ONU
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	olt.TabwriteRegistry()

	return nil
}

func displayServiceProfiles(olt *gopon.LumiaOlt) error {
	// the top level data structure for provisioning services on an OLT is represented by a "Service Profile"
	// we will perform a GET Request to retrieve all currently configured Service Profiles on the OLT
	// a separate object holds a list of the individual profile objects to allow group tabwrite methods
	var err error
	var spl *gopon.ServiceProfileList
	spl, err = olt.GetServiceProfiles()
	if err != nil {
		return err
	}
	spl.Tabwrite()

	var fpl *gopon.FlowProfileList
	fpl, err = olt.GetFlowProfiles()
	if err != nil {
		return err
	}
	fpl.Tabwrite()

	var vpl *gopon.VlanProfileList
	vpl, err = olt.GetVlanProfiles()
	if err != nil {
		return err
	}
	vpl.Tabwrite()

	var ofpl *gopon.OnuFlowProfileList
	ofpl, err = olt.GetOnuFlowProfiles()
	if err != nil {
		return err
	}
	ofpl.Tabwrite()

	var otpl *gopon.OnuTcontProfileList
	otpl, err = olt.GetOnuTcontProfiles()
	if err != nil {
		return err
	}
	otpl.Tabwrite()

	var ovpl *gopon.OnuVlanProfileList
	var ovrl *gopon.OnuVlanRuleList
	ovpl, ovrl, err = olt.GetOnuVlanProfiles()
	if err != nil {
		return err
	}
	ovpl.Tabwrite()
	ovrl.Tabwrite()

	return nil
}

func promptContinue() {
	fmt.Print(">> Continue? (Y/n)")
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)
	input := readFromStdin(reader)
	if input == "" || strings.ToLower(input) == "y" {
		return
	}
	if strings.ToLower(input) == "n" {
		fmt.Println("Exiting...")
		os.Exit(1)
	}
	promptContinue()
}

func readFromStdin(r *bufio.Reader) string {
	a, _, err := r.ReadLine()
	if err == io.EOF {
		return ""
	} else if err != nil {
		panic(err)
	}

	return strings.TrimRight(string(a), "\r\n")
}

func sanitizeSnInput(sn string) string {
	sn = strings.TrimSpace(sn)
	if len(sn) == 12 {
		sn = sanitizeInput(sn)
	}
	if len(sn) == 8 {
		sn = "ISKT" + sn
		sn = sanitizeInput(sn)
	}
	if len(sn) == 12 {
		return sn
	}
	return ""
}

func sanitizeIntfInput(intf string) string {
	str := strings.Split(strings.TrimSpace(intf), "/")
	if len(str) == 2 {
		if str[0] == "0" {
			if len(str[1]) < 3 {
				return intf
			}
		}
	}
	return ""
}

func sanitizeInput(input string) string {
	var allowedChar = []rune{
		'-', '_', '/', '.', '[', ']', '(', ')',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	}
	var isOk bool
	var output string
	for _, c := range input {
		for _, f := range allowedChar {
			if f == c {
				isOk = true
			}
		}
		if isOk {
			output += string(c)
			isOk = false
		}
	}
	return output
}
