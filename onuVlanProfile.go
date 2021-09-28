package goPon

import (
	"fmt"
	"text/tabwriter"
	"encoding/json"
	"os"	
)

type OnuVlanProfile struct {
	Name           string `json:"msanOnuVlanProfileName"`
	DownstreamMode int    `json:"msanOnuVlanProfileDownstreamMode"`
	InputTPID      int    `json:"msanOnuVlanProfileInputTPID"`
	OutputTPID     int    `json:"msanOnuVlanProfileOutputTPID"`
	Usage          int    `json:"msanOnuVlanProfileUsage"`
	Rules          *OnuVlanRuleList
}

type OnuVlanRule struct {
	Name               string `json:"msanOnuVlanProfileName"`
	RuleID             int    `json:"msanOnuVlanProfileRuleId"`
	RuleMatchSVlanID   int    `json:"msanOnuVlanProfileRuleMatchSVlanId"`	// def [4096,4096,-1]
	RuleMatchSPcp      int    `json:"msanOnuVlanProfileRuleMatchSPcp"`		// def [-1,-1,-1]
	RuleMatchSTPID     int    `json:"msanOnuVlanProfileRuleMatchSTPID"`		// def [0, 0, 0]
	RuleMatchCVlanID   int    `json:"msanOnuVlanProfileRuleMatchCVlanId"`	// def [4096, -1, -1]
	RuleMatchCPcp      int    `json:"msanOnuVlanProfileRuleMatchCPcp"`		// def [-1, -1, -1]
	RuleMatchCTPID     int    `json:"msanOnuVlanProfileRuleMatchCTPID"`		// def [0, 0, 0]
	RuleMatchEthertype int    `json:"msanOnuVlanProfileRuleMatchEthertype"`	// def [0, 0, 0]
	RuleRemoveTags     int    `json:"msanOnuVlanProfileRuleRemoveTags"`		// def [1, 1, 1]
	RuleAddSTag        int    `json:"msanOnuVlanProfileRuleAddSTag"`		// def [2, 2, 2]
	RuleAddSPcp        int    `json:"msanOnuVlanProfileRuleAddSPcp"`		// def [0, 0, 0]
	RuleAddSVlanID     int    `json:"msanOnuVlanProfileRuleAddSVlanId"`		// def [0, 0, 0]
	RuleAddSTPID       int    `json:"msanOnuVlanProfileRuleAddSTPID"`		// def [1, 1, 1]
	RuleAddCTag        int    `json:"msanOnuVlanProfileRuleAddCTag"`		// def [2, 2, 2]
	RuleAddCPcp        int    `json:"msanOnuVlanProfileRuleAddCPcp"`		// def [0, 0, 0]
	RuleAddCVlanID     int    `json:"msanOnuVlanProfileRuleAddCVlanId"`		// def [0, 0, 0]
	RuleAddCTPID       int    `json:"msanOnuVlanProfileRuleAddCTPID"`		// def [1, 1, 1]
}

type OnuVlanProfileList struct {
	Entry []*OnuVlanProfile
}

type OnuVlanRuleList struct {
	Entry []*OnuVlanRule
}

func NewOnuVlanProfile(name string) *OnuVlanProfile {
	p := &OnuVlanProfile{
		Name: name,
		DownstreamMode: 1,	// enabled, def
		InputTPID: 33024,	// def val 0x8100
		OutputTPID: 34984,	// def value 0x88a8
	}
	return p
}

func (p *OnuVlanProfile) GetName() string {
	return p.Name
}

func (p *OnuVlanProfile) GetDsMode() string {
	switch p.DownstreamMode {
	case 1:
		return "Enabled"
	case 2:
		return "Disabled"
	default:
		return ""
	}
}

func (p *OnuVlanProfile) IsUsed() bool {
	return p.Usage == 1
}

var OnuVlanProfileHeaders = []string{
	"Name",
	"DsMode",
	"TPID-In",
	"TPID-Out",
	"Rules",
}

func (p *OnuVlanProfile) ListEssentialParams() map[string]interface{} {
	var EssentialOnuVlanProfile = map[string]interface{}{
		OnuVlanProfileHeaders[0]:	p.GetName(),
		OnuVlanProfileHeaders[1]:	p.GetDsMode(),
		OnuVlanProfileHeaders[2]:	p.InputTPID,
		OnuVlanProfileHeaders[3]:	p.OutputTPID,
		OnuVlanProfileHeaders[4]:	p.GetRulesString(),
	}
	return EssentialOnuVlanProfile
}

func (p *OnuVlanProfile) GetRulesString() string {
	var ruleString string
	//fmt.Printf("Length of Rules: %d\n", len(p.Rules.Entry))
	for i := 0; i < len(p.Rules.Entry); i++ {
		ruleString += fmt.Sprintf("%d, ", p.Rules.Entry[i].RuleID)
	}
	return ruleString
}

func (p *OnuVlanProfile) GetRules() (*OnuVlanRuleList, error) {
	var list *OnuVlanRuleList
	for i := 0; i < len(p.Rules.Entry); i++ {
		list.Entry = append(list.Entry, p.Rules.Entry[i])
	}
	if len(list.Entry) < 1 {
		return nil, ErrNotExists
	}
	return list, nil
}

func (p *OnuVlanProfile) GetRuleById(id int) (*OnuVlanRule, error) {
	// if exists, return it, if not, err not exists
	for i := 0; i < len(p.Rules.Entry); i++ {
		if p.Rules.Entry[i].RuleID == id {
			return p.Rules.Entry[i], nil
		}
	}
	return nil, ErrNotExists
}

func (pl *OnuVlanProfileList) GetProfileByName(name string) (*OnuVlanProfile, error) {
	for _, p := range pl.Entry {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, ErrNotExists
}

func (pl *OnuVlanProfileList) GetRulesByProfileName(name string) (*OnuVlanRuleList, error) {
	for _, p := range pl.Entry {
		if p.Name == name {
			return p.GetRules()
		}
	}
	return nil, ErrNotExists
}

/*
func (p *OnuVlanProfile) NewRule(id int) *OnuVlanRule {
	// first check if that rule exists for that profile
	rule, err := p.GetRuleById(id)
	if err == nil {
		return nil, ErrExists
	}
	rule = &OnuVlanRule{
		Name:   profile,
		RuleID: id,
	}
	return rule, nil
}

func (p *OnuVlanProfile) DeleteRule(id int) error {
	rule, err := p.GetRuleById(id)
	if err != nil {
		return ErrNotExists
	}
	return rule.Delete()
}

func (r *OnuVlanRule) Delete() error {
	// not sure how to do this
}
*/

var OnuVlanRuleHeaders = []string{
	"Profile",
	"ID",
	"Match Criteria",
	"Add/Remove Actions",
}

func (r *OnuVlanRule) ListEssentialParams() map[string]interface{} {
	var EssentialOnuVlanRules = map[string]interface{}{
		OnuVlanRuleHeaders[0]: r.Name,
		OnuVlanRuleHeaders[1]: r.RuleID,
		OnuVlanRuleHeaders[2]: r.GetMatchCriteriaString(),
		OnuVlanRuleHeaders[3]: r.GetActionListString(),
	}
	return EssentialOnuVlanRules
}

var DefRule97MatchCriteria = map[string]int{
	OnuVlanRuleMatchCriteria[0]: 4096,
	OnuVlanRuleMatchCriteria[1]: -1,
	OnuVlanRuleMatchCriteria[2]: 0,
	OnuVlanRuleMatchCriteria[3]: 4096,
	OnuVlanRuleMatchCriteria[4]: -1,
	OnuVlanRuleMatchCriteria[5]: 0,
	OnuVlanRuleMatchCriteria[6]: 0,
}

var DefRule98MatchCriteria = map[string]int{
	OnuVlanRuleMatchCriteria[0]: 4096,
	OnuVlanRuleMatchCriteria[1]: -1,
	OnuVlanRuleMatchCriteria[2]: 0,
	OnuVlanRuleMatchCriteria[3]: -1,
	OnuVlanRuleMatchCriteria[4]: -1,
	OnuVlanRuleMatchCriteria[5]: 0,
	OnuVlanRuleMatchCriteria[6]: 0,
}

var DefRule99MatchCriteria = map[string]int{
	OnuVlanRuleMatchCriteria[0]: -1,
	OnuVlanRuleMatchCriteria[1]: -1,
	OnuVlanRuleMatchCriteria[2]: 0,
	OnuVlanRuleMatchCriteria[3]: -1,
	OnuVlanRuleMatchCriteria[4]: -1,
	OnuVlanRuleMatchCriteria[5]: 0,
	OnuVlanRuleMatchCriteria[6]: 0,
}

var OnuVlanRuleMatchCriteria = []string {
	"S-VID",
	"S-PCP",
	"S-TPID",
	"C-VID",
	"C-PCP",
	"C-TPID",
	"Ethertype",
}

// SVlanId, SPcp, STPID, CVlanId, CPcp, CTPID, Ethertype
func (r *OnuVlanRule) GetMatchCriteria() []int {
	list := []int{
		r.RuleMatchSVlanID,
		r.RuleMatchSPcp,
		r.RuleMatchSTPID,
		r.RuleMatchCVlanID,
		r.RuleMatchCPcp,
		r.RuleMatchCTPID,
		r.RuleMatchEthertype,
	}
	return list
}

func (r *OnuVlanRule) GetMatchCriteriaString() string {
	list := r.GetMatchCriteria()
	var matchString string
	for i, v := range list {
		if v != DefRule97MatchCriteria[OnuVlanRuleMatchCriteria[i]] {
			if v != DefRule98MatchCriteria[OnuVlanRuleMatchCriteria[i]] {
				if v != DefRule99MatchCriteria[OnuVlanRuleMatchCriteria[i]] {
					matchString += fmt.Sprintf("%s:%v, ", OnuVlanRuleMatchCriteria[i], v)
				}
			}
		}
	}
	if matchString == "" {
		matchString = "Default"
	}
	return matchString
}

func (r *OnuVlanRule) SetMatchCriteria(list []int) error {
	switch len(list) {
	case 7:
		for i := 0; i < 7; i++ {
			if list[i] < -1 {
				list[i] = -1
			}
			if list[i] > 4096 {
				list[i] = 4096
			}
		}
	default:
		return ErrNotInput
	}
	r.RuleMatchSVlanID = list[0]
	r.RuleMatchSPcp = list[1]
	r.RuleMatchSTPID = list[2]
	r.RuleMatchCVlanID = list[3]
	r.RuleMatchCPcp = list[4]
	r.RuleMatchCTPID = list[5]
	r.RuleMatchEthertype = list[6]
	return nil
}

/*
RuleID
RuleMatchSVlanID [4096,4096,-1]
RuleMatchSPcp [-1,-1,-1]
RuleMatchSTPID [0, 0, 0]
RuleMatchCVlanID [4096, -1, -1]
RuleMatchCPcp [-1, -1, -1]
RuleMatchCTPID [0, 0, 0]
RuleMatchEthertype [0, 0, 0]
RuleRemoveTags [1, 1, 1]
RuleAddSTag [2, 2, 2]
RuleAddSPcp [0, 0, 0]
RuleAddSVlanID [0, 0, 0]
RuleAddSTPID [1, 1, 1]
RuleAddCTag [2, 2, 2]
RuleAddCPcp [0, 0, 0]
RuleAddCVlanID [0, 0, 0]
RuleAddCTPID [1, 1, 1]
*/

var DefaultActionList = map[string]int{
	OnuVlanRuleActionList[0]: 1,
	OnuVlanRuleActionList[1]: 2,
	OnuVlanRuleActionList[2]: 0,
	OnuVlanRuleActionList[3]: 0,
	OnuVlanRuleActionList[4]: 1,
	OnuVlanRuleActionList[5]: 2,
	OnuVlanRuleActionList[6]: 0,
	OnuVlanRuleActionList[7]: 0,
	OnuVlanRuleActionList[8]: 1,
}

var OnuVlanRuleActionList = []string{
	"Rem Tags",
	"Add S-Tag",
	"Add S-PCP",
	"Add S-VID",
	"Add S-TPID",
	"Add C-Tag",
	"Add C-PCP",
	"Add C-VID",
	"Add C-TPID",
}

// RemoveTags, AddSTag, AddSPcp, AddSVlanId, AddSTPID, AddCTag, AddCPcp, AddCVlanId, AddCTPID
func (r *OnuVlanRule) GetActionList() []int {
	list := []int{
		r.RuleRemoveTags,
		r.RuleAddSTag,
		r.RuleAddSPcp,
		r.RuleAddSVlanID,
		r.RuleAddSTPID,
		r.RuleAddCTag,
		r.RuleAddCPcp,
		r.RuleAddCVlanID,
		r.RuleAddCTPID,
	}
	return list
}

func (r *OnuVlanRule) GetActionListString() string {
	list := r.GetActionList()
	var actionList string
	for i, v := range list {
		if v != DefaultActionList[OnuVlanRuleActionList[i]] {
			actionList += fmt.Sprintf("%s:%v, ", OnuVlanRuleActionList[i], v)
		}
	}
	if actionList == "" {
		actionList = "Default"
	}
	return actionList
}

func (r *OnuVlanRule) SetActions(list []int) error {
	switch len(list) {
	case 9:
		for i := 0; i < 7; i++ {
			if list[i] < 0 {
				list[i] = 0
			}
			if list[i] > 4096 {
				list[i] = 4096
			}
		}
	default:
		return ErrNotInput
	}
	r.RuleRemoveTags = list[0]
	r.RuleAddSTag = list[1]
	r.RuleAddSPcp = list[2]
	r.RuleAddSVlanID = list[3]
	r.RuleAddSTPID = list[4]
	r.RuleAddCTag = list[5]
	r.RuleAddCPcp = list[6]
	r.RuleAddCVlanID = list[7]
	r.RuleAddCTPID = list[8]
	return nil
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (p *OnuVlanProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", data
	}
	return p.Name, data
}

func (ovpl *OnuVlanProfileList) Tabwrite() {
	fmt.Println("|| ONU VLAN Profile List ||")
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range OnuVlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuVlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	for _, vp := range ovpl.Entry {
		// first get the data as a map
		l := vp.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range OnuVlanProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuVlanProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}

func (ovrl *OnuVlanRuleList) Tabwrite() {
	fmt.Println("|| ONU VLAN Rule List ||")
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range OnuVlanRuleHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuVlanRuleHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	for _, vr := range ovrl.Entry {
		// first get the data as a map
		l := vr.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range OnuVlanRuleHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuVlanRuleHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}
