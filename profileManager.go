package goPon

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
                        MsanOnuInfoTable struct {
                                MsanOnuInfoEntry []OnuInfo `json:"msanOnuInfoEntry"`
                        } `json:"msanOnuInfoTable"`
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
