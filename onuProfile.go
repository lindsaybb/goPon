package gopon

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

// this data structure is very simplistic, the real logic is held by the lumiaOlt object

type OnuProfile struct {
	IfName             string `json:"ifName"`
	ServiceProfileName string `json:"msanServiceProfileName"`
}

type OnuProfileList struct {
	Entry []*OnuProfile
}

func NewOnuProfile(intf, spName string) *OnuProfile {
	o := &OnuProfile{
		IfName:             intf,
		ServiceProfileName: spName,
	}
	return o
}

func (op *OnuProfile) GenerateJson() (name string, data []byte) {
	data, err := json.Marshal(op)
	if err != nil {
		return "", data
	}
	return op.IfName, data
}

var OnuProfileHeaders = []string{
	"Interface",
	"ServiceProfile",
}

func (op *OnuProfile) ListEssentialParams() map[string]interface{} {
	var EssentialOnuProfile = map[string]interface{}{
		OnuProfileHeaders[0]: op.IfName,
		OnuProfileHeaders[1]: op.ServiceProfileName,
	}

	return EssentialOnuProfile
}

/*
func (opl *OnuProfileList) CombineSameInterfaces() map[string][]string {
	combine := make(map[string][]string)
	for _, op := range opl.Entry {
		combine[op.IfName] = append(combine[op.IfName], op.ServiceProfileName)
	}
	return combine
}
*/
//var OnuProfileMap map[string][]string

func (opl *OnuProfileList) Tabwrite() {
	// create the writer
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	// write tab-separated header values to tw buffer
	for _, v := range OnuProfileHeaders {
		fmt.Fprintf(tw, "%v\t", v)
	}
	fmt.Fprintf(tw, "\n")
	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")

	for _, op := range opl.Entry {
		// first get the data as a map
		l := op.ListEssentialParams()
		// iterate over the map using the header as string key
		for _, v := range OnuProfileHeaders {
			fmt.Fprintf(tw, "%v\t", l[v])
		}
		fmt.Fprintf(tw, "\n")
	}

	// write tab-separated spacers (-) reflecting the length of the headers
	for _, v := range OnuProfileHeaders {
		fmt.Fprintf(tw, "%v\t", fs(v))
	}
	fmt.Fprintf(tw, "\n")
	// calculate column width and print table from tw buffer
	tw.Flush()
}
