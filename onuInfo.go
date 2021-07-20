package gopon

import (
        "encoding/json"
)

type OnuInfo struct {
        IfName                      string `json:"msanOnuInfoIfName"`
        OnuInfoPrimaryStatus        int    `json:"msanOnuInfoOnuInfoPrimaryStatus"`
        EqualizationDelay           int    `json:"msanOnuInfoEqualizationDelay"`
        PowerLevel                  int    `json:"msanOnuInfoPowerLevel"`
        VendorID                    string `json:"msanOnuInfoVendorId"`
        Version                     string `json:"msanOnuInfoVersion"`
        TrafficManagementOption     int    `json:"msanOnuInfoTrafficManagementOption"`
        OperState                   int    `json:"msanOnuInfoOperState"`
        EquipmentID                 string `json:"msanOnuInfoEquipmentId"`
        OmccVersion                 string `json:"msanOnuInfoOmccVersion"`
        OnuHardwareType             int    `json:"msanOnuInfoOnuHardwareType"`
        HardwareRevision            int    `json:"msanOnuInfoHardwareRevision"`
        SecurityCapability          int    `json:"msanOnuInfoSecurityCapability"`
        TotalPriorityQueueNumber    int    `json:"msanOnuInfoTotalPriorityQueueNumber"`
        TotalTrafficSchedulerNumber int    `json:"msanOnuInfoTotalTrafficSchedulerNumber"`
        TotalGemPortNumber          int    `json:"msanOnuInfoTotalGemPortNumber"`
        TotalTcontNumber            int    `json:"msanOnuInfoTotalTcontNumber"`
        TotalEthernetUniNumber      int    `json:"msanOnuInfoTotalEthernetUniNumber"`
        TotalPotsUniNumber          int    `json:"msanOnuInfoTotalPotsUniNumber"`
        SysUpTime                   int    `json:"msanOnuInfoSysUpTime"`
        OnuImageInstance0Version    string `json:"msanOnuInfoOnuImageInstance0Version"`
        OnuImageInstance0Valid      int    `json:"msanOnuInfoOnuImageInstance0Valid"`
        OnuImageInstance0Activate   int    `json:"msanOnuInfoOnuImageInstance0Activate"`
        OnuImageInstance0Commit     int    `json:"msanOnuInfoOnuImageInstance0Commit"`
        OnuImageInstance1Version    string `json:"msanOnuInfoOnuImageInstance1Version"`
        OnuImageInstance1Valid      int    `json:"msanOnuInfoOnuImageInstance1Valid"`
        OnuImageInstance1Activate   int    `json:"msanOnuInfoOnuImageInstance1Activate"`
        OnuImageInstance1Commit     int    `json:"msanOnuInfoOnuImageInstance1Commit"`
        OnuMacAddress               string `json:"msanOnuInfoOnuMacAddress"`
        OnuDhcpMode                 int    `json:"msanOnuInfoOnuDhcpMode"`
        OnuIPAddress                string `json:"msanOnuInfoOnuIpAddress"`
        OnuIPMask                   string `json:"msanOnuInfoOnuIpMask"`
        OnuDefaultGateway           string `json:"msanOnuInfoOnuDefaultGateway"`
        OnuFastLeaveCapability      int    `json:"msanOnuInfoOnuFastLeaveCapability"`
        SerialNumber                string `json:"msanOnuInfoSerialNumber"`
        Password                    string `json:"msanOnuInfoPassword"`
        RxPower                     int    `json:"msanOnuInfoRxPower"`
        TxPower                     int    `json:"msanOnuInfoTxPower"`
        OltRxPower                  int    `json:"msanOnuInfoOltRxPower"`
        Temp                        int    `json:"msanOnuInfoTemp"`
}

type OnuInfoList struct {
        Entry []*OnuInfo
}

// IsUp returns whether the OperState is 1
func (o *OnuInfo) IsUp() bool {
        return o.OperState == 1
}

// GenerateJson serializes the data structure so it can be set with Restconf
func (o *OnuInfo) GenerateJson() (intf string, data []byte) {
        data, err := json.Marshal(o)
        if err != nil {
                return "", data
        }
        return o.IfName, data
}
