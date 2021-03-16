package gopon

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

type IgmpProfile struct {
	Name                     string `json:"msanMulticastProfileName"`
	IgmpSnooping             int    `json:"msanMulticastProfileIgmpSnooping"`
	IgmpSnoopingFastLeave    int    `json:"msanMulticastProfileIgmpSnoopingFastLeave"`
	IgmpSnoopingSuppression  int    `json:"msanMulticastProfileIgmpSnoopingSuppression"`
	IgmpProxy                int    `json:"msanMulticastProfileIgmpProxy"`
	IgmpProxyIPAddress       string `json:"msanMulticastProfileIgmpProxyIpAddress"`
	IgmpFiltering            int    `json:"msanMulticastProfileIgmpFiltering"`
	MulticastGroupLimit      int    `json:"msanMulticastProfileMulticastGroupLimit"`
	Mvr                      int    `json:"msanMulticastProfileMvr"`
	IgmpProxyProtocolVersion int    `json:"msanMulticastProfileIgmpProxyProtocolVersion"`
	Usage                    int    `json:"msanMulticastProfileUsage"`
}

type IgmpProfileList struct {
	Entry []*IgmpProfile
}

func NewIgmpProfile(name string) *IgmpProfile {
	p := &IgmpProfile{
		Name:                     name,
		IgmpSnooping:             0,
		IgmpSnoopingFastLeave:    1,
		IgmpSnoopingSuppression:  0,
		IgmpProxy:                0,
		IgmpProxyIPAddress:       "", // defaults to OLT Management IP
		IgmpFiltering:            1,
		MulticastGroupLimit:      0,
		Mvr:                      0,
		IgmpProxyProtocolVersion: 2,
	}
	return p
}

func (p *IgmpProfile) GetName() string {
	return p.Name
}

// Copy returns a copy of the profile object with a new name and Usage set to 2
func (p *IgmpProfile) Copy(newName string) (*IgmpProfile, error) {
	if p.Name == newName {
		return nil, ErrExists
	}
	np := p
	np.Name = newName
	np.Usage = 2
	return np, nil
}

func (p *IgmpProfile) GetIgmpSnooping() bool {
	return p.IgmpSnooping == 1
}

func (p *IgmpProfile) SetIgmpSnooping(state bool) {
	if state {
		p.IgmpSnooping = 1
	} else {
		p.IgmpSnooping = 0
	}
}

func (p *IgmpProfile) GetFastLeave() bool {
	return p.IgmpSnoopingFastLeave == 1
}

func (p *IgmpProfile) SetFastLeave(state bool) {
	if state {
		p.IgmpSnoopingFastLeave = 1
	} else {
		p.IgmpSnoopingFastLeave = 0
	}
}

func (p *IgmpProfile) GetSuppression() bool {
	return p.IgmpSnoopingSuppression == 1
}

func (p *IgmpProfile) SetSuppression(state bool) {
	if state {
		p.IgmpSnoopingSuppression = 1
	} else {
		p.IgmpSnoopingSuppression = 0
	}
}

func (p *IgmpProfile) GetIgmpProxy() bool {
	return p.IgmpProxy == 1
}

// Optional to specify the Proxy IP Address when enabling
func (p *IgmpProfile) SetIgmpProxy(state bool, address ...string) {
	if state {
		p.IgmpProxy = 1
		if len(address) != 0 {
			p.IgmpProxyIPAddress = address[0]
		}
	} else {
		p.IgmpSnooping = 0
	}
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (p *IgmpProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", data
	}
	return p.Name, data
}

/*
Name                     string `json:"msanMulticastProfileName"`
	IgmpSnooping             int    `json:"msanMulticastProfileIgmpSnooping"`
	IgmpSnoopingFastLeave    int    `json:"msanMulticastProfileIgmpSnoopingFastLeave"`
	IgmpSnoopingSuppression  int    `json:"msanMulticastProfileIgmpSnoopingSuppression"`
	IgmpProxy                int    `json:"msanMulticastProfileIgmpProxy"`
	IgmpProxyIPAddress       string `json:"msanMulticastProfileIgmpProxyIpAddress"`
	IgmpFiltering            int    `json:"msanMulticastProfileIgmpFiltering"`
	MulticastGroupLimit      int    `json:"msanMulticastProfileMulticastGroupLimit"`
	Mvr                      int    `json:"msanMulticastProfileMvr"`
	IgmpProxyProtocolVersion
*/
var IgmpProfileHeaders = []string{
	"Name",
	"Snooping",
	"Fast-Leave",
	"Suppression",
	"Proxy",
	"Filtering",
	"Version",
}

// ListEssentialParams returns a map of the essential IgmpProfile parameters
func (p *IgmpProfile) ListEssentialParams() map[string]interface{} {
	var EssentialIgmpProfile = map[string]interface{}{
		IgmpProfileHeaders[0]: p.GetName(),
		IgmpProfileHeaders[1]: p.GetIgmpSnooping(),
		IgmpProfileHeaders[2]: p.GetFastLeave(),
		IgmpProfileHeaders[3]: p.GetSuppression(),
		IgmpProfileHeaders[4]: p.GetIgmpProxy(),
		IgmpProfileHeaders[5]: p.IgmpProxyProtocolVersion,
	}
	// I want all of these Bools to return strings of "Enabled/Disabled"
	return EssentialIgmpProfile
}

// Tabwrite displays the essential information of VlanProfile in organized columns
func (p *IgmpProfile) Tabwrite() {
	fmt.Println("|| IGMP Profile ||")
	l := p.ListEssentialParams()
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	for _, v := range IgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range IgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range IgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", l[v])
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range IgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	tw.Flush()
}

// Separate is a method to maintain backward-compatability
func (ipl *IgmpProfileList) Separate() []*IgmpProfile {
	var entry *IgmpProfile
	var list []*IgmpProfile
	for _, e := range ipl.Entry {
		entry = e
		list = append(list, entry)
	}
	return list
}

// Tabwrite displays the essential information of a list of Flow Profiles in organized columns
func (ipl *IgmpProfileList) Tabwrite() {
	fmt.Println("|| IGMP Profile List ||")
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range IgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range IgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	ips := ipl.Separate()
	for _, ip := range ips {
		// first get the data as a map
		l := ip.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range IgmpProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range IgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}
