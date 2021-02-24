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
	err := manuallyDeregisterOnu(olt)
	fmt.Println(err)
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
	olt.TabwriteRegistry()

	// but as a "stateful" check we can query to OLT directly
	var opl *gopon.OnuProfileList
	opl, err = olt.GetOnuProfileUsage()
	if err != nil {
		return err
	}
	opl.Tabwrite()
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
