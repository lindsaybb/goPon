package goPon

import (
	"encoding/json"
)

type OnuConfig struct {
	IfName                 string `json:"msanOnuCfgIfName"`
	Password               string `json:"msanOnuCfgPassword"`
	EnablePm               int    `json:"msanOnuCfgEnablePm"`
	SerialNumber           string `json:"msanOnuCfgSerialNumber"`
	AdminState             int    `json:"msanOnuCfgAdminState"`
	OnuDhcpMode            int    `json:"msanOnuCfgOnuDhcpMode"`
	OnuIPAddress           string `json:"msanOnuCfgOnuIpAddress"`
	OnuIPMask              string `json:"msanOnuCfgOnuIPMask"`
	OnuDefaultGateway      string `json:"msanOnuCfgOnuDefaultGateway"`
	OnuReset               int    `json:"msanOnuCfgOnuReset"`
	OnuResetBackupImage    int    `json:"msanOnuCfgOnuResetBackupImage"`
	DefaultConfigFile      string `json:"msanOnuCfgDefaultConfigFile"`
	SendConfig             int    `json:"msanOnuCfgSendConfig"`
	SendConfigStatus       int    `json:"msanOnuCfgSendConfigStatus"`
	OnuResync              int    `json:"msanOnuCfgOnuResync"`
	OnuResetFactoryDefault int    `json:"msanOnuCfgOnuResetFactoryDefault"`
}

type OnuConfigList struct {
	Entry []*OnuConfig
}

// NewOnuConfig accepts an Onu Serial Number and an Onu Interface (Not Url Encoded 0/x/y)
// and provides an OnuConfig object
func NewOnuConfig(sn, intf string) *OnuConfig {
	o := &OnuConfig{
		IfName:                 intf,
		Password:               "",
		EnablePm:               0, // disabled by default
		SerialNumber:           sn,
		AdminState:             1,
		OnuDhcpMode:            1,
		OnuIPAddress:           "0.0.0.0",
		OnuIPMask:              "0.0.0.0",
		OnuDefaultGateway:      "0.0.0.0",
		DefaultConfigFile:      "",
		SendConfig:             2,
		SendConfigStatus:       6,
		OnuResync:              2,
		OnuResetFactoryDefault: 2,
	}
	//fmt.Println(o)
	return o
}

// GenerateBlankConfig allows a configured Onu to be removed by patching the interface with a blank serial number
func GenerateBlankConfig(intf string) *OnuConfig {
	o := &OnuConfig{
		IfName:       intf,
		Password:     "",
		EnablePm:     2,
		SerialNumber: "",
		AdminState:   1,
	}
	//fmt.Println(o)
	return o
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (o *OnuConfig) GenerateJson() (intf string, data []byte) {
	data, err := json.Marshal(o)
	if err != nil {
		return "", data
	}
	return o.IfName, data
}
