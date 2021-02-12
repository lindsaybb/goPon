package gopon

import (
	"encoding/json"
)

type OnuConfig struct {
	IfName            string `json:"msanOnuCfgIfName"`
	Password          string `json:"msanOnuCfgPassword"`
	EnablePm          int    `json:"msanOnuCfgEnablePm"`
	SerialNumber      string `json:"msanOnuCfgSerialNumber"`
	AdminState        int    `json:"msanOnuCfgAdminState"`
	OnuDhcpMode       int    `json:"msanOnuCfgOnuDhcpMode"`
	OnuIPAddress      string `json:"msanOnuCfgOnuIpAddress"`
	OnuIPMask         string `json:"msanOnuCfgOnuIPMask"`
	OnuDefaultGateway string `json:"msanOnuCfgOnuDefaultGateway"`
	//OnuReset				 int	`json:"msanOnuCfgOnuReset"`
	//OnuResetBackupImage	 int	`json:"msanOnuCfgOnuResetBackupImage"`
	//DefaultConfigFile      string `json:"msanOnuCfgDefaultConfigFile"`
	//SendConfig             int    `json:"msanOnuCfgSendConfig"`
	//SendConfigStatus       int    `json:"msanOnuCfgSendConfigStatus"`
	//OnuResync              int    `json:"msanOnuCfgOnuResync"`
	//OnuResetFactoryDefault int    `json:"msanOnuCfgOnuResetFactoryDefault"`
}

type OnuConfigList struct {
	Entry []*OnuConfig
}

// NewOnuConfig accepts an Onu Serial Number and an Onu Interface(Url Encoded 0/x/y) and provides an OnuConfig object
func NewOnuConfig(sn, intf string) *OnuConfig {
	//if !isUrlEncoded(intf) {
	//	fmt.Println("Supplied Interface is not Url Encoded")
	//	return nil
	//}
	o := &OnuConfig{
		IfName:            intf,
		Password:          "",
		EnablePm:          2, // shouldn't be enabled by default
		SerialNumber:      sn,
		AdminState:        1,
		OnuDhcpMode:       1,
		OnuIPAddress:      "",
		OnuIPMask:         "",
		OnuDefaultGateway: "",
	}
	//fmt.Println(o)
	return o
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (o *OnuConfig) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(o)
	if err != nil {
		return "", data
	}
	return o.IfName, data
}
