package goPon

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

// FlowProfile is the complete Flow profile data struct (ordered in order it appears as json)
type FlowProfile struct {
	Name                      string `json:"msanServiceFlowProfileName"`
	MatchUsAny                int    `json:"msanServiceFlowProfileMatchUsAny"`
	MatchUsMacDestAddr        string `json:"msanServiceFlowProfileMatchUsMacDestAddr"`
	MatchUsMacDestMask        string `json:"msanServiceFlowProfileMatchUsMacDestMask"`
	MatchUsMacSrcAddr         string `json:"msanServiceFlowProfileMatchUsMacSrcAddr"`
	MatchUsMacSrcMask         string `json:"msanServiceFlowProfileMatchUsMacSrcMask"`
	MatchUsCPcp               int    `json:"msanServiceFlowProfileMatchUsCPcp"`
	MatchUsSPcp               int    `json:"msanServiceFlowProfileMatchUsSPcp"`
	MatchUsVlanProfile        int    `json:"msanServiceFlowProfileMatchUsVlanProfile"`
	MatchUsCVlanIDRange       string `json:"msanServiceFlowProfileMatchUsCVlanIdRange"`
	MatchUsSVlanIDRange       string `json:"msanServiceFlowProfileMatchUsSVlanIdRange"`
	MatchUsEthertype          int    `json:"msanServiceFlowProfileMatchUsEthertype"`
	MatchUsIPProtocol         int    `json:"msanServiceFlowProfileMatchUsIpProtocol"`
	MatchUsIPSrcAddr          string `json:"msanServiceFlowProfileMatchUsIpSrcAddr"`
	MatchUsIPSrcMask          string `json:"msanServiceFlowProfileMatchUsIpSrcMask"`
	MatchUsIPDestAddr         string `json:"msanServiceFlowProfileMatchUsIpDestAddr"`
	MatchUsIPDestMask         string `json:"msanServiceFlowProfileMatchUsIpDestMask"`
	MatchUsIPDscp             int    `json:"msanServiceFlowProfileMatchUsIpDscp"`
	MatchUsIPCsc              int    `json:"msanServiceFlowProfileMatchUsIpCsc"`
	MatchUsIPDropPrecedence   int    `json:"msanServiceFlowProfileMatchUsIpDropPrecedence"`
	MatchUsTCPSrcPort         int    `json:"msanServiceFlowProfileMatchUsTcpSrcPort"`
	MatchUsTCPDestPort        int    `json:"msanServiceFlowProfileMatchUsTcpDestPort"`
	MatchUsUDPSrcPort         int    `json:"msanServiceFlowProfileMatchUsUdpSrcPort"`
	MatchUsUDPDstPort         int    `json:"msanServiceFlowProfileMatchUsUdpDstPort"`
	MatchUsIpv6SrcAddr        string `json:"msanServiceFlowProfileMatchUsIpv6SrcAddr"`
	MatchUsIpv6SrcAddrMaskLen int    `json:"msanServiceFlowProfileMatchUsIpv6SrcAddrMaskLen"`
	MatchUsIpv6DstAddr        string `json:"msanServiceFlowProfileMatchUsIpv6DstAddr"`
	MatchUsIpv6DstAddrMaskLen int    `json:"msanServiceFlowProfileMatchUsIpv6DstAddrMaskLen"`
	MatchDsAny                int    `json:"msanServiceFlowProfileMatchDsAny"`
	MatchDsMacDestAddr        string `json:"msanServiceFlowProfileMatchDsMacDestAddr"`
	MatchDsMacDestMask        string `json:"msanServiceFlowProfileMatchDsMacDestMask"`
	MatchDsMacSrcAddr         string `json:"msanServiceFlowProfileMatchDsMacSrcAddr"`
	MatchDsMacSrcMask         string `json:"msanServiceFlowProfileMatchDsMacSrcMask"`
	MatchDsCPcp               int    `json:"msanServiceFlowProfileMatchDsCPcp"`
	MatchDsSPcp               int    `json:"msanServiceFlowProfileMatchDsSPcp"`
	MatchDsVlanProfile        int    `json:"msanServiceFlowProfileMatchDsVlanProfile"`
	MatchDsCVlanIDRange       string `json:"msanServiceFlowProfileMatchDsCVlanIdRange"`
	MatchDsSVlanIDRange       string `json:"msanServiceFlowProfileMatchDsSVlanIdRange"`
	MatchDsEthertype          int    `json:"msanServiceFlowProfileMatchDsEthertype"`
	MatchDsIPProtocol         int    `json:"msanServiceFlowProfileMatchDsIpProtocol"`
	MatchDsIPSrcAddr          string `json:"msanServiceFlowProfileMatchDsIpSrcAddr"`
	MatchDsIPSrcMask          string `json:"msanServiceFlowProfileMatchDsIpSrcMask"`
	MatchDsIPDestAddr         string `json:"msanServiceFlowProfileMatchDsIpDestAddr"`
	MatchDsIPDestMask         string `json:"msanServiceFlowProfileMatchDsIpDestMask"`
	MatchDsIPDscp             int    `json:"msanServiceFlowProfileMatchDsIpDscp"`
	MatchDsIPCsc              int    `json:"msanServiceFlowProfileMatchDsIpCsc"`
	MatchDsIPDropPrecedence   int    `json:"msanServiceFlowProfileMatchDsIpDropPrecedence"`
	MatchDsTCPSrcPort         int    `json:"msanServiceFlowProfileMatchDsTcpSrcPort"`
	MatchDsTCPDestPort        int    `json:"msanServiceFlowProfileMatchDsTcpDestPort"`
	MatchDsUDPSrcPort         int    `json:"msanServiceFlowProfileMatchDsUdpSrcPort"`
	MatchDsUDPDstPort         int    `json:"msanServiceFlowProfileMatchDsUdpDstPort"`
	MatchDsIpv6SrcAddr        string `json:"msanServiceFlowProfileMatchDsIpv6SrcAddr"`
	MatchDsIpv6SrcAddrMaskLen int    `json:"msanServiceFlowProfileMatchDsIpv6SrcAddrMaskLen"`
	MatchDsIpv6DstAddr        string `json:"msanServiceFlowProfileMatchDsIpv6DstAddr"`
	MatchDsIpv6DstAddrMaskLen int    `json:"msanServiceFlowProfileMatchDsIpv6DstAddrMaskLen"`
	UsCdr                     int    `json:"msanServiceFlowProfileUsCdr"`
	UsCdrBurstSize            int    `json:"msanServiceFlowProfileUsCdrBurstSize"`
	UsPdr                     int    `json:"msanServiceFlowProfileUsPdr"`
	UsPdrBurstSize            int    `json:"msanServiceFlowProfileUsPdrBurstSize"`
	UsMarkPcp                 int    `json:"msanServiceFlowProfileUsMarkPcp"`
	UsMarkPcpValue            int    `json:"msanServiceFlowProfileUsMarkPcpValue"`
	UsMarkDscp                int    `json:"msanServiceFlowProfileUsMarkDscp"`
	UsMarkDscpValue           int    `json:"msanServiceFlowProfileUsMarkDscpValue"`
	DsCdr                     int    `json:"msanServiceFlowProfileDsCdr"`
	DsCdrBurstSize            int    `json:"msanServiceFlowProfileDsCdrBurstSize"`
	DsPdr                     int    `json:"msanServiceFlowProfileDsPdr"`
	DsPdrBurstSize            int    `json:"msanServiceFlowProfileDsPdrBurstSize"`
	DsMarkPcp                 int    `json:"msanServiceFlowProfileDsMarkPcp"`
	DsMarkPcpValue            int    `json:"msanServiceFlowProfileDsMarkPcpValue"`
	DsMarkDscp                int    `json:"msanServiceFlowProfileDsMarkDscp"`
	DsMarkDscpValue           int    `json:"msanServiceFlowProfileDsMarkDscpValue"`
	DsQueuingPriority         int    `json:"msanServiceFlowProfileDsQueuingPriority"`
	DsSchedulingMode          int    `json:"msanServiceFlowProfileDsSchedulingMode"`
	Usage                     int    `json:"msanServiceFlowProfileUsage"`
}

type FlowProfileList struct {
	Entry []*FlowProfile
}

// NewFlowProfile returns a FlowProfile with default values (function is called in order of YANG model parameter definition)
func NewFlowProfile(name string) *FlowProfile {
	p := &FlowProfile{
		Name:                      name,
		MatchUsAny:                2,     // disabled, default
		MatchUsMacDestAddr:        "",    // nil, def
		MatchUsMacDestMask:        "",    // nil, def
		MatchUsMacSrcAddr:         "",    // nil, def
		MatchUsMacSrcMask:         "",    // nil, def
		MatchUsCPcp:               -1,    // not defined, def
		MatchUsSPcp:               -1,    // not defined, def
		MatchUsVlanProfile:        2,     // default is disabled (2)
		MatchUsCVlanIDRange:       empty, // nil is default, empty string is returned from not setting value
		MatchUsSVlanIDRange:       empty, // nil is default, empty string is returned from not setting value
		MatchUsEthertype:          -1,    // not defined, def
		MatchUsIPProtocol:         -1,    // not defined, nil, [1:icmp, 2:igmp, 4:ip, 6:tcp, 17:udp]
		MatchUsIPSrcAddr:          "",    // nil, def
		MatchUsIPSrcMask:          "",    // nil, def
		MatchUsIPDestAddr:         "",    // nil, def
		MatchUsIPDestMask:         "",    // nil, def
		MatchUsIPDscp:             -1,    // not defined, def
		MatchUsIPCsc:              -1,    // not defined, def
		MatchUsIPDropPrecedence:   -1,    // not defined, def [0:noDrop, 1:lowDrop, 2:mediumDrop, 3:highDrop] (two-bit value 00, 01, 10, 11)
		MatchUsTCPSrcPort:         -1,    // not defined, def
		MatchUsTCPDestPort:        -1,    // not defined, def
		MatchUsUDPSrcPort:         -1,    // not defined, def
		MatchUsUDPDstPort:         -1,    // not defined, def
		MatchDsAny:                2,     // disabled, default
		MatchDsMacDestAddr:        "",    // nil, def
		MatchDsMacDestMask:        "",    // nil, def
		MatchDsMacSrcAddr:         "",    // nil, def
		MatchDsMacSrcMask:         "",    // nil, def
		MatchDsCPcp:               -1,    // not defined, def
		MatchDsSPcp:               -1,    // not defined, def
		MatchDsVlanProfile:        2,     // default is disabled (2)
		MatchDsCVlanIDRange:       empty, // nil is default, empty string is returned from not setting value
		MatchDsSVlanIDRange:       empty, // nil is default, empty string is returned from not setting value
		MatchDsEthertype:          -1,    // not defined, def
		MatchDsIPProtocol:         -1,    // not defined, nil, [1:icmp, 2:igmp, 4:ip, 6:tcp, 17:udp]
		MatchDsIPSrcAddr:          "",    // nil, def
		MatchDsIPSrcMask:          "",    // nil, def
		MatchDsIPDestAddr:         "",    // nil, def
		MatchDsIPDestMask:         "",    // nil, def
		MatchDsIPDscp:             -1,    // not defined, def
		MatchDsIPCsc:              -1,    // not defined, def
		MatchDsIPDropPrecedence:   -1,    // not defined, def [0:noDrop, 1:lowDrop, 2:mediumDrop, 3:highDrop] (two-bit value 00, 01, 10, 11)
		MatchDsTCPSrcPort:         -1,    // not defined, def
		MatchDsTCPDestPort:        -1,    // not defined, def
		MatchDsUDPSrcPort:         -1,    // not defined, def
		MatchDsUDPDstPort:         -1,    // not defined, def
		UsCdr:                     0,     // not defined, def
		UsCdrBurstSize:            0,     // not defined, def
		UsPdr:                     0,     // not defined, def
		UsPdrBurstSize:            0,     // not defined, def
		UsMarkPcp:                 1,     // not defined, def
		UsMarkPcpValue:            -1,    // not defined, def
		UsMarkDscp:                1,     // not defined, def
		UsMarkDscpValue:           -1,    // not defined, def
		DsCdr:                     0,     // not defined, def
		DsCdrBurstSize:            0,     // not defined, def
		DsPdr:                     0,     // not defined, def
		DsPdrBurstSize:            0,     // not defined, def
		DsMarkPcp:                 1,     // not defined, def
		DsMarkPcpValue:            -1,    // not defined, def
		DsMarkDscp:                1,     // not defined, def
		DsMarkDscpValue:           -1,    // not defined, def
		DsQueuingPriority:         0,     //  not defined, def
		DsSchedulingMode:          1,     // weighted, def [2: strict]
		MatchUsIpv6SrcAddr:        "",    // nil, def
		MatchUsIpv6SrcAddrMaskLen: 0,     // not defined, def
		MatchUsIpv6DstAddr:        "",    // nil, def
		MatchUsIpv6DstAddrMaskLen: 0,     // not defined, def
		MatchDsIpv6SrcAddr:        "",    // nil, def
		MatchDsIpv6SrcAddrMaskLen: 0,     // not defined, def
		MatchDsIpv6DstAddr:        "",    // nil, def
		MatchDsIpv6DstAddrMaskLen: 0,     // not defined, def
	}
	return p
}

// GetName returns the name of the FlowProfile
func (p *FlowProfile) GetName() string {
	return p.Name
}

// IsUsed returns true if the profile is currently attached to one or more subscribers
func (p *FlowProfile) IsUsed() bool {
	return p.Usage == 1
}

// Copy returns a copy of the profile object with a new name and Usage set to 2
func (p *FlowProfile) Copy(newName string) (*FlowProfile, error) {
	if p.Name == newName {
		return nil, ErrExists
	}
	np := p
	np.Name = newName
	np.Usage = 2
	return np, nil
}

// GetMatchBothVlanProfile returns a bool of if the FlowProfile logic is set to match what is set in the Vlan Profile in the same Service Profile, the most common scenario
func (p *FlowProfile) GetMatchBothVlanProfile() bool {
	return p.GetMatchUsVlanProfile() && p.GetMatchDsVlanProfile()
}

// GetMatchUsVlanProfile returns a bool of if the FlowProfile logic is set to match Us what is set in the Vlan Profile in the same Service Profile
func (p *FlowProfile) GetMatchUsVlanProfile() bool {
	return p.MatchUsVlanProfile == 1
}

// GetMatchDsVlanProfile returns a bool of if the FlowProfile logic is set to match Ds what is set in the Vlan Profile in the same Service Profile
func (p *FlowProfile) GetMatchDsVlanProfile() bool {
	return p.MatchDsVlanProfile == 1
}

// SetMatchBothVlanProfile sets FlowProfile logic to match that set in the Vlan Profile in the same Service Profile
func (p *FlowProfile) SetMatchBothVlanProfile() {
	p.MatchUsVlanProfile = 1
	p.MatchDsVlanProfile = 1
}

var FlowProfileHeaders = []string{
	"Name",
	"UsMatchVlanProfile",
	"DsMatchVlanProfile",
	"UsMatchOther",
	"DsMatchOther",
	"UsHandling",
	"DsHandling",
	"QueuingPriority",
	"SchedulingMode",
}

// ListEssentialParams returns a map of the essential FlowProfile parameters
func (p *FlowProfile) ListEssentialParams() map[string]interface{} {
	var EssentialFlowProfile = map[string]interface{}{
		FlowProfileHeaders[0]: p.GetName(),               // string
		FlowProfileHeaders[1]: p.GetMatchUsVlanProfile(), // bool
		FlowProfileHeaders[2]: p.GetMatchDsVlanProfile(), // bool
		FlowProfileHeaders[3]: p.GetMatchUsOther(),       // collection of values, returns any non-nil
		FlowProfileHeaders[4]: p.GetMatchDsOther(),
		FlowProfileHeaders[5]: p.GetUsHandling(),
		FlowProfileHeaders[6]: p.GetDsHandling(),
		FlowProfileHeaders[7]: p.GetQueueingPriority(),
		FlowProfileHeaders[8]: p.GetSchedulingMode(),
	}

	return EssentialFlowProfile
}

var FlowProfileUsOther = []string{
	"MatchUsAny",
	"MatchUsMacDestAddr",
	"MatchUsMacDestMask",
	"MatchUsMacSrcAddr",
	"MatchUsMacSrcMask",
	"MatchUsCPcp",
	"MatchUsSPcp",
	"MatchUsCVlanIDRange",
	"MatchUsSVlanIDRange",
	"MatchUsEthertype",
	"MatchUsIPProtocol",
	"MatchUsIPSrcAddr",
	"MatchUsIPSrcMask",
	"MatchUsIPDestAddr",
	"MatchUsIPDestMask",
	"MatchUsIPDscp",
	"MatchUsIPCsc",
	"MatchUsIPDropPrecedence",
	"MatchUsTCPSrcPort",
	"MatchUsTCPDestPort",
	"MatchUsUDPSrcPort",
	"MatchUsUDPDstPort",
	"MatchUsIpv6SrcAddr",
	"MatchUsIpv6SrcAddrMaskLen",
	"MatchUsIpv6DstAddr",
	"MatchUsIpv6SrcAddrMaskLen",
}

func (p *FlowProfile) IsMatchUsAny() bool {
	return p.MatchUsAny == 2
}

func (p *FlowProfile) IsMatchDsAny() bool {
	return p.MatchDsAny == 2
}

// GetMatchUsOther returns a list of any non-default values as a key:value pair
func (p *FlowProfile) GetMatchUsOther() []interface{} {
	var out []interface{}

	if !p.IsMatchUsAny() {
		out = append(out, map[string]int{FlowProfileUsOther[0]: p.MatchUsAny})
	}
	if p.MatchUsMacDestAddr != "" {
		out = append(out, map[string]string{FlowProfileUsOther[1]: p.MatchUsMacDestAddr})
	}
	if p.MatchUsMacDestMask != "" {
		out = append(out, map[string]string{FlowProfileUsOther[2]: p.MatchUsMacDestMask})
	}
	if p.MatchUsMacSrcAddr != "" {
		out = append(out, map[string]string{FlowProfileUsOther[3]: p.MatchUsMacSrcAddr})
	}
	if p.MatchUsMacSrcMask != "" {
		out = append(out, map[string]string{FlowProfileUsOther[4]: p.MatchUsMacSrcMask})
	}
	if p.MatchUsCPcp != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[5]: p.MatchUsCPcp})
	}
	if p.MatchUsSPcp != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[6]: p.MatchUsSPcp})
	}
	if p.MatchUsCVlanIDRange != empty {
		list, err := getVlanFromB64(p.MatchUsCVlanIDRange)
		if err == nil {
			out = append(out, map[string][]int{FlowProfileUsOther[7]: list})
		} else {
			log.Printf("flowProfile: %v\n", err)
		}
	}
	if p.MatchUsSVlanIDRange != empty {
		list, err := getVlanFromB64(p.MatchUsSVlanIDRange)
		if err == nil {
			out = append(out, map[string][]int{FlowProfileUsOther[8]: list})
		} else {
			log.Printf("flowProfile: %v\n", err)
		}
	}
	if p.MatchUsEthertype != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[9]: p.MatchUsEthertype})
	}
	if p.MatchUsIPProtocol != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[10]: p.MatchUsIPProtocol})
	}
	if p.MatchUsIPSrcAddr != "" {
		out = append(out, map[string]string{FlowProfileUsOther[11]: p.MatchUsIPSrcAddr})
	}
	if p.MatchUsIPSrcMask != "" {
		out = append(out, map[string]string{FlowProfileUsOther[12]: p.MatchUsIPSrcMask})
	}
	if p.MatchUsIPDestAddr != "" {
		out = append(out, map[string]string{FlowProfileUsOther[13]: p.MatchUsIPDestAddr})
	}
	if p.MatchUsIPDestMask != "" {
		out = append(out, map[string]string{FlowProfileUsOther[14]: p.MatchUsIPDestMask})
	}
	if p.MatchUsIPDscp != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[15]: p.MatchUsIPDscp})
	}
	if p.MatchUsIPCsc != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[16]: p.MatchUsIPCsc})
	}
	if p.MatchUsIPDropPrecedence != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[17]: p.MatchUsIPDropPrecedence})
	}
	if p.MatchUsTCPSrcPort != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[18]: p.MatchUsTCPSrcPort})
	}
	if p.MatchUsTCPDestPort != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[19]: p.MatchUsTCPDestPort})
	}
	if p.MatchUsUDPSrcPort != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[20]: p.MatchUsUDPSrcPort})
	}
	if p.MatchUsUDPDstPort != -1 {
		out = append(out, map[string]int{FlowProfileUsOther[21]: p.MatchUsUDPDstPort})
	}
	if p.MatchUsIpv6SrcAddr != "" {
		out = append(out, map[string]string{FlowProfileUsOther[22]: p.MatchUsIpv6SrcAddr})
	}
	if p.MatchUsIpv6SrcAddrMaskLen != 0 {
		out = append(out, map[string]int{FlowProfileUsOther[23]: p.MatchUsIpv6SrcAddrMaskLen})
	}
	if p.MatchUsIpv6DstAddr != "" {
		out = append(out, map[string]string{FlowProfileUsOther[24]: p.MatchUsIpv6DstAddr})
	}
	if p.MatchUsIpv6DstAddrMaskLen != 0 {
		out = append(out, map[string]int{FlowProfileUsOther[25]: p.MatchUsIpv6SrcAddrMaskLen})
	}
	return out
}

var FlowProfileDsOther = []string{
	"MatchDsAny",
	"MatchDsMacDestAddr",
	"MatchDsMacDestMask",
	"MatchDsMacSrcAddr",
	"MatchDsMacSrcMask",
	"MatchDsCPcp",
	"MatchDsSPcp",
	"MatchDsCVlanIDRange",
	"MatchDsSVlanIDRange",
	"MatchDsEthertype",
	"MatchDsIPProtocol",
	"MatchDsIPSrcAddr",
	"MatchDsIPSrcMask",
	"MatchDsIPDestAddr",
	"MatchDsIPDestMask",
	"MatchDsIPDscp",
	"MatchDsIpCsc",
	"MatchDsIpDropPrecedence",
	"MatchDsTCPSrcPort",
	"MatchDsTCPDestPort",
	"MatchDsUDPSrcPort",
	"MatchDsUDPDstPort",
	"MatchDsIpv6SrcAddr",
	"MatchDsIpv6SrcAddrMaskLen",
	"MatchDsIpv6DstAddr",
	"MatchDsIpv6DstAddrMaskLen",
}

// GetMatchDsOther returns a list of any non-default values as a key:value pair
func (p *FlowProfile) GetMatchDsOther() []interface{} {
	var out []interface{}

	if p.MatchDsAny != 2 {
		out = append(out, map[string]int{FlowProfileDsOther[0]: p.MatchDsAny})
	}
	if p.MatchDsMacDestAddr != "" {
		out = append(out, map[string]string{FlowProfileDsOther[1]: p.MatchUsMacDestAddr})
	}
	if p.MatchDsMacDestMask != "" {
		out = append(out, map[string]string{FlowProfileDsOther[2]: p.MatchDsMacDestMask})
	}
	if p.MatchDsMacSrcAddr != "" {
		out = append(out, map[string]string{FlowProfileDsOther[3]: p.MatchDsMacSrcAddr})
	}
	if p.MatchDsMacSrcMask != "" {
		out = append(out, map[string]string{FlowProfileDsOther[4]: p.MatchDsMacSrcMask})
	}
	if p.MatchDsCPcp != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[5]: p.MatchDsCPcp})
	}
	if p.MatchDsSPcp != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[6]: p.MatchDsSPcp})
	}
	if p.MatchDsCVlanIDRange != empty {
		list, err := getVlanFromB64(p.MatchDsCVlanIDRange)
		if err == nil {
			out = append(out, map[string][]int{FlowProfileDsOther[7]: list})
		} else {
			log.Printf("flowProfile: %v\n", err)
		}
	}
	if p.MatchDsSVlanIDRange != empty {
		list, err := getVlanFromB64(p.MatchDsSVlanIDRange)
		if err == nil {
			out = append(out, map[string][]int{FlowProfileDsOther[8]: list})
		} else {
			log.Printf("flowProfile: %v\n", err)
		}
	}
	if p.MatchDsEthertype != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[9]: p.MatchDsEthertype})
	}
	if p.MatchDsIPProtocol != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[10]: p.MatchDsIPProtocol})
	}
	if p.MatchDsIPSrcAddr != "" {
		out = append(out, map[string]string{FlowProfileDsOther[11]: p.MatchDsIPSrcAddr})
	}
	if p.MatchDsIPSrcMask != "" {
		out = append(out, map[string]string{FlowProfileDsOther[12]: p.MatchDsIPSrcMask})
	}
	if p.MatchDsIPDestAddr != "" {
		out = append(out, map[string]string{FlowProfileDsOther[13]: p.MatchDsIPDestAddr})
	}
	if p.MatchDsIPDestMask != "" {
		out = append(out, map[string]string{FlowProfileDsOther[14]: p.MatchDsIPDestMask})
	}
	if p.MatchDsIPDscp != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[15]: p.MatchDsIPDscp})
	}
	if p.MatchDsIPCsc != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[16]: p.MatchDsIPCsc})
	}
	if p.MatchDsIPDropPrecedence != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[17]: p.MatchDsIPDropPrecedence})
	}
	if p.MatchDsTCPSrcPort != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[18]: p.MatchDsTCPSrcPort})
	}
	if p.MatchDsTCPDestPort != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[19]: p.MatchDsTCPDestPort})
	}
	if p.MatchDsUDPSrcPort != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[20]: p.MatchDsUDPSrcPort})
	}
	if p.MatchDsUDPDstPort != -1 {
		out = append(out, map[string]int{FlowProfileDsOther[21]: p.MatchDsUDPDstPort})
	}
	if p.MatchDsIpv6SrcAddr != "" {
		out = append(out, map[string]string{FlowProfileDsOther[22]: p.MatchDsIpv6SrcAddr})
	}
	if p.MatchDsIpv6SrcAddrMaskLen != 0 {
		out = append(out, map[string]int{FlowProfileDsOther[23]: p.MatchDsIpv6SrcAddrMaskLen})
	}
	if p.MatchDsIpv6DstAddr != "" {
		out = append(out, map[string]string{FlowProfileDsOther[24]: p.MatchDsIpv6DstAddr})
	}
	if p.MatchDsIpv6DstAddrMaskLen != 0 {
		out = append(out, map[string]int{FlowProfileDsOther[25]: p.MatchDsIpv6DstAddrMaskLen})
	}
	return out
}

var FlowProfileUsHandling = []string{
	"UsCdr",
	"UsCdrBurstSize",
	"UsPdr",
	"UsPdrBurstSize",
	"UsMarkPcp",
	"UsMarkPcpValue",
	"UsMarkDscp",
	"UsMarkDscpValue",
}

// GetUsHandling returns a list of any non-default values as a key:value pair
func (p *FlowProfile) GetUsHandling() []interface{} {
	var out []interface{}

	if p.UsCdr != 0 {
		out = append(out, map[string]int{FlowProfileUsHandling[0]: p.UsCdr})
	}
	if p.UsCdrBurstSize != 0 {
		out = append(out, map[string]int{FlowProfileUsHandling[1]: p.UsCdrBurstSize})
	}
	if p.UsPdr != 0 {
		out = append(out, map[string]int{FlowProfileUsHandling[2]: p.UsPdr})
	}
	if p.UsPdrBurstSize != 0 {
		out = append(out, map[string]int{FlowProfileUsHandling[3]: p.UsPdrBurstSize})
	}
	if p.UsMarkPcp != 1 {
		out = append(out, map[string]int{FlowProfileUsHandling[4]: p.UsMarkPcp})
	}
	if p.UsMarkPcpValue != -1 {
		out = append(out, map[string]int{FlowProfileUsHandling[5]: p.UsMarkPcpValue})
	}
	if p.UsMarkDscp != 1 {
		out = append(out, map[string]int{FlowProfileUsHandling[6]: p.UsMarkDscp})
	}
	if p.UsMarkDscpValue != -1 {
		out = append(out, map[string]int{FlowProfileUsHandling[7]: p.UsMarkDscpValue})
	}

	return out
}

var FlowProfileDsHandling = []string{
	"DsCdr",
	"DsCdrBurstSize",
	"DsPdr",
	"DsPdrBurstSize",
	"DsMarkPcp",
	"DsMarkPcpValue",
	"DsMarkDscp",
	"DsMarkDscpValue",
}

// GetDsHandling returns a list of any non-default values as a key:value pair
func (p *FlowProfile) GetDsHandling() []interface{} {
	var out []interface{}

	if p.DsCdr != 0 {
		out = append(out, map[string]int{FlowProfileDsHandling[0]: p.DsCdr})
	}
	if p.DsCdrBurstSize != 0 {
		out = append(out, map[string]int{FlowProfileDsHandling[1]: p.DsCdrBurstSize})
	}
	if p.DsPdr != 0 {
		out = append(out, map[string]int{FlowProfileDsHandling[2]: p.DsPdr})
	}
	if p.DsPdrBurstSize != 0 {
		out = append(out, map[string]int{FlowProfileDsHandling[3]: p.DsPdrBurstSize})
	}
	if p.DsMarkPcp != 1 {
		out = append(out, map[string]int{FlowProfileDsHandling[4]: p.DsMarkPcp})
	}
	if p.DsMarkPcpValue != -1 {
		out = append(out, map[string]int{FlowProfileDsHandling[5]: p.DsMarkPcpValue})
	}
	if p.DsMarkDscp != 1 {
		out = append(out, map[string]int{FlowProfileDsHandling[6]: p.DsMarkDscp})
	}
	if p.DsMarkDscpValue != -1 {
		out = append(out, map[string]int{FlowProfileDsHandling[7]: p.DsMarkDscpValue})
	}

	return out
}

// GetQueueingPriority returns a string of the Ds Priority Bit, default: 0
func (p *FlowProfile) GetQueueingPriority() string {
	return fmt.Sprintf("%d", p.DsQueuingPriority)
}

var FlowProfileSchedulingModes = []string{
	"nil",
	"Weighted",
	"Strict",
}

// GetSchedulingMode returns a string explaining the profile-configured Ds Scheduling Mode
func (p *FlowProfile) GetSchedulingMode() string {
	return FlowProfileSchedulingModes[p.DsSchedulingMode]
}

// Tabwrite displays the essential information of FlowProfile in organized columns
func (p *FlowProfile) Tabwrite() {
	fmt.Println("|| Flow Profile ||")
	l := p.ListEssentialParams()
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	for _, v := range FlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range FlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range FlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", l[v])
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range FlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	tw.Flush()
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (p *FlowProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", data
	}
	return p.Name, data
}

// Separate is a method to maintain backward-compatability
func (fpl *FlowProfileList) Separate() []*FlowProfile {
	var entry *FlowProfile
	var list []*FlowProfile
	for _, e := range fpl.Entry {
		entry = e
		list = append(list, entry)
	}
	return list
}

// Tabwrite displays the essential information of a list of Flow Profiles in organized columns
func (fpl *FlowProfileList) Tabwrite() {
	fmt.Println("|| Flow Profile List ||")
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range FlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range FlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	fps := fpl.Separate()
	for _, fp := range fps {
		// first get the data as a map
		l := fp.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range FlowProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range FlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}
