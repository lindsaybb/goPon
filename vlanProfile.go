package gopon

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

// VlanProfile is a collection of parameters for creation of a Vlan profile
type VlanProfile struct {
	Name               string `json:"msanVlanProfileName"`
	CVid               string `json:"msanVlanProfileCVid"`
	CVidNative         int    `json:"msanVlanProfileCVidNative"`
	CVidRemark         int    `json:"msanVlanProfileCVidRemark"`
	SVid               int    `json:"msanVlanProfileSVid"`
	SEtherType         int    `json:"msanVlanProfileSEtherType"`
	NetworkPortCTag    int    `json:"msanVlanProfileNetworkPortCTag"`
	CVidExternal       int    `json:"msanVlanProfileCVidExternal"`
	CVidNativeExternal int    `json:"msanVlanProfileCVidNativeExternal"`
	CVidRemarkExternal int    `json:"msanVlanProfileCVidRemarkExternal"`
	SVidExternal       int    `json:"msanVlanProfileSVidExternal"`
	Usage              int    `json:"msanVlanProfileUsage"`
}

type VlanProfileList struct {
	Entry []*VlanProfile
}

// NewVlanProfile returns an empty, initialzed struct to be populated
func NewVlanProfile(name string) *VlanProfile {
	p := &VlanProfile{
		Name:               name,
		CVid:               empty,
		CVidNative:         -1,
		CVidRemark:         -1,
		SVid:               -1,
		SEtherType:         34984,
		NetworkPortCTag:    1,
		CVidExternal:       2,
		CVidNativeExternal: 2,
		CVidRemarkExternal: 2,
		SVidExternal:       2,
	}

	return p
}

// GetName returns the name of the VlanProfile
func (p *VlanProfile) GetName() string {
	return p.Name
}

// Copy returns a copy of the profile object with a new name and Usage set to 2
func (p *VlanProfile) Copy(newName string) (*VlanProfile, error) {
	if p.Name == newName {
		return nil, ErrExists
	}
	np := p
	np.Name = newName
	np.Usage = 2
	return np, nil
}

// GetCVid returns the values set as Customer VLAN ID in the VlanProfile
func (p *VlanProfile) GetCVid() []int {
	list, err := getVlanFromB64(p.CVid)
	if err != nil {
		log.Printf("vlanProfile: %v\n", err)
		return []int{}
	}
	return list
}

// SetCVid allows applying an int slice to the VlanProfile to be used as C-VID values
func (p *VlanProfile) SetCVid(vlans []int) (err error) {
	p.CVid, err = getB64FromVlan(vlans)
	return err
}

// SetCVidFromString wraps the int slice format around a string of comma-separated VLANs
func (p *VlanProfile) SetCVidFromString(vlans string) error {
	list, err := generateVlanList(vlans)
	if err != nil || len(list) == 0 {
		return ErrNotInput
	}
	return p.SetCVid(list)
}

var VlanProfileHeaders = []string{
	"Name",
	"C-Vid",
	"C-Vid Native",
	"S-Vid",
	"S-Ethertype",
}

// ListEssentialParams returns a map of the essential VlanProfile parameters
func (p *VlanProfile) ListEssentialParams() map[string]interface{} {
	var EssentialVlanProfile = map[string]interface{}{
		VlanProfileHeaders[0]: p.GetName(),
		VlanProfileHeaders[1]: p.GetCVid(),
		VlanProfileHeaders[2]: p.CVidNative,
		VlanProfileHeaders[3]: p.SVid,
		VlanProfileHeaders[4]: p.SEtherType,
	}

	return EssentialVlanProfile
}

// Tabwrite displays the essential information of VlanProfile in organized columns
func (p *VlanProfile) Tabwrite() {
	l := p.ListEssentialParams()
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	for _, v := range VlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range VlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range VlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", l[v])
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range VlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	tw.Flush()
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (p *VlanProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", data
	}
	return p.Name, data
}

// Separate is a method to maintain backward-compatability
func (vpl *VlanProfileList) Separate() []*VlanProfile {
	var entry *VlanProfile
	var list []*VlanProfile
	for _, e := range vpl.Entry {
		entry = e
		list = append(list, entry)
	}
	return list
}

// Tabwrite displays the essential information of a list of Flow Profiles in organized columns
func (vpl *VlanProfileList) Tabwrite() {
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range VlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range VlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	vps := vpl.Separate()
	for _, vp := range vps {
		// first get the data as a map
		l := vp.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range VlanProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range VlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}

/*
func generateVlanProfile(h []string) ([]byte) {
	//vlanProfiles: []string{"RestTest", "CVidList"},
	//vlanProfiles: []string{"RestTest", "110"},
	cvids, err := generateVlanList(h[1])
	if err != nil {
		cvids = []int{1}
	}
	cvidString := getB64FromVlan(cvids)
	w := new(VlanProfile)
	w.MsanVlanProfileName = h[0]
	w.MsanVlanProfileCVid = cvidString
	w.MsanVlanProfileCVidNative = -1 //, _ = strconv.Atoi(h[2])	//"-1"
	w.MsanVlanProfileCVidRemark = -1 //, _ = strconv.Atoi(h[3])	//"-1"
	w.MsanVlanProfileSVid = -1 //, _ = strconv.Atoi(h[4])	//"-1"
	w.MsanVlanProfileSEtherType = 34984 //, _ = strconv.Atoi(h[5])	//"34984"
	w.MsanVlanProfileNetworkPortCTag = 1 //, _ = strconv.Atoi(h[6])	//"1"
	w.MsanVlanProfileCVidExternal = 2 //, _ = strconv.Atoi(h[7])	//"2"
	w.MsanVlanProfileCVidNativeExternal = 2 //, _ = strconv.Atoi(h[8])	//"2"
	w.MsanVlanProfileCVidRemarkExternal = 2 //, _ = strconv.Atoi(h[9])	//"2"
	w.MsanVlanProfileSVidExternal = 2 //, _ = strconv.Atoi(h[10])	//"2"

	//fmt.Println(w)
	//t := new(IskratelMsan)
	//t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanVlanProfileTable.MsanVlanProfileEntry = append(t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanVlanProfileTable.MsanVlanProfileEntry, *w)

	data, err := json.Marshal(w)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}
*/
