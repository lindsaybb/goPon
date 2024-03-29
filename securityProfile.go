package goPon

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type SecurityProfile struct {
	Name                   string `json:"msanSecurityProfileName"`
	ProtectedPort          int    `json:"msanSecurityProfileProtectedPort"`
	MacSg                  int    `json:"msanSecurityProfileMacSg"`
	MacLimit               int    `json:"msanSecurityProfileMacLimit"`
	PortSecurity           int    `json:"msanSecurityProfilePortSecurity"`
	ArpInspect             int    `json:"msanSecurityProfileArpInspec"`
	IPSg                   int    `json:"msanSecurityProfileIpSg"`
	IPSgIpv6               int    `json:"msanSecurityProfileIpSgIpv6"`
	IPSgFilteringMode      int    `json:"msanSecurityProfileIpSgFilteringMode"`
	IPSgBindingLimit       int    `json:"msanSecurityProfileIpSgBindingLimit"`
	IPSgBindingLimitDhcpv6 int    `json:"msanSecurityProfileIpSgBindingLimitDhcpv6"`
	IPSgBindingLimitND     int    `json:"msanSecurityProfileIpSgBindingLimitND"`
	StormControlBroadcast  int    `json:"msanSecurityProfileStormControlBroadcast"`
	StormControlMulticast  int    `json:"msanSecurityProfileStormControlMulticast"`
	StormControlUnicast    int    `json:"msanSecurityProfileStormControlUnicast"`
	AppRateLimitDhcp       int    `json:"msanSecurityProfileAppRateLimitDhcp"`
	AppRateLimitIgmp       int    `json:"msanSecurityProfileAppRateLimitIgmp"`
	AppRateLimitPppoe      int    `json:"msanSecurityProfileAppRateLimitPppoe"`
	AppRateLimitStp        int    `json:"msanSecurityProfileAppRateLimitStp"`
	AppRateLimitMn         int    `json:"msanSecurityProfileAppRateLimitMn"`
	Usage                  int    `json:"msanSecurityProfileUsage"`
}

type SecurityProfileList struct {
	Entry []*SecurityProfile
}

func NewSecurityProfile(name string) *SecurityProfile {
	p := &SecurityProfile{
		Name:                   name,
		ProtectedPort:          1,
		MacSg:                  0,
		MacLimit:               0,
		PortSecurity:           0,
		ArpInspect:             0,
		IPSg:                   0,
		IPSgIpv6:               0,
		IPSgFilteringMode:      2,
		IPSgBindingLimit:       4,
		IPSgBindingLimitDhcpv6: 4,
		IPSgBindingLimitND:     4,
		StormControlBroadcast:  -1,
		StormControlMulticast:  -1,
		StormControlUnicast:    100,
		AppRateLimitDhcp:       5,
		AppRateLimitIgmp:       5,
		AppRateLimitPppoe:      5,
		AppRateLimitStp:        5,
		AppRateLimitMn:         1000,
	}
	return p
}

func (p *SecurityProfile) GetName() string {
	return p.Name
}

// IsUsed
func (p *SecurityProfile) IsUsed() bool {
	return p.Usage == 1
}

// Copy returns a copy of the profile object with a new name and Usage set to 2
func (p *SecurityProfile) Copy(newName string) (*SecurityProfile, error) {
	if p.Name == newName {
		return nil, ErrExists
	}
	np := p
	np.Name = newName
	np.Usage = 2
	return np, nil
}

var SecStmCtlList = []string{
	"Broadcast",
	"Unknown-Unicast",
	"Multicast",
}

// BUM order
func (p *SecurityProfile) GetStormControl() []int {
	var stmctl = []int{
		p.StormControlBroadcast,
		p.StormControlUnicast,
		p.StormControlMulticast,
	}
	return stmctl
}

func (p *SecurityProfile) GetStormControlString() string {
	list := p.GetStormControl()
	return fmt.Sprintf("%s: %d, %s: %d, %s: %d", SecStmCtlList[0], list[0], SecStmCtlList[1], list[1], SecStmCtlList[2], list[2])
}

func (p *SecurityProfile) SetStormControl(list []int) error {
	switch len(list) {
	case 0:
		return ErrNotInput
	case 1:
		list = append(list, []int{0, 0}...)
	case 2:
		list = append(list, 0)
	}
	for i := 0; i < 3; i++ {
		if list[i] > 65535 {
			list[i] = 65535
		}
	}
	p.StormControlBroadcast = list[0]
	p.StormControlUnicast = list[1]
	p.StormControlMulticast = list[2]
	return nil
}

var SecArlList = []string{
	"DHCP",
	"IGMP",
	"PPPOE",
	"STP",
	"MN",
}

// DHCP, IGMP, PPPoE, STP, MN
func (p *SecurityProfile) GetAppRateLimit() []int {
	var arl = []int{
		p.AppRateLimitDhcp,
		p.AppRateLimitIgmp,
		p.AppRateLimitPppoe,
		p.AppRateLimitStp,
		p.AppRateLimitMn,
	}
	return arl
}

func (p *SecurityProfile) GetAppRateLimitString() string {
	list := p.GetAppRateLimit()
	return fmt.Sprintf("%s: %d, %s: %d, %s: %d, %s: %d, %s: %d", SecArlList[0], list[0], SecArlList[1], list[1], SecArlList[2], list[2], SecArlList[3], list[3], SecArlList[4], list[4])
}

// allows setting an individual ARL value by providing the key:value of the desired setting
func (p *SecurityProfile) SetAppRateLimit(key string, value int) error {
	key = strings.ToUpper(key)
	if value > 1000 {
		value = 1000
	}
	for i, v := range SecArlList {
		if key == v {
			return p.SetARL(i, value)
		}
	}
	return ErrNotInput
}

// change the ARL individually according to a key relating to the SecArlList, with supplied value (-1 for default value)
func (p *SecurityProfile) SetARL(key int, value int) error {
	switch key {
	case 0:
		p.AppRateLimitDhcp = value
	case 1:
		p.AppRateLimitIgmp = value
	case 2:
		p.AppRateLimitPppoe = value
	case 3:
		p.AppRateLimitStp = value
	case 4:
		p.AppRateLimitMn = value
	default:
		return ErrNotInput
	}

	return nil
}

func (p *SecurityProfile) DefaultAppRateLimit() {
	p.AppRateLimitDhcp = 5
	p.AppRateLimitIgmp = 5
	p.AppRateLimitPppoe = 5
	p.AppRateLimitStp = 5
	p.AppRateLimitMn = 1000
}

// returns true if ProtectedPort is enabled
func (p *SecurityProfile) GetProtectedPort() bool {
	return p.ProtectedPort == 1
}

// supply desired state as bool where enabled=true
func (p *SecurityProfile) SetProtectedPort(state bool) {
	if state {
		p.ProtectedPort = 1
	} else {
		p.ProtectedPort = 0
	}
}

func (p *SecurityProfile) GetMacSG() bool {
	return p.MacSg == 1
}

func (p *SecurityProfile) SetMacSG(state bool) {
	if state {
		p.MacSg = 1
	} else {
		p.MacSg = 0
	}
}

func (p *SecurityProfile) GetMacLimit() int {
	return p.MacLimit
}

func (p *SecurityProfile) SetMacLimit(limit int) {
	if limit < 1 {
		limit = 1
	}
	if limit > 16 {
		limit = 16
	}
	p.MacLimit = limit
}

func (p *SecurityProfile) GetPortSecurity() bool {
	return p.PortSecurity == 1
}

func (p *SecurityProfile) SetPortSecurity(state bool) {
	if state {
		p.PortSecurity = 1
	} else {
		p.PortSecurity = 0
	}
}

func (p *SecurityProfile) GetArpInspect() bool {
	return p.ArpInspect == 1
}

func (p *SecurityProfile) SetArpInspect(state bool) {
	if state {
		p.ArpInspect = 1
	} else {
		p.ArpInspect = 0
	}
}

var SecIpSgList = []string{
	"v4Enable",
	"v6Enable",
	"FilterMode",
	"v4BindingLimit",
	"v6BindingLimitDHCP",
	"v6BindingLimitND",
}

// expand this to include elements of configuration
func (p *SecurityProfile) GetIPv4SG() bool {
	return p.IPSg == 1
}

func (p *SecurityProfile) SetIPv4SG(state bool) {
	if state {
		p.IPSg = 1
	} else {
		p.IPSg = 0
	}
}

func (p *SecurityProfile) GetFilterMode() bool {
	return p.IPSgFilteringMode == 2
}

func (p *SecurityProfile) GetFilterModeString() string {
	switch p.IPSgFilteringMode {
	case 1:
		return "IP"
	case 2:
		return "IP & MAC"
	}
	return ""
}

func (p *SecurityProfile) SetFilterMode(state bool) {
	if state {
		p.IPSgFilteringMode = 2
	} else {
		p.IPSgFilteringMode = 1
	}
}

func (p *SecurityProfile) GetIPv6SG() bool {
	return p.IPSgIpv6 == 1
}

func (p *SecurityProfile) SetIPv6SG(state bool) {
	if state {
		p.IPSgIpv6 = 1
	} else {
		p.IPSgIpv6 = 0
	}
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (p *SecurityProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", data
	}
	return p.Name, data
}

var SecurityProfileHeaders = []string{
	"Name",
	"Port-Protect",
	"MAC-SG",
	"MAC-Limit",
	"Port-Sec",
	"Arp-Inspect",
	"IPv4-SG",
	"IPv6-SG",
	"Storm-Control",
	"AppRateLimit",
}

// ListEssentialParams returns a map of the essential VlanProfile parameters
func (p *SecurityProfile) ListEssentialParams() map[string]interface{} {
	var EssentialSecurityProfile = map[string]interface{}{
		SecurityProfileHeaders[0]: p.GetName(),
		SecurityProfileHeaders[1]: p.GetProtectedPort(),
		SecurityProfileHeaders[2]: p.GetMacSG(),
		SecurityProfileHeaders[3]: p.GetMacLimit(),
		SecurityProfileHeaders[4]: p.GetPortSecurity(),
		SecurityProfileHeaders[5]: p.GetArpInspect(),
		SecurityProfileHeaders[6]: p.GetIPv4SG(),
		SecurityProfileHeaders[7]: p.GetIPv6SG(),
		SecurityProfileHeaders[8]: p.GetStormControl(),
		SecurityProfileHeaders[9]: p.GetAppRateLimit(),
	}
	// I want all of these Bools to return strings of "Enabled/Disabled"
	return EssentialSecurityProfile
}

// Tabwrite displays the essential information of VlanProfile in organized columns
func (p *SecurityProfile) Tabwrite() {
	fmt.Println("|| Security Profile ||")
	l := p.ListEssentialParams()
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	for _, v := range SecurityProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range SecurityProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range SecurityProfileHeaders {
		fmt.Fprintf(tw, "%v\t", l[v])
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range SecurityProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	tw.Flush()
}

// Separate is a method to maintain backward-compatability
func (spl *SecurityProfileList) Separate() []*SecurityProfile {
	var entry *SecurityProfile
	var list []*SecurityProfile
	for _, e := range spl.Entry {
		entry = e
		list = append(list, entry)
	}
	return list
}

// Tabwrite displays the essential information of a list of Flow Profiles in organized columns
func (spl *SecurityProfileList) Tabwrite() {
	fmt.Println("|| Security Profile List ||")
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range SecurityProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range SecurityProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	sps := spl.Separate()
	for _, sp := range sps {
		// first get the data as a map
		l := sp.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range SecurityProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range SecurityProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}
