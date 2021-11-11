package goPon

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type OnuBlacklist struct {
	IfName       string `json:"msanOnuBlackListIfName"`
	SerialNumber string `json:"msanOnuBlackListSerialNumber"`
	Password     string `json:"msanOnuBlackListPassword"`
	Cause        int    `json:"msanOnuBlackListCause"`
}

type OnuBlacklistList struct {
	Entry []*OnuBlacklist
}

func (bl *OnuBlacklist) GetOltInterface() string {
	return bl.IfName
}

func (bl *OnuBlacklist) GetOnuSerialNumber() string {
	return bl.SerialNumber
}

func (bl *OnuBlacklist) GetPassword() string {
	return bl.Password
}

func (bl *OnuBlacklist) GetBlCause() string {
	switch {
	case bl.Cause == 1 {
		return "Invalid"
	case bl.Cause == 2 {
		return "SN Not Known"
	case bl.Cause == 3 {
		return "Password Mismatch"
	case bl.Cause == 6 {
		return "PON Link Mismatch"
	default:
		return "Unknown"
	}
}

var OnuBlacklistHeaders = []string{
	"Interface",
	"Serial Number",
	"Pass",
	"Cause",
}

// ListEssentialParams returns a map of the essential OnuFlowProfile parameters
func (bl *OnuBlacklist) ListEssentialParams() map[string]interface{} {
	var EssentialOnuBlacklist = map[string]interface{}{
		OnuBlacklistHeaders[0]: bl.GetOltInterface(),
		OnuBlacklistHeaders[1]: bl.GetOnuSerialNumber(),
		OnuBlacklistHeaders[2]: bl.GetPassword(),
		OnuBlacklistHeaders[3]: bl.GetBlCause(),
	}

	return EssentialOnuBlacklist
}

// Separate is a method to maintain backward-compatability
func (bll *OnuBlacklistList) Separate() []*OnuBlacklist {
	var entry *OnuBlacklist
	var list []*OnuBlacklist
	for _, e := range bll.Entry {
		entry = e
		list = append(list, entry)
	}
	return list
}

func (bll *OnuBlacklistList) Tabwrite() {
	fmt.Println("|| ONU Blacklist List ||")
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range OnuBlacklistHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuBlacklistHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	obl := bll.Separate()
	for _, bl := range obl {
		// first get the data as a map
		l := bl.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range OnuBlacklistHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuBlacklistHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}
