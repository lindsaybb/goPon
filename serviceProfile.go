package goPon

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

// need to fully implement bitmap logic for UNI ports, two examples given
// create a convenience method that returns ALL sub-profile objects in one go

// ServiceProfile is a collection of the sub-profiles needed to enable a Service on an ONU
type ServiceProfile struct {
	Name                            string `json:"msanServiceProfileName"`
	FlowProfileName                 string `json:"msanServiceProfileServiceFlowProfileName"`
	MulticastProfileName            string `json:"msanServiceProfileMulticastProfileName"`
	VlanProfileName                 string `json:"msanServiceProfileVlanProfileName"`
	L2cpProfileName                 string `json:"msanServiceProfileL2cpProfileName"`
	SecurityProfileName             string `json:"msanServiceProfileSecurityProfileName"`
	OnuFlowProfileName              string `json:"msanServiceProfileOnuFlowProfileName"`
	OnuVlanProfileName              string `json:"msanServiceProfileOnuVlanProfileName"`
	OnuMulticastProfileName         string `json:"msanServiceProfileOnuMulticastProfileName"`
	OnuTcontProfileName             string `json:"msanServiceProfileOnuTcontProfileName"`
	OnuVirtGemPortID                int    `json:"msanServiceProfileOnuVirtGemPortId"`
	OnuTpType                       int    `json:"msanServiceProfileOnuTpType"`
	OnuTpUniBitMap                  string `json:"msanServiceProfileOnuTpUniBitMap"`
	DhcpRa                          int    `json:"msanServiceProfileDhcpRa"`
	DhcpRaTrustClients              int    `json:"msanServiceProfileDhcpRaTrustClients"`
	DhcpRaOpt82UnicastExtension     int    `json:"msanServiceProfileDhcpRaOpt82UnicastExtension"`
	DhcpRaOpt82Insert               int    `json:"msanServiceProfileDhcpRaOpt82Insert"`
	DhcpRaRateLimit                 int    `json:"msanServiceProfileDhcpRaRateLimit"`
	DhcpRaCircuitIDCustomFormat     string `json:"msanServiceProfileDhcpRaCircuitIdCustomFormat"`
	DhcpRaRemoteIDCustomFormat      string `json:"msanServiceProfileDhcpRaRemoteIdCustomFormat"`
	DhcpRaCircuitIDType             int    `json:"msanServiceProfileDhcpRaCircuitIdType"`
	Dhcpv6Ra                        int    `json:"msanServiceProfileDhcpv6Ra"`
	Dhcpv6RaTrustClients            int    `json:"msanServiceProfileDhcpv6RaTrustClients"`
	Dhcpv6RaRemoteIDEnterpriseNum   int    `json:"msanServiceProfileDhcpv6RaRemoteIdEnterpriseNum"`
	Dhcpv6RaInterfaceIDType         int    `json:"msanServiceProfileDhcpv6RaInterfaceIdType"`
	Dhcpv6RaInterfaceIDCustomFormat string `json:"msanServiceProfileDhcpv6RaInterfaceIdCustomFormat"`
	Dhcpv6RaRemoteIDCustomFormat    string `json:"msanServiceProfileDhcpv6RaRemoteIdCustomFormat"`
	PppoeIA                         int    `json:"msanServiceProfilePppoeIA"`
	PppoeIARateLimit                int    `json:"msanServiceProfilePppoeIARateLimit"`
	PPPoeIACircuitIDType            int    `json:"msanServiceProfilePPPoeIACircuitIdType"`
	PPPoeIACircuitIDCustomFormat    string `json:"msanServiceProfilePPPoeIACircuitIdCustomFormat"`
	PPPoeIARemoteIDCustomFormat     string `json:"msanServiceProfilePPPoeIARemoteIdCustomFormat"`
	Usage                           int    `json:"msanServiceProfileUsage"`
}

type ServiceProfileList struct {
	Entry []*ServiceProfile
}

var OnuTpTypeList = []string{
	"nil",
	"VEIP",
	"IPHOST",
	"UNI",
}

// ConvertOnuTPToString is a helper function to convert the logic used to represent termination point to a readable format
func ConvertOnuTPToString(tp int) string {
	if tp > len(OnuTpTypeList) - 1 || tp < 1 {
		return ""
	}
	return OnuTpTypeList[tp]
}

// ConvertOnuTPUniBitMapToInt is a helper function to convert the logic used to represent UNI physical port to a readable format
func ConvertOnuTPUniBitMapToInt(bitmap string) int {
	// haven't implemented this logic yet
	if bitmap == "QAAA" {
		return 1
	}
	if bitmap == "IAAA" {
		return 2
	}
	return 0
}

// ConvertOnuTPUniBitMapFromInt is a helper function to convert from int to a bitmap using the required logic for representing a UNI physical port
func ConvertOnuTPUniBitMapFromInt(id int) string {
	if id == 1 {
		return "QAAA"
	}
	if id == 2 {
		return "IAAA"
	}
	return "AAAA"
}

// NewServiceProfile returns an empty, initialzed struct to be populated
func NewServiceProfile(name string) *ServiceProfile {
	sp := &ServiceProfile{
		Name:                            name,
		FlowProfileName:                 "",
		MulticastProfileName:            "",
		VlanProfileName:                 "",
		L2cpProfileName:                 "",
		SecurityProfileName:             "",
		OnuFlowProfileName:              "",
		OnuVlanProfileName:              "",
		OnuMulticastProfileName:         "",
		OnuTcontProfileName:             "",
		OnuVirtGemPortID:                1,
		OnuTpType:                       1,
		OnuTpUniBitMap:                  "AAAA",
		DhcpRa:                          0,
		DhcpRaTrustClients:              0,
		DhcpRaOpt82UnicastExtension:     0,
		DhcpRaOpt82Insert:               0,
		DhcpRaRateLimit:                 5,
		DhcpRaCircuitIDCustomFormat:     "",
		DhcpRaRemoteIDCustomFormat:      "",
		DhcpRaCircuitIDType:             1,
		Dhcpv6Ra:                        0,
		Dhcpv6RaTrustClients:            0,
		Dhcpv6RaRemoteIDEnterpriseNum:   1332,
		Dhcpv6RaInterfaceIDType:         2,
		Dhcpv6RaInterfaceIDCustomFormat: "",
		Dhcpv6RaRemoteIDCustomFormat:    "",
		PppoeIA:                         0,
		PppoeIARateLimit:                5,
		PPPoeIACircuitIDType:            1,
		PPPoeIACircuitIDCustomFormat:    "",
		PPPoeIARemoteIDCustomFormat:     "",
	}
	return sp
}

func (sp *ServiceProfile) GetName() string {
	return sp.Name
}

func (sp *ServiceProfile) IsUsed() bool {
	return sp.Usage == 1
}

// Copy returns a copy of the profile object with a new name and Usage set to 2
func (sp *ServiceProfile) Copy(newName string) (*ServiceProfile, error) {
	if sp.Name == newName {
		return nil, ErrExists
	}

	nsp := sp
	nsp.Name = newName
	nsp.Usage = 2
	return nsp, nil
}

// ServiceProfileEssentialHeaders ensure correct order of entries is maintained for Tabwriter
var ServiceProfileEssentialHeaders = []string{
	"Name",
	"Flow Profile",
	"VLAN Profile",
	"ONU Flow Profile",
	"ONU TCONT Profile",
	"ONU VLAN Profile",
	"Virtual GEM Port",
	"ONU TP Type",
}

// ListEssentialSubProfiles lists currently provisioned values of a Service Profile in a map of profile:name
func (sp *ServiceProfile) ListEssentialSubProfiles() map[string]interface{} {
	var EssentialServiceProfile = map[string]interface{}{
		ServiceProfileEssentialHeaders[0]: sp.GetName(),
		ServiceProfileEssentialHeaders[1]: sp.FlowProfileName,
		ServiceProfileEssentialHeaders[2]: sp.VlanProfileName,
		ServiceProfileEssentialHeaders[3]: sp.OnuFlowProfileName,
		ServiceProfileEssentialHeaders[4]: sp.OnuTcontProfileName,
		ServiceProfileEssentialHeaders[5]: sp.OnuVlanProfileName,
		ServiceProfileEssentialHeaders[6]: sp.OnuVirtGemPortID,
		ServiceProfileEssentialHeaders[7]: ConvertOnuTPToString(sp.OnuTpType),
	}

	return EssentialServiceProfile
}

// ServiceProfileHeaders ensure correct order of entries is maintained for Tabwriter
var ServiceProfileHeaders = []string{
	"Name",
	"Flow Profile",
	"VLAN Profile",
	"ONU Flow Profile",
	"ONU TCONT Profile",
	"ONU VLAN Profile",
	"Virtual GEM Port",
	"ONU TP Type",
	"Security Profile",
	"IGMP Profile",
	"ONU IGMP Profile",
	"L2CP Profile",
	"DHCP RA",
	"PPPoE IA",
}

// ListSubProfiles lists currently provisioned values of a Service Profile in a map of profile:name
func (sp *ServiceProfile) ListSubProfiles() map[string]interface{} {
	var ServiceProfile = map[string]interface{}{
		ServiceProfileHeaders[0]:  sp.GetName(),
		ServiceProfileHeaders[1]:  sp.FlowProfileName,
		ServiceProfileHeaders[2]:  sp.VlanProfileName,
		ServiceProfileHeaders[3]:  sp.OnuFlowProfileName,
		ServiceProfileHeaders[4]:  sp.OnuTcontProfileName,
		ServiceProfileHeaders[5]:  sp.OnuVlanProfileName,
		ServiceProfileHeaders[6]:  sp.OnuVirtGemPortID,
		ServiceProfileHeaders[7]:  ConvertOnuTPToString(sp.OnuTpType),
		ServiceProfileHeaders[8]:  sp.SecurityProfileName,
		ServiceProfileHeaders[9]:  sp.MulticastProfileName,
		ServiceProfileHeaders[10]: sp.OnuMulticastProfileName,
		ServiceProfileHeaders[11]: sp.L2cpProfileName,
		ServiceProfileHeaders[12]: sp.GetDhcpRaNonDefaults(),
		ServiceProfileHeaders[13]: sp.GetPppoeIaNonDefaults(),
	}

	return ServiceProfile
}

// GetDhcpRaNonDefaults returns a slice of any values that are not default for parameters related to DHCP Relay Agent
func (sp *ServiceProfile) GetDhcpRaNonDefaults() []interface{} {
	var out []interface{}

	if sp.DhcpRa != 0 {
		out = append(out, map[string]int{"DhcpRa": sp.DhcpRa})
	}
	if sp.DhcpRaTrustClients != 0 {
		out = append(out, map[string]int{"TrustClients": sp.DhcpRaTrustClients})
	}
	if sp.DhcpRaOpt82UnicastExtension != 0 {
		out = append(out, map[string]int{"Opt82UnicastExtension": sp.DhcpRaOpt82UnicastExtension})
	}
	if sp.DhcpRaOpt82Insert != 0 {
		out = append(out, map[string]int{"Opt82Insert": sp.DhcpRaOpt82Insert})
	}
	if sp.DhcpRaRateLimit != 5 {
		out = append(out, map[string]int{"RateLimit": sp.DhcpRaRateLimit})
	}
	if sp.DhcpRaCircuitIDCustomFormat != "" {
		out = append(out, map[string]string{"CircuitIDCustomFormat": sp.DhcpRaCircuitIDCustomFormat})
	}
	if sp.DhcpRaRemoteIDCustomFormat != "" {
		out = append(out, map[string]string{"RemoteIDCustomFormat": sp.DhcpRaRemoteIDCustomFormat})
	}
	if sp.DhcpRaCircuitIDType != 1 {
		out = append(out, map[string]int{"CircuitIDType": sp.DhcpRaCircuitIDType})
	}
	if sp.Dhcpv6Ra != 0 {
		out = append(out, map[string]int{"Dhcpv6Ra": sp.Dhcpv6Ra})
	}
	if sp.Dhcpv6RaTrustClients != 0 {
		out = append(out, map[string]int{"v6TrustClients": sp.Dhcpv6RaTrustClients})
	}
	if sp.Dhcpv6RaRemoteIDEnterpriseNum != 1332 {
		out = append(out, map[string]int{"v6RemoteIDEnterpriseNum": sp.Dhcpv6RaRemoteIDEnterpriseNum})
	}
	if sp.Dhcpv6RaInterfaceIDType != 2 {
		out = append(out, map[string]int{"v6InterfaceIDType": sp.Dhcpv6RaInterfaceIDType})
	}
	if sp.Dhcpv6RaInterfaceIDCustomFormat != "" {
		out = append(out, map[string]string{"v6InterfaceIDCustomFormat": sp.Dhcpv6RaInterfaceIDCustomFormat})
	}
	if sp.Dhcpv6RaRemoteIDCustomFormat != "" {
		out = append(out, map[string]string{"v6RemoteIDCustomFormat": sp.Dhcpv6RaRemoteIDCustomFormat})
	}

	return out
}

// GetPppoeIaNonDefaults returns a slice of any values that are not default for parameters related to PPPoE Intermediate Agent
func (sp *ServiceProfile) GetPppoeIaNonDefaults() []interface{} {
	var out []interface{}

	if sp.PppoeIA != 0 {
		out = append(out, map[string]int{"PppoeIA": sp.PppoeIA})
	}
	if sp.PppoeIARateLimit != 5 {
		out = append(out, map[string]int{"RateLimit": sp.PppoeIARateLimit})
	}
	if sp.PPPoeIACircuitIDType != 1 {
		out = append(out, map[string]int{"CircuitIDType": sp.PPPoeIACircuitIDType})
	}
	if sp.PPPoeIACircuitIDCustomFormat != "" {
		out = append(out, map[string]string{"CircuitIDCustomFormat": sp.PPPoeIACircuitIDCustomFormat})
	}
	if sp.PPPoeIARemoteIDCustomFormat != "" {
		out = append(out, map[string]string{"RemoteIDCustomFormat": sp.PPPoeIARemoteIDCustomFormat})
	}

	return out
}

// SetFlowProfile allows a name to be specified for the ServiceFlowProfile parameter, will error on Post if doesn't already exist in NE
func (sp *ServiceProfile) SetFlowProfile(name string) {
	sp.FlowProfileName = name
}

// SetMulticastProfile allows a name to be specified for the MulticastProfile parameter, will error on Post if doesn't already exist in NE
func (sp *ServiceProfile) SetMulticastProfile(name string) {
	sp.MulticastProfileName = name
}

// SetVlanProfile allows a name to be specified for the VlanProfile parameter, will error on Post if doesn't already exist in NE
func (sp *ServiceProfile) SetVlanProfile(name string) {
	sp.VlanProfileName = name
}

// SetL2cpProfile allows a name to be specified for the L2cpProfile parameter, will error on Post if doesn't already exist in NE
func (sp *ServiceProfile) SetL2cpProfile(name string) {
	sp.L2cpProfileName = name
}

// SetSecurityProfile allows a name to be specified for the SecurityProfile parameter, will error on Post if doesn't already exist in NE
func (sp *ServiceProfile) SetSecurityProfile(name string) {
	sp.SecurityProfileName = name
}

// SetOnuFlowProfile allows a name to be specified for the OnuFlowProfile parameter, will error on Post if doesn't already exist in NE
func (sp *ServiceProfile) SetOnuFlowProfile(name string) {
	sp.OnuFlowProfileName = name
}

// SetOnuVlanProfile allows a name to be specified for the OnuVlanProfile parameter, will error on Post if doesn't already exist in NE
func (sp *ServiceProfile) SetOnuVlanProfile(name string) {
	sp.OnuVlanProfileName = name
}

// SetOnuMulticastProfile allows a name to be specified for the OnuMulticastProfile parameter, will error on Post if doesn't already exist in NE
func (sp *ServiceProfile) SetOnuMulticastProfile(name string) {
	sp.OnuMulticastProfileName = name
}

// SetOnuTcontProfile allows a name to be specified for the OnuTcontProfile parameter, will error on Post if doesn't already exist in NE
func (sp *ServiceProfile) SetOnuTcontProfile(name string) {
	sp.OnuTcontProfileName = name
}

// SetVirtualGemPort allows a number from 1-32 to be specified for the OnuVirtGemPortID parameter
func (sp *ServiceProfile) SetVirtualGemPort(id int) {
	if id < 1 || id > 32 {
		id = 1
	}
	sp.OnuVirtGemPortID = id
}

// SetOnuTpType allows a number from 1-3 to be specified for the OnuTpType parameter, defined as [1: VEIP, 2:IPHOST, 3: UNI]
func (sp *ServiceProfile) SetOnuTpType(id int) {
	if id < 1 || id > 3 {
		id = 1
	}
	sp.OnuTpType = id
}

// SetOnuTpUniBitMap allows a number from 1-16 to be specified for the OnuTpUniBitMap parameter, mapping to bitmap is handled indirectly
func (sp *ServiceProfile) SetOnuTpUniBitMap(id int) {
	if id < 1 || id > 16 {
		id = 1
	}
	sp.OnuTpUniBitMap = ConvertOnuTPUniBitMapFromInt(id)
}

// Tabwrite displays the essential information of Service Profile in organized columns
func (sp *ServiceProfile) Tabwrite() {
	fmt.Println("|| Service Profile ||")
	// first get the data as a map
	l := sp.ListEssentialSubProfiles()
	// initiate a tabwriter
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values
	for _, v := range ServiceProfileEssentialHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range ServiceProfileEssentialHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// use the Headers as key in the data map to display values in columns
	for _, v := range ServiceProfileEssentialHeaders {
		fmt.Fprintf(tw, "%v\t", l[v])
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range ServiceProfileEssentialHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}

// Tabwrite displays the essential information of Service Profile in organized columns
func (sp *ServiceProfile) TabwriteFull() {
	fmt.Println("|| Service Profile ||")
	// first get the data as a map
	l := sp.ListSubProfiles()
	// initiate a tabwriter
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values
	for _, v := range ServiceProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range ServiceProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// use the Headers as key in the data map to display values in columns
	for _, v := range ServiceProfileHeaders {
		fmt.Fprintf(tw, "%v\t", l[v])
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range ServiceProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (sp *ServiceProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(sp)
	if err != nil {
		return "", data
	}
	return sp.Name, data
}

// Separate is a method to maintain backward-compatability
func (spl *ServiceProfileList) Separate() []*ServiceProfile {
	var entry *ServiceProfile
	var list []*ServiceProfile
	for _, e := range spl.Entry {
		entry = e
		list = append(list, entry)
	}
	return list
}

// NameSlice generate a []string of the Service Profile Names
func (spl *ServiceProfileList) NameSlice() []string {
	var sl []string
	for i := 0; i < len(spl.Entry); i++ {
		sl = append(sl, spl.Entry[i].Name)
	}
	return sl
}

// ProfileExists ensures the supplied Service Profile name is already provisioned on the OLT
func (spl *ServiceProfileList) ProfileExists(name string) bool {
	sl := spl.NameSlice()
	for _, sp := range sl {
		if name == sp {
			return true
		}
	}
	return false
}

// Tabwrite displays the essential information of a list of Service Profiles in organized columns
func (spl *ServiceProfileList) Tabwrite() {
	fmt.Println("|| Service Profile List ||")
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range ServiceProfileEssentialHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range ServiceProfileEssentialHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	sps := spl.Separate()
	for _, sp := range sps {
		// for each service profile get the data as a map
		l := sp.ListEssentialSubProfiles()
		// iterate over the map using the header as string key
		for _, v := range ServiceProfileEssentialHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range ServiceProfileEssentialHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}

// Tabwrite displays the full information of a list of Service Profiles in organized columns
func (spl *ServiceProfileList) TabwriteFull() {
	fmt.Println("|| Service Profile List ||")
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range ServiceProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range ServiceProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	sps := spl.Separate()
	for _, sp := range sps {
		// for each service profile get the data as a map
		l := sp.ListSubProfiles()
		// iterate over the map using the header as string key
		for _, v := range ServiceProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range ServiceProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}
