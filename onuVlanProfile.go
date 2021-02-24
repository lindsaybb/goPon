package gopon

import "encoding/json"

type OnuVlanProfile struct {
	Name           string `json:"msanOnuVlanProfileName"`
	DownstreamMode int    `json:"msanOnuVlanProfileDownstreamMode"`
	InputTPID      int    `json:"msanOnuVlanProfileInputTPID"`
	OutputTPID     int    `json:"msanOnuVlanProfileOutputTPID"`
	Usage          int    `json:"msanOnuVlanProfileUsage"`
	Rules          []*OnuVlanRule
}

type OnuVlanRule struct {
	Name               string `json:"msanOnuVlanProfileName"`
	RuleID             int    `json:"msanOnuVlanProfileRuleId"`
	RuleMatchSVlanID   int    `json:"msanOnuVlanProfileRuleMatchSVlanId"`
	RuleMatchSPcp      int    `json:"msanOnuVlanProfileRuleMatchSPcp"`
	RuleMatchSTPID     int    `json:"msanOnuVlanProfileRuleMatchSTPID"`
	RuleMatchCVlanID   int    `json:"msanOnuVlanProfileRuleMatchCVlanId"`
	RuleMatchCPcp      int    `json:"msanOnuVlanProfileRuleMatchCPcp"`
	RuleMatchCTPID     int    `json:"msanOnuVlanProfileRuleMatchCTPID"`
	RuleMatchEthertype int    `json:"msanOnuVlanProfileRuleMatchEthertype"`
	RuleRemoveTags     int    `json:"msanOnuVlanProfileRuleRemoveTags"`
	RuleAddSTag        int    `json:"msanOnuVlanProfileRuleAddSTag"`
	RuleAddSPcp        int    `json:"msanOnuVlanProfileRuleAddSPcp"`
	RuleAddSVlanID     int    `json:"msanOnuVlanProfileRuleAddSVlanId"`
	RuleAddSTPID       int    `json:"msanOnuVlanProfileRuleAddSTPID"`
	RuleAddCTag        int    `json:"msanOnuVlanProfileRuleAddCTag"`
	RuleAddCPcp        int    `json:"msanOnuVlanProfileRuleAddCPcp"`
	RuleAddCVlanID     int    `json:"msanOnuVlanProfileRuleAddCVlanId"`
	RuleAddCTPID       int    `json:"msanOnuVlanProfileRuleAddCTPID"`
}

func NewOnuVlanProfile(name string) *OnuVlanProfile {
	p := &OnuVlanProfile{
		Name: name,
	}
	return p
}

func (p *OnuVlanProfile) GetName() string {
	return p.Name
}

/* NEEDS WORK
func (p *OnuVlanProfile) GetRules() ([]*OnuVlanRule, error) {
	return p.GetRuleByName(p.Name)
}

func (p *OnuVlanProfile) GetRuleById(id int) (*OnuVlanRule, error) {
	// if exists, return it, if not, err not exists
	if rule, ok := p.Rules.RuleID[id]; !ok {
		return nil, ErrNotExists
	}
	return rule, nil
}

func (p *OnuVlanProfile) GetRuleByName(name string) ([]*OnuVlanRule, error) {
	var list []*OnuVlanRule
	// this must be http-interactive as they are updated with creation
}

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

// RemoveTags, AddSTag, AddSPcp, AddSVlanId, AddSTPID, AddCTag, AddCPcp, AddCVlanId, AddCTPID
func (r *OnuVlanRule) GetActions() []int {
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
