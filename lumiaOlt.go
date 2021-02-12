package gopon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type LumiaOlt struct {
	Host            string
	Current         *IskratelMsan
	Cache           *IskratelMsan
	AuthorizedOnuSn []string
	RegisteredOnuSn map[string]string
}

// OnuRegistrationMap pairs the Onu Logical Interface(0/x/y) with the attached SerialNumber
var OnuRegistrationMap map[string]string

// NewIskratelMsan sets up a data structure for the specified Host
func NewLumiaOlt(host string) *LumiaOlt {
	var t = &LumiaOlt{
		Host:    host,
		Current: NewIskratelMsan(),
		Cache:   NewIskratelMsan(),
	}
	return t
}

// HostIsReachable is a helper method to ensure HTTPS access on port 443 to specified Host address
func (l *LumiaOlt) HostIsReachable() bool {
	err := CheckHost(l.Host, 1)
	return err == nil
}

func (l *LumiaOlt) CacheSwap() {
	// logprinting as a placeholder for some form of logging of past caches/debug visibility
	//log.Println(l.Cache)
	l.Cache = l.Current
}

func (l *LumiaOlt) CacheBack() {
	// reversed from cache if error in call
	l.Current = l.Cache
}

/*
func (l *LumiaOlt) CacheDiff() []interface{} {
	// compare l.Current and l.Cache
	return // any items in one but not the other
}

// GetIskratelMsan returns a populated struct of all active profiles on the OLT
func (l *LumiaOlt) GetIskratelMsan() *IskratelMsan {
	p := GetAllProfiles()
	l.CacheSwap()
	l.Current = p
	return p
}

// GetAllProfiles is a helper/merger function that combines all individual IskratelMsan GET logic
func (l *LumiaOlt) GetAllProfiles() (*IskratelMsan, error) {
	sp := NewIskratelMsan()
}
*/

// GetCurrentLogs accepts a path as output directory location
func (l *LumiaOlt) GetCurrentLogs(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	info, _ := os.Stat(absPath)
	if !info.IsDir() {
		absPath = filepath.Dir(absPath)
	}
	cl, err := NewFtpClient(l.Host, auth)
	if err != nil {
		return err
	}
	//numFiles, err := GetOltLogs(cl, absPath, false)
	_, err = GetOltLogs(cl, absPath, false)
	cl.Close()
	return err
}

// GetCurrentLogs accepts a path as output directory location, and prints Tail (last 10 lines) of each log as it is downloading
func (l *LumiaOlt) GetCurrentLogsTail(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	info, _ := os.Stat(absPath)
	if !info.IsDir() {
		absPath = filepath.Dir(absPath)
	}
	cl, err := NewFtpClient(l.Host, auth)
	if err != nil {
		return err
	}
	//numFiles, err := GetOltLogs(cl, absPath, false)
	_, err = GetOltLogs(cl, absPath, true)
	cl.Close()
	return err
}

// UploadConfig uses Ftp to transfer Script (.scr) or InnboxConfig (.conf) to Olt from supplied path
func (l *LumiaOlt) UploadConfig(path string) error {
	var isScr, isConf bool
	switch {
	case strings.HasSuffix(path, ".scr"):
		isScr = true
	case strings.HasSuffix(path, ".conf"):
		isConf = true
	default:
		return ErrNotInput
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	file, err := os.Open(absPath)
	if err != nil {
		return err
	}
	cl, err := NewFtpClient(l.Host, auth)
	if err != nil {
		return err
	}
	if isScr {
		err = PutOltConfig(cl, file)
	} else if isConf {
		err = PutInnboxConfig(cl, file)
	}
	cl.Close()
	return err
}

// DeleteConfig uses Ftp to remove Script (.scr) or InnboxConfig (.conf) from Olt ny name
func (l *LumiaOlt) DeleteConfig(path string) error {
	var isScr, isConf bool
	switch {
	case strings.HasSuffix(path, ".scr"):
		isScr = true
	case strings.HasSuffix(path, ".conf"):
		isConf = true
	default:
		return ErrNotInput // while it is possible to delete other files, we will limit scope
	}
	// ensure the file does not contain other path info harmful to the operation
	path = filepath.Base(path)

	cl, err := NewFtpClient(l.Host, auth)
	if err != nil {
		return err
	}
	if isScr {
		err = DeleteOltConfig(cl, path)
	} else if isConf {
		err = DeleteInnboxConfig(cl, path)
	}
	cl.Close()
	return err
}

func (l *LumiaOlt) GetOnuBlacklist() (*OnuBlacklistList, error) {
	rawJson, err := RestGetProfiles(l.Host, onuBlacklist)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(rawJson))
	l.CacheSwap()
	err = json.Unmarshal(rawJson, &l.Current)
	if err != nil {
		l.CacheBack()
		return nil, err
	}
	var list OnuBlacklistList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuBlackListTable.MsanOnuBlackListEntry); i++ {
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuBlackListTable.MsanOnuBlackListEntry[i])
	}
	return &list, nil
}

// LoadOnuAuthList opens the supplied filepath and reads line-separated entries to build a slice of registered serial numbers
// this should be modified to work with csv/excel files for simpler handling (maybe an external handler)
func (l *LumiaOlt) LoadOnuAuthList(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	authFile, err := os.Open(absPath)
	if err != nil {
		return err
	}
	var reg []string
	s := bufio.NewScanner(authFile)
	for s.Scan() {
		sn := strings.TrimSpace(s.Text())
		if len(sn) == 12 {
			reg = append(reg, sn)
		}
		if len(sn) == 8 {
			sn = "ISKT" + sn
			reg = append(reg, sn)
		}
	}
	// should this be an append or an overwrite?
	l.AuthorizedOnuSn = reg
	fmt.Printf("Authorized Onu SerialNumber list now has %d entries\n", len(l.AuthorizedOnuSn))
	return nil
}

// UpdateRegisteredOnuList updates the Olt's record of the Onu Serial Numbers currently active in the system
// this list may differ from the AuthorizedOnuList if devices are pre-authorized but not yet deployed
func (l *LumiaOlt) UpdateRegisteredOnuList() error {
	rawJson, err := RestGetProfiles(l.Host, onuConfig)
	if err != nil {
		return err
	}
	//fmt.Println(string(rawJson))
	l.CacheSwap()
	err = json.Unmarshal(rawJson, &l.Current)
	if err != nil {
		l.CacheBack()
		return err
	}
	reg := make(map[string]string)
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuCfgTable.MsanOnuCfgEntry); i++ {
		reg[l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuCfgTable.MsanOnuCfgEntry[i].IfName] = l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuCfgTable.MsanOnuCfgEntry[i].SerialNumber
		//list = append(list, l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuCfgTable.MsanOnuCfgEntry[i].SerialNumber)
	}
	//if len(list) < 1 {
	//	return ErrNotExists
	//}
	l.RegisteredOnuSn = reg

	return nil
}

// ValidateSn loops over the list of Authorized Onu Serial Numbers in the Olt and returns a bool if the supplied value exists in list
func (l *LumiaOlt) ValidateSn(sn string) bool {
	for _, v := range l.AuthorizedOnuSn {
		if sn == v {
			return true
		}
	}
	return false
}

// NextAvailableInterface receives an Olt Interface (0/x) and checks the OnuRegistrationMap to find the next available Onu Subinterface(0/x/y)
// returns nil if there are no available interfaces left (1-128); otherwise returns the given Olt interface with an Onu Subinterface attached
func (l *LumiaOlt) NextAvailableOnuInterface(intf string) string {
	var list []int
	var entry int
	for k := range l.RegisteredOnuSn {
		if strings.Contains(k, intf) {
			add := strings.Split(k, "/")
			// this is a controlled list, where we control the values and can assume the length will be 3
			// the third segment represents the in-use sub-interface on the filtered olt port
			entry, _ = strconv.Atoi(add[2])
			list = append(list, entry)
		}
	}
	if len(list) < 1 {
		fmt.Println("No used ONU on this interface")
		return fmt.Sprintf("%s/%d", intf, 1)
	}
	sort.Ints(list)
	var counter int
	for _, v := range list {
		counter++
		if counter != v {
			break
		}
	}
	if counter > 128 {
		return ""
	}
	return fmt.Sprintf("%s/%d", intf, counter)
}

// AuthorizeOnuCheckBlacklist accepts a single OnuConfig object and attempts to register the device after checking the blacklist that it exists
// unintended consequence is "onu params" being inserted to running config. Find out why and use this as a feature/function
func (l *LumiaOlt) AuthorizeOnuCheckBlacklist(ocfg *OnuConfig) error {
	// checking Blacklist is an unnecessary step here, good precaution but useless, server will reject if doesn't
	obll, err := l.GetOnuBlacklist()
	if err != nil {
		return err
	}
	if len(obll.Entry) == 0 {
		return ErrNotExists
	}
	for _, e := range obll.Entry {
		if e.SerialNumber == ocfg.SerialNumber {
			if l.ValidateSn(ocfg.SerialNumber) {
				ifName, jsonData := ocfg.GenerateJson()
				if ifName == "" {
					return ErrNotStruct
				}
				//fmt.Println(jsonData)
				resp, err := RestPatchProfile(l.Host, onuConfig, UrlEncodeInterface(ifName), jsonData)
				if err != nil {
					return err
				}
				if resp != responseOk {
					fmt.Println(resp)
					return ErrNotStatusOk
				}
				return nil
			} else {
				return ErrNotAuthorized
			}
		}
	}
	return ErrNotExists
}

// AuthorizeOnu accepts a single OnuConfig object with a ServiceProfile name and attempts to register the device
func (l *LumiaOlt) AuthorizeOnu(ocfg *OnuConfig) error {
	// unintended consequence is "onu params" being inserted to running config. Find out why and use this as a feature/function
	if l.ValidateSn(ocfg.SerialNumber) {
		ifName, jsonData := ocfg.GenerateJson()
		if ifName == "" {
			return ErrNotStruct
		}
		//fmt.Println(jsonData)
		resp, err := RestPatchProfile(l.Host, onuConfig, UrlEncodeInterface(ifName), jsonData)
		if err != nil {
			return err
		}
		if resp != responseOk {
			fmt.Println(resp)
			return ErrNotStatusOk
		}
		return nil
	} else {
		return ErrNotAuthorized
	}
}

func (l *LumiaOlt) GetOnuProfileUsage() (*OnuProfileList, error) {
	rawJson, err := RestGetProfiles(l.Host, onuProfiles)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(rawJson))
	l.CacheSwap()
	err = json.Unmarshal(rawJson, &l.Current)
	if err != nil {
		l.CacheBack()
		return nil, err
	}
	var list OnuProfileList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanServicePortProfileTable.MsanServicePortProfileEntry); i++ {
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanServicePortProfileTable.MsanServicePortProfileEntry[i])
	}
	return &list, nil
}

func (l *LumiaOlt) PostOnuProfile(op *OnuProfile) error {
	ifName, jsonData := op.GenerateJson()
	if ifName == "" {
		return ErrNotStruct
	}
	//fmt.Println(jsonData)
	resp, err := RestPostProfile(l.Host, onuProfiles, UrlEncodeInterface(ifName), jsonData)
	if err != nil {
		return err
	}
	if resp != responseOk {
		fmt.Println(resp)
		return ErrNotStatusOk
	}
	return nil
}

// RemoveOnuProfileUsage is actually a
func (l *LumiaOlt) RemoveOnuProfileUsage(intf, spName string) error {
	removalQuery := UrlEncodeInterface(intf) + "," + spName
	status, err := RestDeleteProfile(l.Host, onuProfiles, removalQuery)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// GetServiceProfiles performs a Get Request to the l.Host and returns a list of the ServiceProfile struct
func (l *LumiaOlt) GetServiceProfiles() (*ServiceProfileList, error) {
	rawJson, err := RestGetProfiles(l.Host, serviceProfiles)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(rawJson))
	l.CacheSwap()
	err = json.Unmarshal(rawJson, &l.Current)
	if err != nil {
		l.CacheBack()
		return nil, err
	}
	var list ServiceProfileList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanServiceProfileTable.MsanServiceProfileEntry); i++ {
		//fmt.Println(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanServiceProfileTable.MsanServiceProfileEntry[i])
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanServiceProfileTable.MsanServiceProfileEntry[i])
	}
	return &list, nil
}

// GetServiceProfileByName is a helper method that returns a single ServiceProfile struct by name, if exists
func (l *LumiaOlt) GetServiceProfileByName(name string) (*ServiceProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetServiceProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list.Entry {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteServiceProfile removes the named ServiceProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteServiceProfile(name string) error {

	status, err := RestDeleteProfile(l.Host, serviceProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteServiceProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostServiceProfile performs a Post request to l.Host containing serialized data from a ServiceProfile struct, if the name is not already used
func (l *LumiaOlt) PostServiceProfile(name string, data []byte) error {

	status, err := RestPostProfile(l.Host, serviceProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostServiceProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil
}

// GetFlowProfiles performs a Get Request to the l.Host and returns a list of the FlowProfile struct
func (l *LumiaOlt) GetFlowProfiles() (*FlowProfileList, error) {
	rawJson, err := RestGetProfiles(l.Host, flowProfiles)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list FlowProfileList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanServiceFlowProfileTable.MsanServiceFlowProfileEntry); i++ {
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanServiceFlowProfileTable.MsanServiceFlowProfileEntry[i])
	}
	return &list, nil
}

// GetFlowProfileByName is a helper method that returns a single FlowProfile struct by name, if exists
func (l *LumiaOlt) GetFlowProfileByName(name string) (*FlowProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetFlowProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list.Entry {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteFlowProfile removes the named FlowProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteFlowProfile(name string) error {

	status, err := RestDeleteProfile(l.Host, flowProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteFlowProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostFlowProfile performs a Post request to l.Host containing serialized data from a FlowProfile struct, if the name is not already used
func (l *LumiaOlt) PostFlowProfile(name string, data []byte) error {

	status, err := RestPostProfile(l.Host, flowProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostFlowProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil

}

// GetVlanProfiles performs a Get Request to the l.Host and returns a list of the VlanProfile struct
func (l *LumiaOlt) GetVlanProfiles() (*VlanProfileList, error) {
	rawJson, err := RestGetProfiles(l.Host, vlanProfiles)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list VlanProfileList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanVlanProfileTable.MsanVlanProfileEntry); i++ {
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanVlanProfileTable.MsanVlanProfileEntry[i])
	}
	return &list, nil

}

// GetVlanProfileByName is a helper method that returns a single VlanProfile struct by name, if exists
func (l *LumiaOlt) GetVlanProfileByName(name string) (*VlanProfile, error) {
	list, err := l.GetVlanProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list.Entry {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteVlanProfile removes the named VlanProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteVlanProfile(name string) error {

	status, err := RestDeleteProfile(l.Host, vlanProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteVlanProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostVlanProfile performs a Post request to l.Host containing serialized data from a VlanProfile struct, if the name is not already used
func (l *LumiaOlt) PostVlanProfile(name string, data []byte) error {

	status, err := RestPostProfile(l.Host, vlanProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostVlanProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil
}

// GetOnuFlowProfiles performs a Get Request to the l.Host and returns a list of the OnuFlowProfile struct
func (l *LumiaOlt) GetOnuFlowProfiles() (*OnuFlowProfileList, error) {
	rawJson, err := RestGetProfiles(l.Host, onuFlowProfiles)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list OnuFlowProfileList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuFlowProfileTable.MsanOnuFlowProfileEntry); i++ {
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuFlowProfileTable.MsanOnuFlowProfileEntry[i])
	}
	return &list, nil

}

// GetOnuFlowProfileByName is a helper method that returns a single OnuFlowProfile struct by name, if exists
func (l *LumiaOlt) GetOnuFlowProfileByName(name string) (*OnuFlowProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetOnuFlowProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list.Entry {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteOnuFlowProfile removes the named OnuFlowProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteOnuFlowProfile(name string) error {

	status, err := RestDeleteProfile(l.Host, onuFlowProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuFlowProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostOnuFlowProfile performs a Post request to l.Host containing serialized data from a OnuFlowProfile struct, if the name is not already used
func (l *LumiaOlt) PostOnuFlowProfile(name string, data []byte) error {

	status, err := RestPostProfile(l.Host, onuFlowProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuFlowProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil
}

// GetOnuTcontProfiles performs a Get Request to the l.Host and returns a list of the OnuTcontProfile struct
func (l *LumiaOlt) GetOnuTcontProfiles() (*OnuTcontProfileList, error) {
	rawJson, err := RestGetProfiles(l.Host, onuTcontProfiles)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list OnuTcontProfileList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuTcontProfileTable.MsanOnuTcontProfileEntry); i++ {
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuTcontProfileTable.MsanOnuTcontProfileEntry[i])
	}
	return &list, nil

}

// GetOnuTcontProfileByName is a helper method that returns a single OnuTcontProfile struct by name, if exists
func (l *LumiaOlt) GetOnuTcontProfileByName(name string) (*OnuTcontProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetOnuTcontProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list.Entry {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteOnuTcontProfile removes the named OnuTcontProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteOnuTcontProfile(name string) error {

	status, err := RestDeleteProfile(l.Host, onuTcontProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuTcontProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostOnuTcontProfile performs a Post request to l.Host containing serialized data from a OnuTcontProfile struct, if the name is not already used
func (l *LumiaOlt) PostOnuTcontProfile(name string, data []byte) error {

	status, err := RestPostProfile(l.Host, onuTcontProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuTcontProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil
}

// GetSecurityProfiles performs a Get Request to the l.Host and returns a list of the SecurityProfile struct
func (l *LumiaOlt) GetSecurityProfiles() (*SecurityProfileList, error) {
	rawJson, err := RestGetProfiles(l.Host, securityProfiles)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list SecurityProfileList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanSecurityProfileTable.MsanSecurityProfileEntry); i++ {
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanSecurityProfileTable.MsanSecurityProfileEntry[i])
	}
	return &list, nil
}

// GetSecurityProfileByName is a helper method that returns a single SecurityProfile struct by name, if exists
func (l *LumiaOlt) GetSecurityProfileByName(name string) (*SecurityProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetSecurityProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list.Entry {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteSecurityProfile removes the named SecurityProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteSecurityProfile(name string) error {

	status, err := RestDeleteProfile(l.Host, securityProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteSecurityProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostSecurityProfile performs a Post request to l.Host containing serialized data from a SecurityProfile struct, if the name is not already used
func (l *LumiaOlt) PostSecurityProfile(name string, data []byte) error {

	status, err := RestPostProfile(l.Host, securityProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostSecurityProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil
}

// GetMulticastProfiles performs a Get Request to the l.Host and returns a list of the IgmpProfile struct
func (l *LumiaOlt) GetMulticastProfiles() (*IgmpProfileList, error) {
	rawJson, err := RestGetProfiles(l.Host, igmpProfiles)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list IgmpProfileList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanMulticastProfileTable.MsanMulticastProfileEntry); i++ {
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanMulticastProfileTable.MsanMulticastProfileEntry[i])
	}
	return &list, nil

}

// GetMulticastProfileByName is a helper method that returns a single IgmpProfile struct by name, if exists
func (l *LumiaOlt) GetMulticastProfileByName(name string) (*IgmpProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetMulticastProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list.Entry {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteMulticastProfile removes the named IgmpProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteMulticastProfile(name string) error {

	status, err := RestDeleteProfile(l.Host, igmpProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteMulticastProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostMulticastProfile performs a Post request to l.Host containing serialized data from a IgmpProfile struct, if the name is not already used
func (l *LumiaOlt) PostMulticastProfile(name string, data []byte) error {

	status, err := RestPostProfile(l.Host, igmpProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostMulticastProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil
}

// GetOnuMulticastProfiles performs a Get Request to the l.Host and returns a list of the OnuIgmpProfile struct
func (l *LumiaOlt) GetOnuMulticastProfiles() (*OnuIgmpProfileList, error) {
	rawJson, err := RestGetProfiles(l.Host, onuIgmpProfiles)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list OnuIgmpProfileList
	for i := 0; i < len(l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuMulticastProfileTable.MsanOnuMulticastProfileEntry); i++ {
		list.Entry = append(list.Entry, &l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuMulticastProfileTable.MsanOnuMulticastProfileEntry[i])
	}
	return &list, nil
}

// GetOnuMulticastProfileByName is a helper method that returns a single OnuIgmpProfile struct by name, if exists
func (l *LumiaOlt) GetOnuMulticastProfileByName(name string) (*OnuIgmpProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetOnuMulticastProfiles()
	if err != nil {
		return nil, err
	}
	for _, v := range list.Entry {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrNotExists
}

// DeleteOnuMulticastProfile removes the named OnuIgmpProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteOnuMulticastProfile(name string) error {

	status, err := RestDeleteProfile(l.Host, onuIgmpProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuMulticastProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostOnuMulticastProfile performs a Post request to l.Host containing serialized data from a OnuIgmpProfile struct, if the name is not already used
func (l *LumiaOlt) PostOnuMulticastProfile(name string, data []byte) error {

	status, err := RestPostProfile(l.Host, onuIgmpProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuMulticastProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil

}

// GetOnuVlanProfiles performs a Get Request to the l.Host and returns a list of the OnuVlanProfile struct
func (l *LumiaOlt) GetOnuVlanProfiles() ([]*OnuVlanProfile, error) {
	rawJson, err := RestGetProfiles(l.Host, onuVlanProfiles)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list []*OnuVlanProfile
	for _, v := range l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuVlanProfileTable.MsanOnuVlanProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetOnuVlanProfileByName is a helper method that returns a single OnuVlanProfile struct by name, if exists
func (l *LumiaOlt) GetOnuVlanProfileByName(name string) (*OnuVlanProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetOnuVlanProfiles()
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

// DeleteOnuVlanProfile removes the named OnuVlanProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteOnuVlanProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := l.GetOnuVlanProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(l.Host, onuVlanProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuVlanProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostOnuVlanProfile performs a Post request to l.Host containing serialized data from a OnuVlanProfile struct, if the name is not already used
func (l *LumiaOlt) PostOnuVlanProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := l.GetOnuVlanProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The OnuVlanProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(l.Host, onuVlanProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuVlanProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil
}

/*
// GetOnuVlanRules performs a Get Request to the l.Host and returns a list of the OnuVlanRule struct
func (l *LumiaOlt) GetOnuVlanRules() ([]*OnuVlanRule, error) {
	rawJson, err := RestGetProfiles(l.Host, onuVlanRules)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list []*OnuVlanRule
	for _, v := range l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanOnuVlanProfileRuleTable.MsanOnuVlanProfileRuleEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetOnuVlanRuleByName is a helper method that returns a single OnuVlanRule struct by name, if exists
func (l *LumiaOlt) GetOnuVlanRuleByName(name string) (*OnuVlanRule, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetOnuVlanRules()
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

// DeleteOnuVlanRule removes the named OnuVlanRule from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteOnuVlanRule(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := l.GetOnuVlanRuleByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(l.Host, onuVlanRules, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteOnuVlanRule:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// CANNOT POST RULES INDEPENDENT OF THE ONUVLANPROFILE, THESE FUNCTIONS ARE A SUBSET OF THAT PROFILE'S OPERATIONS
// PostOnuVlanRule performs a Post request to l.Host containing serialized data from a OnuVlanRule struct, if the name is not already used
func (l *LumiaOlt) PostOnuVlanRule(name string, data []byte) error {
	// check if name is already in use
	_, err := l.GetOnuVlanRuleByName(name)
	if err == nil {
		return ErrExists
	}
	// The OnuVlanRule has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(l.Host, onuVlanRules, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostOnuVlanRule:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil
}
*/

// GetL2cpProfiles performs a Get Request to the l.Host and returns a list of the L2cpProfile struct
func (l *LumiaOlt) GetL2cpProfiles() ([]*L2cpProfile, error) {
	rawJson, err := RestGetProfiles(l.Host, l2cpProfiles)
	if err != nil {
		return nil, err
	}
	l.CacheSwap()
	err = json.Unmarshal(rawJson, l.Current)
	if err != nil {
		return nil, err
	}
	var list []*L2cpProfile
	for _, v := range l.Current.ISKRATELMSANMIB.ISKRATELMSANMIB.MsanL2CpProfileTable.MsanL2CpProfileEntry {
		list = append(list, &v)
	}
	return list, nil
}

// GetL2cpProfileByName is a helper method that returns a single L2cpProfile struct by name, if exists
func (l *LumiaOlt) GetL2cpProfileByName(name string) (*L2cpProfile, error) {
	if name == "" {
		return nil, ErrNotInput
	}
	list, err := l.GetL2cpProfiles()
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

// DeleteL2cpProfile removes the named L2cpProfile from the l.Host if it exists and is not in use by a ServiceProfile
func (l *LumiaOlt) DeleteL2cpProfile(name string) error {
	// get individual profile by supplied name value, if exists
	p, err := l.GetL2cpProfileByName(name)
	if err != nil {
		return err
	}
	// cannot delete in-use profile
	if p.Usage == 1 {
		return ErrInUse
	}
	// perform the delete operation
	status, err := RestDeleteProfile(l.Host, l2cpProfiles, name)
	if err != nil {
		return err
	}
	fmt.Println("DeleteL2cpProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotDelete }
	return nil
}

// PostL2cpProfile performs a Post request to l.Host containing serialized data from a L2cpProfile struct, if the name is not already used
func (l *LumiaOlt) PostL2cpProfile(name string, data []byte) error {
	// check if name is already in use
	_, err := l.GetL2cpProfileByName(name)
	if err == nil {
		return ErrExists
	}
	// The L2cpProfile has a method called GenerateJson() that serializes the data as input
	// perform the post operation
	status, err := RestPostProfile(l.Host, l2cpProfiles, name, data)
	if err != nil {
		return err
	}
	fmt.Println("PostL2cpProfile:", status)
	// check for ErrNotStatusOk
	//chg := l.CacheDiff()
	//if chg == nil { return ErrNotPost }
	return nil
}
