package gopon

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// required to maintain order, variables initialized once as constants
var Endpoints = []string{
	serviceProfiles,
	flowProfiles,
	vlanProfiles,
	igmpProfiles,
	securityProfiles,
	onuFlowProfiles,
	onuTcontProfiles,
	onuVlanProfiles,
	onuVlanRules,
	onuIgmpProfiles,
	l2cpProfiles,
	onuBlacklist, // GET only
	onuConfig,
	onuProfiles,
}

// to post to the endpoint, use the endpoint at a key to get the endpoint entry string
var endpointEntry = map[string]string{
	serviceProfiles:  "msanServiceProfileEntry",
	flowProfiles:     "msanServiceFlowProfileEntry",
	vlanProfiles:     "msanVlanProfileEntry",
	igmpProfiles:     "msanMulticastProfileEntry",
	securityProfiles: "msanSecurityProfileEntry",
	onuFlowProfiles:  "msanOnuFlowProfileEntry",
	onuTcontProfiles: "msanOnuTcontProfileEntry",
	onuVlanProfiles:  "msanOnuVlanProfileEntry",
	onuVlanRules:     "msanOnuVlanProfileRuleEntry",
	onuIgmpProfiles:  "msanOnuMulticastProfileEntry",
	l2cpProfiles:     "msanL2cpProfileEntry",
	onuConfig:        "msanOnuCfgEntry",
	onuProfiles:      "msanServicePortProfileEntry",
}

// checkHost checks if host accepts tcp connection on hardcoded https port
func CheckHost(host string, timeout int) (err error) {
	to := time.Duration(timeout) * time.Second
	_, err = net.DialTimeout("tcp", host+":443", to)
	return
}

func RestGetProfiles(host string, ep string) ([]byte, error) {
	reqUrl := fmt.Sprintf("https://%s/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/%s", host, ep)
	fmt.Printf("------------\nGET Request: %s\n------------\n", reqUrl)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", auth)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// host is ip address, ep is endpoint, epi is endpoint [*]unit of entry, name is profile name
// could return http Response directly
func RestPostProfile(host, ep, name string, data []byte) (string, error) {
	reqUrl := fmt.Sprintf("https://%s/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/%s/%s=%s", host, ep, endpointEntry[ep], name)
	fmt.Printf("------------\nPOST Request: %s\n------------\n", reqUrl)
	//fmt.Println(string(data))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", auth)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(string(resp.Status))
		return "", err
	}
	defer resp.Body.Close()
	//fmt.Println(string(resp.Status))

	return string(resp.Status), nil
}

func RestPatchProfile(host, ep, name string, data []byte) (string, error) {
	reqUrl := fmt.Sprintf("https://%s/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/%s/%s=%s", host, ep, endpointEntry[ep], name)
	fmt.Printf("------------\nPATCH Request: %s\n------------\n", reqUrl)
	//fmt.Println(string(data))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("PATCH", reqUrl, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", auth)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println(string(resp.Status))
		return "", err
	}
	defer resp.Body.Close()
	//fmt.Println(string(resp.Status))

	return string(resp.Status), nil
}

// returning http Response requires the profileManager to import HTTP
func RestDeleteProfile(host string, ep string, name string) (string, error) {
	reqUrl := fmt.Sprintf("https://%s/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/%s/%s=%s", host, ep, endpointEntry[ep], name)
	fmt.Printf("------------\nDELETE Request: %s\n------------\n", reqUrl)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("DELETE", reqUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Cookie", auth)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		//fmt.Println(resp.Status)
		return resp.Status, err
	}
	defer resp.Body.Close()
	//fmt.Println(string(resp.Status))
	return resp.Status, nil
}
