package gopon

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

type OnuIgmpProfile struct {
	Name                  string `json:"msanOnuMulticastProfileName"`
	IgmpMode              int    `json:"msanOnuMulticastProfileIgmpMode"`
	IgmpProxy             int    `json:"msanOnuMulticastProfileIgmpProxy"`
	IgmpSnoopingFastLeave int    `json:"msanOnuMulticastProfileIgmpSnoopingFastLeave"`
	UsIgmpTciVlanID       int    `json:"msanOnuMulticastProfileUsIgmpTciVlanId"`
	UsIgmpTciPcpValue     int    `json:"msanOnuMulticastProfileUsIgmpTciPcpValue"`
	UsIgmpTciCtrlMode     int    `json:"msanOnuMulticastProfileUsIgmpTciCtrlMode"`
	DsVlanTagging         int    `json:"msanOnuMulticastProfileDsVlanTagging"`
	DsGemPort             int    `json:"msanOnuMulticastProfileDsGemPort"`
	Usage                 int    `json:"msanOnuMulticastProfileUsage"`
}

type OnuIgmpProfileList struct {
	Entry []*OnuIgmpProfile
}

func NewOnuIgmpProfile(name string) *OnuIgmpProfile {
	p := &OnuIgmpProfile{
		Name:                  name,
		IgmpMode:              1,
		IgmpProxy:             2,
		IgmpSnoopingFastLeave: 1,
		UsIgmpTciVlanID:       0,
		UsIgmpTciPcpValue:     0,
		UsIgmpTciCtrlMode:     5,
		DsVlanTagging:         2,
		DsGemPort:             4000,
	}
	return p
}

func (p *OnuIgmpProfile) GetName() string {
	return p.Name
}

// Copy returns a copy of the profile object with a new name and Usage set to 2
func (p *OnuIgmpProfile) Copy(newName string) (*OnuIgmpProfile, error) {
	if p.Name == newName {
		return nil, ErrExists
	}
	np := p
	np.Name = newName
	np.Usage = 2
	return np, nil
}

func (p *OnuIgmpProfile) GetMode() string {
	if p.IgmpMode == 1 {
		return "Flooding"
	}
	if p.IgmpMode == 2 {
		return "Snooping"
	}
	return "Unknown"
}

func (p *OnuIgmpProfile) SetMode(value int) {
	if value < 1 {
		value = 1
	}
	if value > 2 {
		value = 2
	}
	p.IgmpMode = value
}

func (p *OnuIgmpProfile) GetProxy() bool {
	// 1 is enabled, default is 2 disabled
	return p.IgmpProxy == 1
}

func (p *OnuIgmpProfile) SetProxy(value bool) {
	if value {
		p.IgmpProxy = 1
	} else {
		p.IgmpProxy = 2
	}
}

func (p *OnuIgmpProfile) GetFastLeave() bool {
	return p.IgmpSnoopingFastLeave == 1
}

func (p *OnuIgmpProfile) SetFastLeave() {
	p.IgmpSnoopingFastLeave = 1
}

func (p *OnuIgmpProfile) GetGEMPort() int {
	return p.DsGemPort
}

func (p *OnuIgmpProfile) SetGEMPort(value int) {
	if value < 3800 || value >= 4000 {
		value = 3999
	}
	p.DsGemPort = value
}

func (p *OnuIgmpProfile) GetIgmpSnooping() bool {
	return p.IgmpMode == 2
}

func (p *OnuIgmpProfile) SetIgmpSnooping() {
	p.IgmpMode = 2
}

func (p *OnuIgmpProfile) GetIgmpFlooding() bool {
	return p.IgmpMode == 1
}

func (p *OnuIgmpProfile) SetIgmpFlooding() {
	p.IgmpMode = 1
}

func (p *OnuIgmpProfile) GetIgmpSnoopingFastLeave() bool {
	return p.IgmpSnoopingFastLeave == 1
}

func (p *OnuIgmpProfile) SetIgmpSnoopingFastLeave(state bool) {
	if state {
		p.IgmpSnoopingFastLeave = 1
	} else {
		p.IgmpSnoopingFastLeave = 0
	}
}

func (p *OnuIgmpProfile) GetUsTci() [3]int {
	return [3]int{p.UsIgmpTciVlanID, p.UsIgmpTciPcpValue, p.UsIgmpTciCtrlMode}
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (p *OnuIgmpProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", data
	}
	return p.Name, data
}

/*
Name: name,
		IgmpMode: 1,
		IgmpProxy: 2,
		IgmpSnoopingFastLeave: 1,
		UsIgmpTciVlanID: 0,
		UsIgmpTciPcpValue: 0,
		UsIgmpTciCtrlMode: 5,
		DsVlanTagging: 2,
		DsGemPort: 4000,
*/
var OnuIgmpProfileHeaders = []string{
	"Name",
	"Mode",
	"Proxy",
	"Fast-Leave",
	"UsTci",
	"DsGem",
}

// ListEssentialParams returns a map of the essential OnuIgmpProfile parameters
func (p *OnuIgmpProfile) ListEssentialParams() map[string]interface{} {
	var EssentialOnuIgmpProfile = map[string]interface{}{
		OnuIgmpProfileHeaders[0]: p.GetName(),
		OnuIgmpProfileHeaders[1]: p.GetMode(),
		OnuIgmpProfileHeaders[2]: p.GetProxy(),
		OnuIgmpProfileHeaders[3]: p.GetFastLeave(),
		OnuIgmpProfileHeaders[4]: p.GetUsTci(),
		OnuIgmpProfileHeaders[5]: p.DsGemPort,
	}
	// I want all of these Bools to return strings of "Enabled/Disabled"
	return EssentialOnuIgmpProfile
}

// Tabwrite displays the essential information of VlanProfile in organized columns
func (p *OnuIgmpProfile) Tabwrite() {
	l := p.ListEssentialParams()
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	for _, v := range OnuIgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range OnuIgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range OnuIgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", l[v])
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range OnuIgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	tw.Flush()
}

// Separate is a method to maintain backward-compatability
func (oipl *OnuIgmpProfileList) Separate() []*OnuIgmpProfile {
	var entry *OnuIgmpProfile
	var list []*OnuIgmpProfile
	for _, e := range oipl.Entry {
		entry = e
		list = append(list, entry)
	}
	return list
}

// Tabwrite displays the essential information of a list of Flow Profiles in organized columns
func (oipl *OnuIgmpProfileList) Tabwrite() {
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range OnuIgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuIgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	oips := oipl.Separate()
	for _, oip := range oips {
		// first get the data as a map
		l := oip.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range OnuIgmpProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuIgmpProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}

/*
// Optional to specify the Proxy IP Address when enabling
func (p *OnuIgmpProfile) SetIgmpProxy(state bool, address ...string) {
	if state {
		p.IgmpProxy = 1
		if len(address) != 0 {
			p.IgmpProxyIPAddress = address[0]
		}
	} else {
		p.IgmpSnooping = 0
	}
	return
}
*/
/*
// Support ticket open for why TCI VLAN ID cannot use default value of 0
func generateOnuIgmpProfile(h []string) ([]byte) {
	//onuIgmpProfiles: []string{"Name", "IgmpMode", "IgmpProxy", "SnoopingFastLeave", "DsGemPort"},
	//onuIgmpProfiles: []string{"RestTest", "2", "2", "1", "3998"},
	w := new(OnuIgmpProfile)
	w.MsanOnuMulticastProfileName = h[0]
	w.MsanOnuMulticastProfileIgmpMode = 2 // snooping, default is flooding [1:flood, 2:snoop]
	w.MsanOnuMulticastProfileIgmpProxy = 2 	// disabled, def [1:enable]
	w.MsanOnuMulticastProfileIgmpSnoopingFastLeave = 1 	// enabled, def [2:disable]
	w.MsanOnuMulticastProfileUsIgmpTciVlanID = 1 	//# 0 is supposed to be not defined/def (tag control information) but isn't allowed, 1 works, int32 value 0..4094
	w.MsanOnuMulticastProfileUsIgmpTciPcpValue = 0 	// not defined, def
	w.MsanOnuMulticastProfileUsIgmpTciCtrlMode = 5 	// smart, def [1:transparent, 2:add, 3:replace, 4:replace-vid-only, 5:smart]
	w.MsanOnuMulticastProfileDsVlanTagging = 2 	// disable, def
	w.MsanOnuMulticastProfileDsGemPort = 3998 	// 4000 is default, range 3800-4000, cannot overlap on same device

	//fmt.Println(w)
	//t := new(IskratelMsan)
	//t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuMulticastProfileTable.MsanOnuMulticastProfileEntry = append(t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuMulticastProfileTable.MsanOnuMulticastProfileEntry, *w)
	data, err := json.Marshal(w)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}
*/
