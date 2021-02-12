package gopon

// this file acts as a handler between the larger data structure and the individual profiles
// combination and separation of profiles is done here, using the master struct

// should implement a cache method so repetitive calls for the same element are not duplicated
// but for now, for each new call, a full query will be made and parsed directly

// the operations related to GET, SET and DELETE must reference a HOST config
// all of these individual sub-,methods involve separating data out of a bulk-request for the entire data structure
// that's not how I've handled it in the past, need to specify the endpoint with each GET
// all of these methods begin with the root data structure then fragment their own sub-set, coming back to resolve it completely
// this way of handling the data ensures the proper level of nesting, while dealing with simpler objects

// change all fmt prints to logwriter

// handling of OnuVlanProfile and subset of OnuVlanRules needs some special attention

// create a conveince method for tabwriting all service profiles as a list

type IskratelMsan struct {
	ISKRATELMSANMIB struct {
		ISKRATELMSANMIB struct {
			MsanServiceProfileTable struct {
				MsanServiceProfileEntry []ServiceProfile `json:"msanServiceProfileEntry"`
			} `json:"msanServiceProfileTable"`
			MsanServiceFlowProfileTable struct {
				MsanServiceFlowProfileEntry []FlowProfile `json:"msanServiceFlowProfileEntry"`
			} `json:"msanServiceFlowProfileTable"`
			MsanVlanProfileTable struct {
				MsanVlanProfileEntry []VlanProfile `json:"msanVlanProfileEntry"`
			} `json:"msanVlanProfileTable"`
			MsanMulticastProfileTable struct {
				MsanMulticastProfileEntry []IgmpProfile `json:"msanMulticastProfileEntry"`
			} `json:"msanMulticastProfileTable"`
			MsanSecurityProfileTable struct {
				MsanSecurityProfileEntry []SecurityProfile `json:"msanSecurityProfileEntry"`
			} `json:"msanSecurityProfileTable"`
			MsanOnuFlowProfileTable struct {
				MsanOnuFlowProfileEntry []OnuFlowProfile `json:"msanOnuFlowProfileEntry"`
			} `json:"msanOnuFlowProfileTable"`
			MsanOnuTcontProfileTable struct {
				MsanOnuTcontProfileEntry []OnuTcontProfile `json:"msanOnuTcontProfileEntry"`
			} `json:"msanOnuTcontProfileTable"`
			MsanOnuVlanProfileTable struct {
				MsanOnuVlanProfileEntry []OnuVlanProfile `json:"msanOnuVlanProfileEntry"`
			} `json:"msanOnuVlanProfileTable"`
			MsanOnuVlanProfileRuleTable struct {
				MsanOnuVlanProfileRuleEntry []OnuVlanRule `json:"msanOnuVlanProfileRuleEntry"`
			} `json:"msanOnuVlanProfileRuleTable"`
			MsanOnuMulticastProfileTable struct {
				MsanOnuMulticastProfileEntry []OnuIgmpProfile `json:"msanOnuMulticastProfileEntry"`
			} `json:"msanOnuMulticastProfileTable"`
			MsanL2CpProfileTable struct {
				MsanL2CpProfileEntry []L2cpProfile `json:"msanL2cpProfileEntry"`
			} `json:"msanL2cpProfileTable"`
			MsanOnuBlackListTable struct {
				MsanOnuBlackListEntry []OnuBlacklist `json:"msanOnuBlackListEntry"`
			} `json:"msanOnuBlackListTable"`
			MsanOnuCfgTable struct {
				MsanOnuCfgEntry []OnuConfig `json:"msanOnuCfgEntry"`
			} `json:"msanOnuCfgTable"`
			MsanServicePortProfileTable struct {
				MsanServicePortProfileEntry []OnuProfile `json:"msanServicePortProfileEntry"`
			} `json:"msanServicePortProfileTable"`
		} `json:"ISKRATEL-MSAN-MIB"`
	} `json:"ISKRATEL-MSAN-MIB:"`
}

// NewIskratelMsan sets up a data structure for the specified Host
func NewIskratelMsan() *IskratelMsan {
	var t = &IskratelMsan{}

	return t
}

/*
// GetServiceProfiles performs a Get Request to the t.Host and returns a list of the ServiceProfile struct
func (t *IskratelMsan) GetServiceProfiles() ([]*ServiceProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, serviceProfiles)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(rawJson))
	err = json.Unmarshal(rawJson, t)
	if err != nil {
		return nil, err
	}
	var list []*ServiceProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanServiceProfileTable.MsanServiceProfileEntry {
		fmt.Println(&v)
		list = append(list, &v)
	}
	return list, nil
}

// GetServiceProfileByName is a helper method that returns a single ServiceProfile struct by name, if exists
func (t *IskratelMsan) GetServiceProfileByName(name string) (*ServiceProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetServiceProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteServiceProfile removes the named ServiceProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteServiceProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetServiceProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, serviceProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteServiceProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostServiceProfile performs a Post request to t.Host containing serialized data from a ServiceProfile struct, if the name is not already used
func (t *IskratelMsan) PostServiceProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetServiceProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The ServiceProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, serviceProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostServiceProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// GetFlowProfiles performs a Get Request to the t.Host and returns a list of the FlowProfile struct
func (t *IskratelMsan) GetFlowProfiles() ([]*FlowProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, flowProfiles)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*FlowProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanServiceFlowProfileTable.MsanServiceFlowProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetFlowProfileByName is a helper method that returns a single FlowProfile struct by name, if exists
func (t *IskratelMsan) GetFlowProfileByName(name string) (*FlowProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetFlowProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteFlowProfile removes the named FlowProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteFlowProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetFlowProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, flowProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteFlowProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostFlowProfile performs a Post request to t.Host containing serialized data from a FlowProfile struct, if the name is not already used
func (t *IskratelMsan) PostFlowProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetFlowProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The FlowProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, flowProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostFlowProfile:", status)
	// check for ErrNotStatusOk
	return nil

}

// GetVlanProfiles performs a Get Request to the t.Host and returns a list of the VlanProfile struct
func (t *IskratelMsan) GetVlanProfiles() ([]*VlanProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, vlanProfiles)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*VlanProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanVlanProfileTable.MsanVlanProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetVlanProfileByName is a helper method that returns a single VlanProfile struct by name, if exists
func (t *IskratelMsan) GetVlanProfileByName(name string) (*VlanProfile, error) {
	list, err := t.GetVlanProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteVlanProfile removes the named VlanProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteVlanProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetVlanProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, vlanProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteVlanProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostVlanProfile performs a Post request to t.Host containing serialized data from a VlanProfile struct, if the name is not already used
func (t *IskratelMsan) PostVlanProfile(name string, data []byte) error {
	// check if name is already in use, or let NE handle the conflict itself... ?
	_, err := t.GetVlanProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The VlanProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, vlanProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostVlanProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// GetMulticastProfiles performs a Get Request to the t.Host and returns a list of the IgmpProfile struct
func (t *IskratelMsan) GetMulticastProfiles() ([]*IgmpProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, igmpProfiles)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*IgmpProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanMulticastProfileTable.MsanMulticastProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetMulticastProfileByName is a helper method that returns a single IgmpProfile struct by name, if exists
func (t *IskratelMsan) GetMulticastProfileByName(name string) (*IgmpProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetMulticastProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteMulticastProfile removes the named IgmpProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteMulticastProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetMulticastProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, igmpProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteMulticastProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostMulticastProfile performs a Post request to t.Host containing serialized data from a IgmpProfile struct, if the name is not already used
func (t *IskratelMsan) PostMulticastProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetMulticastProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The IgmpProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, igmpProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostMulticastProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// GetSecurityProfiles performs a Get Request to the t.Host and returns a list of the SecurityProfile struct
func (t *IskratelMsan) GetSecurityProfiles() ([]*SecurityProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, securityProfiles)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*SecurityProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanSecurityProfileTable.MsanSecurityProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetSecurityProfileByName is a helper method that returns a single SecurityProfile struct by name, if exists
func (t *IskratelMsan) GetSecurityProfileByName(name string) (*SecurityProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetSecurityProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteSecurityProfile removes the named SecurityProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteSecurityProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetSecurityProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, securityProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteSecurityProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostSecurityProfile performs a Post request to t.Host containing serialized data from a SecurityProfile struct, if the name is not already used
func (t *IskratelMsan) PostSecurityProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetSecurityProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The SecurityProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, securityProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostSecurityProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// GetOnuFlowProfiles performs a Get Request to the t.Host and returns a list of the OnuFlowProfile struct
func (t *IskratelMsan) GetOnuFlowProfiles() ([]*OnuFlowProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, onuFlowProfiles)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*OnuFlowProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuFlowProfileTable.MsanOnuFlowProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetOnuFlowProfileByName is a helper method that returns a single OnuFlowProfile struct by name, if exists
func (t *IskratelMsan) GetOnuFlowProfileByName(name string) (*OnuFlowProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetOnuFlowProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteOnuFlowProfile removes the named OnuFlowProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteOnuFlowProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetOnuFlowProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, onuFlowProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuFlowProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostOnuFlowProfile performs a Post request to t.Host containing serialized data from a OnuFlowProfile struct, if the name is not already used
func (t *IskratelMsan) PostOnuFlowProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetOnuFlowProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The OnuFlowProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, onuFlowProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuFlowProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// GetOnuTcontProfiles performs a Get Request to the t.Host and returns a list of the OnuTcontProfile struct
func (t *IskratelMsan) GetOnuTcontProfiles() ([]*OnuTcontProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, onuTcontProfiles)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*OnuTcontProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuTcontProfileTable.MsanOnuTcontProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetOnuTcontProfileByName is a helper method that returns a single OnuTcontProfile struct by name, if exists
func (t *IskratelMsan) GetOnuTcontProfileByName(name string) (*OnuTcontProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetOnuTcontProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteOnuTcontProfile removes the named OnuTcontProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteOnuTcontProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetOnuTcontProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, onuTcontProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuTcontProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostOnuTcontProfile performs a Post request to t.Host containing serialized data from a OnuTcontProfile struct, if the name is not already used
func (t *IskratelMsan) PostOnuTcontProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetOnuTcontProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The OnuTcontProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, onuTcontProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuTcontProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// GetOnuVlanProfiles performs a Get Request to the t.Host and returns a list of the OnuVlanProfile struct
func (t *IskratelMsan) GetOnuVlanProfiles() ([]*OnuVlanProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, onuVlanProfiles)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*OnuVlanProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuVlanProfileTable.MsanOnuVlanProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetOnuVlanProfileByName is a helper method that returns a single OnuVlanProfile struct by name, if exists
func (t *IskratelMsan) GetOnuVlanProfileByName(name string) (*OnuVlanProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetOnuVlanProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteOnuVlanProfile removes the named OnuVlanProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteOnuVlanProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetOnuVlanProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, onuVlanProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuVlanProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostOnuVlanProfile performs a Post request to t.Host containing serialized data from a OnuVlanProfile struct, if the name is not already used
func (t *IskratelMsan) PostOnuVlanProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetOnuVlanProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The OnuVlanProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, onuVlanProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuVlanProfile:", status)
	// check for ErrNotStatusOk
	return nil
}
*/
/*
// GetOnuVlanRules performs a Get Request to the t.Host and returns a list of the OnuVlanRule struct
func (t *IskratelMsan) GetOnuVlanRules() ([]*OnuVlanRule, error) {
	rawJson, err := RestGetProfiles(t.Host, onuVlanRules)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*OnuVlanRule
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuVlanProfileRuleTable.MsanOnuVlanProfileRuleEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetOnuVlanRuleByName is a helper method that returns a single OnuVlanRule struct by name, if exists
func (t *IskratelMsan) GetOnuVlanRuleByName(name string) (*OnuVlanRule, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetOnuVlanRules()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteOnuVlanRule removes the named OnuVlanRule from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteOnuVlanRule(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetOnuVlanRuleByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, onuVlanRules, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuVlanRule:", status)
	// check for ErrNotStatusOk
	return nil
}

// CANNOT POST RULES INDEPENDENT OF THE ONUVLANPROFILE, THESE FUNCTIONS ARE A SUBSET OF THAT PROFILE'S OPERATIONS
// PostOnuVlanRule performs a Post request to t.Host containing serialized data from a OnuVlanRule struct, if the name is not already used
func (t *IskratelMsan) PostOnuVlanRule(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetOnuVlanRuleByName(name)
	if err == nil {
		return ErrExists
	}
	// The OnuVlanRule has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, onuVlanRules, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuVlanRule:", status)
	// check for ErrNotStatusOk
	return nil
}
*/
/*
// GetOnuMulticastProfiles performs a Get Request to the t.Host and returns a list of the OnuIgmpProfile struct
func (t *IskratelMsan) GetOnuMulticastProfiles() ([]*OnuIgmpProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, onuIgmpProfiles)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*OnuIgmpProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuMulticastProfileTable.MsanOnuMulticastProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetOnuMulticastProfileByName is a helper method that returns a single OnuIgmpProfile struct by name, if exists
func (t *IskratelMsan) GetOnuMulticastProfileByName(name string) (*OnuIgmpProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetOnuMulticastProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteOnuMulticastProfile removes the named OnuIgmpProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteOnuMulticastProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetOnuMulticastProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, onuIgmpProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuMulticastProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostOnuMulticastProfile performs a Post request to t.Host containing serialized data from a OnuIgmpProfile struct, if the name is not already used
func (t *IskratelMsan) PostOnuMulticastProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetOnuMulticastProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The OnuIgmpProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, onuIgmpProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuMulticastProfile:", status)
	// check for ErrNotStatusOk
	return nil

}

// GetL2cpProfiles performs a Get Request to the t.Host and returns a list of the L2cpProfile struct
func (t *IskratelMsan) GetL2cpProfiles() ([]*L2cpProfile, error) {
	rawJson, err := RestGetProfiles(t.Host, l2cpProfiles)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rawJson, &t)
	if err != nil {
		return nil, err
	}
	var list []*L2cpProfile
	for _, v := range t.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanL2CpProfileTable.MsanL2CpProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetL2cpProfileByName is a helper method that returns a single L2cpProfile struct by name, if exists
func (t *IskratelMsan) GetL2cpProfileByName(name string) (*L2cpProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := t.GetL2cpProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteL2cpProfile removes the named L2cpProfile from the t.Host if it exists and is not in use by a ServiceProfile
func (t *IskratelMsan) DeleteL2cpProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := t.GetL2cpProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(t.Host, l2cpProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteL2cpProfile:", status)
	// check for ErrNotStatusOk
	return nil
}

// PostL2cpProfile performs a Post request to t.Host containing serialized data from a L2cpProfile struct, if the name is not already used
func (t *IskratelMsan) PostL2cpProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := t.GetL2cpProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The L2cpProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(t.Host, l2cpProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostL2cpProfile:", status)
	// check for ErrNotStatusOk
	return nil
}
*/
/*
// getStructFields performs a GET request on the host with endpoint as the desired profile, returns populated struct or error
func getStructFields(host, profile string) (t *IskratelMsan, err error) {
	// return struct data from rest endpoint
	jsonData, err := getProfiles(host, profile)
	if err != nil {
		return nil, err
	}

	// create the 'master' struct that nests all possible sub-structs
	structData := new(IskratelMsan)
	// io Copy method maps JSON to Struct sub-template according to `json` struct tags
	json.Unmarshal(jsonData, structData)

	// return populated with relevant endpoint
	return structData, nil
}

func postHandler(host string, ep string) (r string, err error) {
	testName := "RestTest"
	switch ep {
	case serviceProfiles:
		json := generateServiceProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case flowProfiles:
		json := generateFlowProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case vlanProfiles:
		json := generateVlanProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case igmpProfiles:
		json := generateIgmpProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case securityProfiles:
		json := generateSecurityProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case onuFlowProfiles:
		json := generateOnuFlowProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case onuTcontProfiles:
		json := generateOnuTcontProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case onuVlanProfiles: // how is this handled in respect to rule assignment within the profile header?
		json := generateOnuVlanProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case onuVlanRules:
		json := generateOnuVlanRule(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case onuIgmpProfiles:
		json := generateOnuIgmpProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	case l2cpProfiles:
		json := generateL2cpProfile(testData[ep])
		r, err = postProfile(host, ep, endpointEntry[ep], testName, json)
		return r, err
	default:
		fmt.Println("Error posting to ", ep)
	}
	return "", nil
}


// previous method, now all profiles have their own tabwriter method that defines the headers

// tidier way of specifying column headers; requires output awareness of # of fields before compilation
var endpointFields = map[string][]string{
	serviceProfiles:  []string{"Name", "Flow Profile", "Vlan Profile", "ONU Flow Profile", "ONU T-CONT Profile", "ONU VLAN Profile", "vGEM Port", "ONU Term", "Usage"},
	flowProfiles:     []string{"Name", "US Match", "DS Match", "DS Peak Rate", "DS Assured Rate", "DS Queueing Priority", "DS Scheduling Mode", "Usage"},
	vlanProfiles:     []string{"Name", "C-VID", "S-VID", "S-Ethertype", "Native", "Remark", "Usage"},
	igmpProfiles:     []string{"Name", "Snooping", "Proxy", "Usage"},
	securityProfiles: []string{"Name", "Protected", "MAC-SG", "MAC Limit", "Port-Sec", "DAI", "IP-SG", "AppRateLimit", "Storm Control", "Usage"},
	onuFlowProfiles:  []string{"Name", "C-VID Range", "C-PCP", "US CDR", "US PDR", "US Priority", "Usage"},
	onuTcontProfiles: []string{"Name", "Type", "ID", "Fixed", "Assured", "Max", "Usage"},
	onuVlanProfiles:  []string{"Name", "Input TPID", "Output TPID", "Usage"},
	onuVlanRules:     []string{"Name", "Rule ID", "C-VID", "S-VID", "C-TPID", "EtherType", "C-Tag", "S-Tag", "C-PCP", "S-PCP", "C-TPID", "S-TPID"},
	onuIgmpProfiles:  []string{"Name", "Mode", "Proxy", "Fast-Leave", "Usage"},
	l2cpProfiles:     []string{"Name", "Description", "Usage"},
}

// each profile has own methods to modify specific details before updating them to the NE

// example data to fill Create block, modifying this data structure would be the API for generating custom changes, can be parsed from text file
var testData = map[string][]string{
	serviceProfiles:  []string{"RestTest", "RestTest", "", "RestTest", "", "RestTest", "", "", "RestTest", "22", "2", "AAAA"},              //serviceProfiles: []string{"Service", "Flow", "Multicast", "Vlan", "Security", "Onu Flow", "Onu Vlan", "Onu Multicast", "Onu Tcont", "vGem", "Onu-Tp", "UniBit"}
	flowProfiles:     []string{"RestTest", "1", "1", "0", "0", "0", "0"},                                                                   // flowProfiles: []string{"Name", "MatchUsVlanProfile", "MatchDsVlanProfile", "DsPdr", "DsCdr", "UsPdr", "UsCdr"},
	vlanProfiles:     []string{"RestTest", "110"},                                                                                          //vlanProfiles: []string{"RestTest", "CVidList"},
	igmpProfiles:     []string{"RestTest", "1", "0", "0", "0.0.0.0"},                                                                       //igmpProfiles: []string{"Name", "Snooping", "FastLeave", "Proxy", "Proxy-IP"},
	securityProfiles: []string{"RestTest", "1", "0", "0", "0", "0", "1", "0", "-1, -1, 100", "5, 5, 5, 5, 5"},                              //securityProfiles: []string{"Name", "Protected", "MAC-SG", "MAC-Limit", "Port-Sec", "Arp-Inspect", "IP-SG", "IPv6-SG", "Storm-Ctl", "AppRateLimit"},
	onuFlowProfiles:  []string{"RestTest", "110"},                                                                                          //onuFlowProfiles: []string{"Name", "CVID-List"},
	onuTcontProfiles: []string{"RestTest", "4", "5", "512", "2048", "750000"},                                                              //onuTcontProfiles: []string{"Name", "TcontID", "TcontType", "FixedDataRate", "AssuredDataRate", "MaxDataRate"},
	onuVlanProfiles:  []string{"RestTest", "1", "33024", "34984"},                                                                          //onuVlanProfiles: []string{"Name", "DS Mode", "TPID-IN", "TPID-OUT"},
	onuVlanRules:     []string{"RestTest", "15", "4096", "-1", "0", "4096", "-1", "0", "0", "1", "2", "0", "0", "1", "1", "0", "111", "1"}, //onuVlanRules: []string{"Name", "Rule ID", "MatchSVlanID", "MatchSPcp", "MatchSTPID", "MatchCVlanID", "MatchCPcp", "MatchCTPID", "MatchEthertype", "RemoveTags", "AddSTag", "AddSPcp", "AddSVlanID", "AddSTPID", "AddCTag", "AddCPcp", "AddCVlanID", "AddCTPID"},
	onuIgmpProfiles:  []string{"RestTest", "2", "2", "1", "3998"},                                                                          //onuIgmpProfiles: []string{"Name", "IgmpMode", "IgmpProxy", "SnoopingFastLeave", "DsGemPort"},
	l2cpProfiles:     []string{"RestTest", "some text"},                                                                                    //l2cpProfiles: []string{"Name", "Descr."},
}
*/
