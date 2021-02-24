package gopon

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

// OnuFlowProfile contains the data structure for Onu Flow Profile handling
type OnuFlowProfile struct {
	Name                string `json:"msanOnuFlowProfileName"`
	MatchUsCVlanIDRange string `json:"msanOnuFlowProfileMatchUsCVlanIdRange"`
	MatchUsCPcp         int    `json:"msanOnuFlowProfileMatchUsCPcp"`
	UsCdr               int    `json:"msanOnuFlowProfileUsCdr"`
	UsPdr               int    `json:"msanOnuFlowProfileUsPdr"`
	UsFlowPriority      int    `json:"msanOnuFlowProfileUsFlowPriority"`
	DsFlowPriority      int    `json:"msanOnuFlowProfileDsFlowPriority"`
	Usage               int    `json:"msanOnuFlowProfileUsage"`
}

type OnuFlowProfileList struct {
	Entry []*OnuFlowProfile
}

// NewOnuFlowProfile returns an OnuFlowProfile struct with default values
func NewOnuFlowProfile(name string) *OnuFlowProfile {
	p := &OnuFlowProfile{
		Name:                name,
		MatchUsCVlanIDRange: empty,
		MatchUsCPcp:         -1,
		UsCdr:               128,
		UsPdr:               1244160,
		UsFlowPriority:      0,
		DsFlowPriority:      0,
	}

	return p
}

// GetName returns the name of the OnuFlowProfile
func (p *OnuFlowProfile) GetName() string {
	return p.Name
}

// Copy returns a copy of the profile object with a new name and Usage set to 2
func (p *OnuFlowProfile) Copy(newName string) (*OnuFlowProfile, error) {
	if p.Name == newName {
		return nil, ErrExists
	}
	np := p
	np.Name = newName
	np.Usage = 2
	return np, nil
}

// GetMatchUsCVlanIDRange returns the values set to match with Customer VLAN ID in the OnuFlowProfile
func (p *OnuFlowProfile) GetMatchUsCVlanIDRange() []int {
	list, err := getVlanFromB64(p.MatchUsCVlanIDRange)
	if err != nil {
		log.Printf("onuFlowProfile: %v\n", err)
		return []int{}
	}
	return list
}

// SetMatchUsCVlanIDRange allows applying an int slice to the OnuFlowProfile to be used as Us Match C-VID values
func (p *OnuFlowProfile) SetMatchUsCVlanIDRange(vlans []int) (err error) {
	p.MatchUsCVlanIDRange, err = getB64FromVlan(vlans)
	return err
}

// SetMatchUsCVlanIDRangeFromString wraps the int slice format around a string of comma-separated VLANs
func (p *OnuFlowProfile) SetMatchUsCVlanIDRangeFromString(vlans string) error {
	list, err := generateVlanList(vlans)
	if err != nil || len(list) == 0 {
		return ErrNotInput
	}
	return p.SetMatchUsCVlanIDRange(list)
}

var OnuFlowProfileHeaders = []string{
	"Name",
	"MatchUsC-VidRange",
	"MatchUsCPcp",
	"UsFlowPriority",
}

// ListEssentialParams returns a map of the essential OnuFlowProfile parameters
func (p *OnuFlowProfile) ListEssentialParams() map[string]interface{} {
	var EssentialOnuFlowProfile = map[string]interface{}{
		OnuFlowProfileHeaders[0]: p.GetName(),
		OnuFlowProfileHeaders[1]: p.GetMatchUsCVlanIDRange(),
		OnuFlowProfileHeaders[2]: p.MatchUsCPcp,
		OnuFlowProfileHeaders[3]: p.UsFlowPriority,
	}

	return EssentialOnuFlowProfile
}

// Tabwrite displays the essential information of OnuFlowProfile in organized columns
func (p *OnuFlowProfile) Tabwrite() {
	l := p.ListEssentialParams()
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	for _, v := range OnuFlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range OnuFlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range OnuFlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", l[v])
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range OnuFlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	tw.Flush()
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (p *OnuFlowProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", data
	}
	return p.Name, data
}

// Separate is a method to maintain backward-compatability
func (ofpl *OnuFlowProfileList) Separate() []*OnuFlowProfile {
	var entry *OnuFlowProfile
	var list []*OnuFlowProfile
	for _, e := range ofpl.Entry {
		entry = e
		list = append(list, entry)
	}
	return list
}

// Tabwrite displays the essential information of a list of Flow Profiles in organized columns
func (ofpl *OnuFlowProfileList) Tabwrite() {
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range OnuFlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuFlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	ofps := ofpl.Separate()
	for _, ofp := range ofps {
		// first get the data as a map
		l := ofp.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range OnuFlowProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuFlowProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}
