package goPon

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrNotStruct     = errors.New("Not a valid struct")
	ErrNotField      = errors.New("Not a valid field name")
	ErrNotExported   = errors.New("Not an exported field")
	ErrNotSettable   = errors.New("Not a settable field")
	ErrNotInput      = errors.New("Incorrect input supplied")
	ErrNotExists     = errors.New("Key not found")
	ErrExists        = errors.New("Key already exists")
	ErrInUse         = errors.New("Cannot modify while in use")
	ErrNotStatusOk   = errors.New("Did not receive 200 OK from HTTP server")
	ErrNotReachable  = errors.New("Host not reachable")
	ErrNotAuthorized = errors.New("Onu Sn not on Authorized List")
)

const (
	auth             = "session=em+protection-user=admin&em+protection-pw=admin"
	empty            = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
	responseOk       = "200 OK"
	serviceProfiles  = "msanServiceProfileTable"
	flowProfiles     = "msanServiceFlowProfileTable"
	vlanProfiles     = "msanVlanProfileTable"
	igmpProfiles     = "msanMulticastProfileTable"
	securityProfiles = "msanSecurityProfileTable"
	onuFlowProfiles  = "msanOnuFlowProfileTable"
	onuTcontProfiles = "msanOnuTcontProfileTable"
	onuVlanProfiles  = "msanOnuVlanProfileTable"
	onuVlanRules     = "msanOnuVlanProfileRuleTable"
	onuIgmpProfiles  = "msanOnuMulticastProfileTable"
	l2cpProfiles     = "msanL2cpProfileTable"
	onuBlacklist     = "msanOnuBlackListTable"
	onuConfig        = "msanOnuCfgTable"
	onuProfiles      = "msanServicePortProfileTable"
	onuInfo		 = "msanOnuInfoTable"
)

const (
	// iota counter for KiloBits, MegaBits, GigaBits
	Kb = 1 << (10 * iota)
	Mb
	Gb
)

// formatKbits converts kb values to Mb and Gb as applicable, round down harshly
func formatKbits(n int) string {
	value, unit := float64(n), "k"
	switch {
	case value >= Gb:
		value, unit = value/Gb, "G"
	case value >= Mb:
		value, unit = value/Mb, "M"
	}
	r := strconv.FormatFloat(value, 'f', 0, 64)
	r += unit
	return r
}

// fillSpace fills the tabwriter space by length of the column title
func fs(s string) string {
	return strings.Repeat("-", len(s))
}

// Base64 encoded VLAN entries resolved to a list of VLANs
func getVlanFromB64(b string) (found []int, err error) {
	//fmt.Println(b)
	p, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		return
	}
	for i, x := range p {
		if x != 0 {
			by := fmt.Sprintf("%08b", x)
			bi := strings.SplitN(by, "", 8)
			for o, v := range bi {
				if v != "0" {
					f := (i * 8) + o
					found = append(found, f)
				}
			}
		}
	}
	//fmt.Println(found)
	return
}

// list of VLANs encoded to a Base64 representative string of bits representing VLAN membership 0-4095
// this should have error handling to prevent incorrect numbers supplied
func getB64FromVlan(find []int) (string, error) {
	by := make([]byte, 512)
	for _, z := range find {
		if z < 1 || z > 4095 {
			return empty, ErrNotInput
		}
		a := z / 8 // returns the byte array reference
		b := z % 8 // returns the bit array reference
		c := 7 - b
		//fmt.Println(z, a, b)
		// bitwise OR | critical to add found bit while maintaining existing
		// bitwise shift left << works from the right, so c is necessary to reverse order
		by[a] = by[a] | 1<<c
		//fmt.Printf("%08b\n", by[a])
	}
	f := base64.StdEncoding.EncodeToString(by)
	//fmt.Printf("\nB64 Encode Function Called; Input: %v\nOutput: %s\n", find, f)
	return f, nil
}

// string input (such as from text file or stdin) of space separated ints from 1-4095 as input, int list as output
func generateVlanList(h string) (o []int, err error) {
	c := strings.Split(h, " ") // using a string to hold a list of VLAN ids "100 101 200"
	for x := range c {
		d, err := strconv.Atoi(c[x])
		if err != nil {
			return nil, err
		}
		o = append(o, d)
	}
	return o, nil
}

func toString(v interface{}) string {
	switch vv := v.(type) {
	case []byte:
		return string(vv)
	case string:
		return vv
	case bool:
		return strconv.FormatBool(vv)
	case int:
		return strconv.FormatInt(int64(vv), 10)
	case int8:
		return strconv.FormatInt(int64(vv), 10)
	case int16:
		return strconv.FormatInt(int64(vv), 10)
	case int32:
		return strconv.FormatInt(int64(vv), 10)
	case int64:
		return strconv.FormatInt(int64(vv), 10)
	case uint:
		return strconv.FormatUint(uint64(vv), 10)
	case uint8:
		return strconv.FormatUint(uint64(vv), 10)
	case uint16:
		return strconv.FormatUint(uint64(vv), 10)
	case uint32:
		return strconv.FormatUint(uint64(vv), 10)
	case uint64:
		return strconv.FormatUint(uint64(vv), 10)
	case float32:
		return strconv.FormatFloat(float64(vv), 'f', 2, 64)
	case float64:
		return strconv.FormatFloat(float64(vv), 'f', 2, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}


func UrlEncodeInterface(intf string) string {
	return strings.ReplaceAll(intf, "/", "%2F")
}

func UrlDecodeInterface(intf string) string {
	return strings.ReplaceAll(intf, "%2F", "/")
}

func parseAuth(auth string) (user, pass string) {
	str := strings.Split(auth, "&")
	for _, st := range str {
		s := strings.Split(st, "=")
		if strings.Contains(s[1], "user") {
			user = s[2]
		}
		if strings.Contains(s[0], "pw") {
			pass = s[1]
		}
	}
	return

}
