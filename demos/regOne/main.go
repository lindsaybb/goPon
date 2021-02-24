package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lindsaybb/gopon"
)

var (
	helpFlag = flag.Bool("h", false, "Show this help")
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
	err := manuallyRegisterOnu(olt)
	fmt.Println(err)
}

func manuallyRegisterOnu(olt *gopon.LumiaOlt) error {
	var err error
	var obll *gopon.OnuBlacklistList
	obll, err = olt.GetOnuBlacklist()
	if err != nil {
		return err
	}
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
	// perform GET request on OLT WhiteList and update app's db of currently provisioned ONU
	err = olt.UpdateOnuRegistry()
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
	err = olt.UpdateOnuRegistry()
	if err != nil {
		return err
	}
	// perform GET request on OLT Service Profiles and display them
	var spl *gopon.ServiceProfileList
	spl, err = olt.GetServiceProfiles()
	if err != nil {
		return err
	}
	spl.Tabwrite()
	fmt.Println(">> Provide Service Profile to apply to ONU (comma-separated for multiple):")
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
