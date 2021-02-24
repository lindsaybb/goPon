package gopon

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// OnuTcontProfile contains the profile for defining the speeds and priority of services
type OnuTcontProfile struct {
	Name            string `json:"msanOnuTcontProfileName"`
	TcontID         int    `json:"msanOnuTcontProfileTcontId"`
	TcontType       int    `json:"msanOnuTcontProfileTcontType"`
	FixedDataRate   int    `json:"msanOnuTcontProfileFixedDataRate"`
	AssuredDataRate int    `json:"msanOnuTcontProfileAssuredDataRate"`
	MaxDataRate     int    `json:"msanOnuTcontProfileMaxDataRate"`
	Usage           int    `json:"msanOnuTcontProfileUsage"`
}

type OnuTcontProfileList struct {
	Entry []*OnuTcontProfile
}

// NewOnuTcontProfile returns an OnuTcontProfile struct with default values
func NewOnuTcontProfile(name string) *OnuTcontProfile {
	p := &OnuTcontProfile{
		Name:            name,
		TcontID:         1,
		TcontType:       5,
		FixedDataRate:   0,
		AssuredDataRate: 0,
		MaxDataRate:     256,
	}
	return p
}

// GetName returns the name of the OnuTcontProfile
func (p *OnuTcontProfile) GetName() string {
	return p.Name
}

// Copy returns a copy of the profile object with a new name and Usage set to 2
func (p *OnuTcontProfile) Copy(newName string) (*OnuTcontProfile, error) {
	if p.Name == newName {
		return nil, ErrExists
	}
	np := p
	np.Name = newName
	np.Usage = 2
	return np, nil
}

// GenerateTcontName returns a string of the suggested naming convention based on profile details
func (p *OnuTcontProfile) GenerateTcontName() string {
	// standard naming convention allows scalability and ease of use
	// only called if desired, can name the profiles anything
	max := formatKbits(p.MaxDataRate)
	id := toString(p.TcontID)
	ctype := toString(p.TcontType)

	name := "T" + ctype + "I" + id
	switch p.TcontType {
	case 1:
		name += "F__"
	case 2:
		name += "_A_"
	case 3:
		name += "_AM"
	case 4:
		name += "__M"
	case 5:
		if p.FixedDataRate > 256 {
			name += "F"
		} else {
			name += "_"
		}
		if p.AssuredDataRate >= 256 {
			name += "A"
		} else {
			name += "_"
		}
		name += "M"
	}
	name += "-" + max
	return name
}

// GetTcontDescription returns a string of helpful information about TConts, including parameters specific to the selected profile
func (p *OnuTcontProfile) GetTcontDescription() string {
	// this section will interpret the Tcont package into a short blurb
	// useful for understanding how changes to this profile affect services
	// this can extrapolate from the standard naming convention
	name := p.GenerateTcontName()
	var descr = []string{
		"Transmission Containers are responsible for negotiating customer services in a multiaccess architecture",
		"Since GPON is Asymmetrical, where the Downlink frame is Broadcast/Unicast, T-Conts only control upstream traffic",
		"There are some tricks to using T-Conts that will become more familiar throughout this exercise",
		"T-Cont Types are a value from 1-5 identifying the handling of committed and burst rates",
		"Type 1 allows setting a Fixed rate that is consumed whether or not the customer is using it, useful for TDM emulation",
		"Type 2 allows setting as Assured rate that is not consumed when not in use, but has priority handling over best-effort services otherwise",
		"Type 3 allows setting an Assured and Max rate, so a portion of the Max is given priority handling, and a portion is best-effort",
		"Type 4 allows setting a Max rate for best-effort handling",
		"Type 5 is a special T-Cont that allows setting Fixed, Assured and Maximum rates",
		"The default T-Cont Type if not set is 5",
		"T-Cont IDs are a value from 1-6 that allow stacking multiple T-Conts on the same ONU",
		"T-Cont IDs cannot overlap on the same ONU",
		"Best practices recommend a structure where the same 'types' of services are given the same IDs, as it is less likely for these to be applied to the same ONU",
		"An example would be always specifying CWMP as ID 6, Internet data as type 2, VoIP as type 3, and IPTV as type 4",
		"The default T-Cont ID if not set is 1",
		"T-Conts can be flexibly named, however due to their function a standard naming convention is recommended",
	}
	desc := strings.Join(descr, ". ")
	name = "Based on the T-Cont in this profile, the recommended name is:" + name
	return desc + name
}

// SetTcontType allows setting the TCont type between 1 and 5, carries over previously set rates if applicable
func (p *OnuTcontProfile) SetTcontType(i int) {
	// before changing, preserve original values if exist
	var tmpF, tmpA, tmpM int
	tmpF = p.FixedDataRate
	tmpA = p.AssuredDataRate
	tmpM = p.MaxDataRate
	// best practice would be to set this before setting rates
	if i < 1 || i > 5 {
		i = 5
	}
	p.TcontType = i
	p.SetFAM(tmpF, tmpA, tmpM)
}

// SetTcontID allows setting the TCont Id between 1 and 6
func (p *OnuTcontProfile) SetTcontID(i int) {
	if i < 1 || i > 6 {
		i = 1
	}
	p.TcontType = i
}

// SetFAM allows setting Fixed, Assured and Max Data rates at once, pass 0 if can't or don't want to set value
func (p *OnuTcontProfile) SetFAM(fixed, assured, max int) {
	if p.CanSetFixed() {
		p.FixedDataRate = fixed
	}
	if p.CanSetAssured() {
		p.AssuredDataRate = assured
	}
	if p.CanSetMax() {
		p.MaxDataRate = max
	}
}

// CanSetMax checks if TCont Type allows setting Max value and returns Bool
func (p *OnuTcontProfile) CanSetMax() bool {
	return p.TcontType == 5 || p.TcontType == 4 || p.TcontType == 3
}

// SetMaxRate allows setting a Max rate between 256 and 1244160 (GPON)
func (p *OnuTcontProfile) SetMaxRate(i int) {
	if i < 256 {
		i = 256
	}
	if i > 1244160 {
		i = 1244160
	}
	p.MaxDataRate = i
}

// CanSetAssured checks if TCont Type allows setting Assured value and returns Bool
func (p *OnuTcontProfile) CanSetAssured() bool {
	return p.TcontType == 5 || p.TcontType == 3 || p.TcontType == 2
}

// SetAssuredRate allows setting an Assured rate between 256 and 1244160 (GPON)
func (p *OnuTcontProfile) SetAssuredRate(i int) {
	if i < 256 {
		i = 256
	}
	if i > 1244160 {
		i = 1244160
	}
	p.AssuredDataRate = i
}

// CanSetFixed checks if TCont Type allows setting Fixed value and returns Bool
func (p *OnuTcontProfile) CanSetFixed() bool {
	return p.TcontType == 5 || p.TcontType == 1
}

// SetFixedRate allows setting a Fixed rate between 256 and 1244160 (GPON ver 2.0.0)
func (p *OnuTcontProfile) SetFixedRate(i int) {
	if i < 256 {
		i = 256
	}
	// recommendation to not allow setting higher than 667000
	if i > 1244160 {
		i = 1244160
	}
	p.FixedDataRate = i
}

var OnuTcontProfileHeaders = []string{
	"Name",
	"Description",
	"ID",
	"Type",
	"Fixed",
	"Assured",
	"Max",
}

// ListEssentialParams returns a map of the essential OnuFlowProfile parameters
func (p *OnuTcontProfile) ListEssentialParams() map[string]interface{} {
	var EssentialOnuTcontProfile = map[string]interface{}{
		OnuTcontProfileHeaders[0]: p.GetName(),
		OnuTcontProfileHeaders[1]: p.GenerateTcontName(),
		OnuTcontProfileHeaders[2]: p.TcontID,
		OnuTcontProfileHeaders[3]: p.TcontType,
		OnuTcontProfileHeaders[4]: formatKbits(p.FixedDataRate),
		OnuTcontProfileHeaders[5]: formatKbits(p.AssuredDataRate),
		OnuTcontProfileHeaders[6]: formatKbits(p.MaxDataRate),
	}

	return EssentialOnuTcontProfile
}

// Tabwrite displays the essential information of OnuTcontProfile in organized columns
func (p *OnuTcontProfile) Tabwrite() {
	l := p.ListEssentialParams()
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	for _, v := range OnuTcontProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range OnuTcontProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range OnuTcontProfileHeaders {
		fmt.Fprintf(tw, "%v\t", l[v])
	}
	fmt.Fprintf(tw, "\n")
	for _, v := range OnuTcontProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	tw.Flush()
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (p *OnuTcontProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", data
	}
	return p.Name, data
}

// Separate is a method to maintain backward-compatability
func (otpl *OnuTcontProfileList) Separate() []*OnuTcontProfile {
	var entry *OnuTcontProfile
	var list []*OnuTcontProfile
	for _, e := range otpl.Entry {
		entry = e
		list = append(list, entry)
	}
	return list
}

// Tabwrite displays the essential information of a list of Flow Profiles in organized columns
func (otpl *OnuTcontProfileList) Tabwrite() {
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range OnuTcontProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuTcontProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	otps := otpl.Separate()
	for _, otp := range otps {
		// first get the data as a map
		l := otp.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range OnuTcontProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuTcontProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}
