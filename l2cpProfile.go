package goPon

import "encoding/json"

type L2cpProfile struct {
	Name        string `json:"msanL2cpProfileName"`
	Description string `json:"msanL2cpProfileDescription"`
	Usage       int    `json:"msanL2cpProfileUsage"`
}

func NewL2cpProfile(name string) *L2cpProfile {
	p := &L2cpProfile{
		Name: name,
	}
	return p
}

func (p *L2cpProfile) GetName() string {
	return p.Name
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (p *L2cpProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(p)
	if err != nil {
		return "", data
	}
	return p.Name, data
}
