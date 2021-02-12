package main

import (
	"flag"
	"fmt"
	"time"

	"iskratel/gopon"
)

var (
	helpFlag = flag.Bool("h", false, "Show this help")
	//sleepFlag = flag.Duration("s", time.Millisecond*0, "Specify breaks between actions")
)

const usage = "`gopon` [options] <olt_ip>"

func main() {
	flag.Parse()

	if *helpFlag || flag.NArg() < 1 {
		fmt.Println(usage)
		flag.PrintDefaults()
		return
	}
	var err error
	start := time.Now()
	host := flag.Args()[0]
	err = serviceProfileDemo(host)
	fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
	err = flowProfileDemo(host)
	fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
	err = vlanProfileDemo(host)
	fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
	err = onuFlowProfileDemo(host)
	fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
	err = onuTcontProfileDemo(host)
	fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
	err = securityProfileDemo(host)
	fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
	err = multicastProfileDemo(host)
	fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
	err = onuMulticastProfileDemo(host)
	fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
}

func serviceProfileDemo(host string) error {
	var err error
	// the root object is a "LumiaOlt" which holds and active and cached copy of the "IskratelMsan" complete nested data structure
	olt := gopon.NewLumiaOlt(host)
	// the OLT object has a method to check whether it is reachable on port 443
	// failure here could be due to TCP/IP, or due to the HTTP service being disabled on the OLT
	if !olt.HostIsReachable() {
		fmt.Printf("Host %s is not reachable\n", host)
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	// the top level data structure for provisioning services on an OLT is represented by a "Service Profile"
	// we will perform a GET Request to retrieve all currently configured Service Profiles on the OLT
	// a separate object holds a list of the individual profile objects to allow group tabwrite methods
	var spl *gopon.ServiceProfileList
	spl, err = olt.GetServiceProfiles()
	if err != nil {
		//fmt.Printf("Error getting Service Profiles from OLT: %v\n", err)
		return err
	}
	spl.Tabwrite()
	/*
			// we can extract each ServiceProfile using a "Separate" Method
			var sps []*gopon.ServiceProfile
			sps = spl.Separate()
			fmt.Printf("--> Number of ServiceProfile entries in ServiceProfileList: %d\n", len(sps))

		// or we can iterate over the list object accessing each "Entry"
		// Tabwrite method simply shows we can access to the root object
		spName := "102_DATA"
		for _, sp := range spl.Entry {
			if sp.Name == spName {
				sp.Tabwrite()
			}
		}
	*/
	// the simplest way to test the POST method is to "copy" an existing service profile and change its name
	newSpName := "_TEST"
	var nsp *gopon.ServiceProfile
	for _, sp := range spl.Entry {
		newSpName = sp.Name + newSpName
		nsp, err = sp.Copy(newSpName)
		if err != nil {
			fmt.Printf("Error copying Service Profile: %v\n", err)
		} else {
			break
		}
	}
	err = olt.PostServiceProfile(nsp.GenerateJson())
	if err != nil {
		//fmt.Printf("Error Posting Service Profile '%s': %v\n", nsp.GetName(), err)
		return err
	}

	// we will get a new service profile list to check that this profile was added
	nspl, err := olt.GetServiceProfiles()
	if err != nil {
		//fmt.Printf("Error getting Service Profiles from OLT: %v\n", err)
		return err
	}
	for _, sp := range nspl.Entry {
		if sp.Name == newSpName {
			sp.Tabwrite()
			break
		}
	}

	// then we will delete this profile by name
	err = olt.DeleteServiceProfile(newSpName)
	if err != nil {
		//fmt.Printf("Error deleting Service Profiles from OLT: %v\n", err)
		return err
	}

	// and get the list to make sure it is gone
	nnspl, err := olt.GetServiceProfiles()
	if err != nil {
		//fmt.Printf("Error getting Service Profiles from OLT: %v\n", err)
		return err
	} else {
		fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(spl.Entry), len(nspl.Entry), len(nnspl.Entry))
	}

	/*
	   	// we can shortcut this method by getting a service profile directly by name
	   	// but there is no real efficiency as in the background the entire list is got and filtered
	   	namedSp, err := olt.GetServiceProfileByName(spName)
	   	fmt.Printf("Error getting Service Profile '%s' By Name from OLT: %v\n", spName, err)
	   	namedSp.Tabwrite()
	   	_, err = olt.GetServiceProfileByName("notExists")
	   	fmt.Printf("Error getting Service Profile '%s' By Name from OLT: %v\n", "notExists", err)
	   	_, err = olt.GetServiceProfileByName("")
	   	fmt.Printf("Error getting Service Profile '%s' By Name from OLT: %v\n", "", err)

	   Host 10.5.100.10 is reachable

	   GET Request: msanServiceProfileTable
	   https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceProfileTable
	   ----
	   2021/02/02 09:26:04 &{{{{[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]}}}}
	   Error getting Service Profiles from OLT: <nil>
	   Name          Flow Profile  VLAN Profile  ONU Flow Profile  ONU TCONT Profile  ONU VLAN Profile  Virtual GEM Port  ONU TP Type  Security Profile  IGMP Profile  ONU IGMP Profile  L2CP Profile  DHCP RA  PPPoE IA
	   ----          ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------  --------
	   101_CWMP      MB            101           101               T5I6__M-20M        A101              2                 IPHOST                                                                       []       []
	   102_DATA      MB            102           102               T5I2__M-100M                         10                VEIP                                                                         []       []
	   102_DATA_Acc  MB            102           102               T5I2__M-100M       A102              11                UNI                                                                          []       []
	   DEFAULT       DEFAULT       DEFAULT                                                              1                 VEIP                                                                         []       []
	   ----          ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------  --------
	   Name      Flow Profile  VLAN Profile  ONU Flow Profile  ONU TCONT Profile  ONU VLAN Profile  Virtual GEM Port  ONU TP Type  Security Profile  IGMP Profile  ONU IGMP Profile  L2CP Profile  DHCP RA  PPPoE IA
	   ----      ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------  --------
	   102_DATA  MB            102           102               T5I2__M-100M                         10                VEIP                                                                         []       []
	   ----      ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------  --------

	   GET Request: msanServiceProfileTable
	   https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceProfileTable
	   ----
	   2021/02/02 09:26:10 &{{{{[{101_CWMP MB  101   101 A101  T5I6__M-20M 2 2 AAAA 0 0 0 0 5   1 0 0 1332 2   0 5 1   1} {102_DATA MB  102   102   T5I2__M-100M 10 1 AAAA 0 0 0 0 5   1 0 0 1332 2   0 5 1   1} {102_DATA_Acc MB  102   102 A102  T5I2__M-100M 11 3 QAAA 0 0 0 0 5   1 0 0 1332 2   0 5 1   2} {DEFAULT DEFAULT  DEFAULT       1 1 AAAA 0 0 0 0 5   1 0 0 1332 2   0 5 1   2}]} {[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]}}}}
	   Error getting Service Profile '102_DATA' By Name from OLT: <nil>
	   Name      Flow Profile  VLAN Profile  ONU Flow Profile  ONU TCONT Profile  ONU VLAN Profile  Virtual GEM Port  ONU TP Type  Security Profile  IGMP Profile  ONU IGMP Profile  L2CP Profile  DHCP RA  PPPoE IA
	   ----      ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------  --------
	   102_DATA  MB            102           102               T5I2__M-100M                         10                VEIP                                                                         []       []
	   ----      ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------  --------

	   GET Request: msanServiceProfileTable
	   https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceProfileTable
	   ----
	   2021/02/02 09:26:15 &{{{{[{101_CWMP MB  101   101 A101  T5I6__M-20M 2 2 AAAA 0 0 0 0 5   1 0 0 1332 2   0 5 1   1} {102_DATA MB  102   102   T5I2__M-100M 10 1 AAAA 0 0 0 0 5   1 0 0 1332 2   0 5 1   1} {102_DATA_Acc MB  102   102 A102  T5I2__M-100M 11 3 QAAA 0 0 0 0 5   1 0 0 1332 2   0 5 1   2} {DEFAULT DEFAULT  DEFAULT       1 1 AAAA 0 0 0 0 5   1 0 0 1332 2   0 5 1   2}]} {[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]} {[]}}}}
	   Error getting Service Profile 'notExists' By Name from OLT: Key not found
	   Error getting Service Profile '' By Name from OLT: Incorrect input supplied
	   Elapsed time of Operation: 11.131478242s

	*/
	return nil
}

func flowProfileDemo(host string) error {
	var err error
	olt := gopon.NewLumiaOlt(host)

	if !olt.HostIsReachable() {
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	var fpl *gopon.FlowProfileList

	fpl, err = olt.GetFlowProfiles()
	if err != nil {
		return err
	}
	fpl.Tabwrite()
	// POST demo copies the first profile in the list, appends "_TEST" to the name and reposts it
	newName := "_TEST"
	var newProfile *gopon.FlowProfile
	for _, fp := range fpl.Entry {
		newName = fp.Name + newName
		newProfile, err = fp.Copy(newName)
		if err != nil {
			fmt.Printf("Error copying Flow Profile: %v\n", err)
		} else {
			break
		}
	}
	err = olt.PostFlowProfile(newProfile.GenerateJson())
	if err != nil {
		return err
	}

	// we will get a new service profile list to check that this profile was added
	nfpl, err := olt.GetFlowProfiles()
	if err != nil {
		return err
	}
	for _, fp := range nfpl.Entry {
		if fp.Name == newName {
			fp.Tabwrite()
			break
		}
	}

	// then we will delete this profile by name
	err = olt.DeleteFlowProfile(newName)
	if err != nil {
		return err
	}

	// and get the list to make sure it is gone
	nnfpl, err := olt.GetFlowProfiles()
	if err != nil {
		return err
	}
	fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(fpl.Entry), len(nfpl.Entry), len(nnfpl.Entry))
	return nil
}

func vlanProfileDemo(host string) error {
	var err error
	olt := gopon.NewLumiaOlt(host)

	if !olt.HostIsReachable() {
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	var vpl *gopon.VlanProfileList

	vpl, err = olt.GetVlanProfiles()
	if err != nil {
		return err
	}
	vpl.Tabwrite()
	// POST demo copies the first profile in the list, appends "_TEST" to the name and reposts it
	newName := "_TEST"
	var newProfile *gopon.VlanProfile
	for _, vp := range vpl.Entry {
		newName = vp.Name + newName
		newProfile, err = vp.Copy(newName)
		if err != nil {
			fmt.Printf("Error copying Vlan Profile: %v\n", err)
		} else {
			break
		}
	}
	err = olt.PostVlanProfile(newProfile.GenerateJson())
	if err != nil {
		return err
	}

	// we will get a new service profile list to check that this profile was added
	nvpl, err := olt.GetVlanProfiles()
	if err != nil {
		return err
	}
	for _, vp := range nvpl.Entry {
		if vp.Name == newName {
			vp.Tabwrite()
			break
		}
	}

	// then we will delete this profile by name
	err = olt.DeleteVlanProfile(newName)
	if err != nil {
		return err
	}

	// and get the list to make sure it is gone
	nnvpl, err := olt.GetVlanProfiles()
	if err != nil {
		return err
	}
	fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(vpl.Entry), len(nvpl.Entry), len(nnvpl.Entry))
	return nil
}

func onuFlowProfileDemo(host string) error {
	var err error
	olt := gopon.NewLumiaOlt(host)

	if !olt.HostIsReachable() {
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	var ofpl *gopon.OnuFlowProfileList

	ofpl, err = olt.GetOnuFlowProfiles()
	if err != nil {
		return err
	}
	ofpl.Tabwrite()
	// POST demo copies the first profile in the list, appends "_TEST" to the name and reposts it
	newName := "_TEST"
	var newProfile *gopon.OnuFlowProfile
	for _, ofp := range ofpl.Entry {
		newName = ofp.Name + newName
		newProfile, err = ofp.Copy(newName)
		if err != nil {
			fmt.Printf("Error copying Onu Flow Profile: %v\n", err)
		} else {
			break
		}
	}
	err = olt.PostOnuFlowProfile(newProfile.GenerateJson())
	if err != nil {
		return err
	}

	// we will get a new service profile list to check that this profile was added
	nofpl, err := olt.GetOnuFlowProfiles()
	if err != nil {
		return err
	}
	for _, ofp := range nofpl.Entry {
		if ofp.Name == newName {
			ofp.Tabwrite()
			break
		}
	}

	// then we will delete this profile by name
	err = olt.DeleteOnuFlowProfile(newName)
	if err != nil {
		return err
	}

	// and get the list to make sure it is gone
	nnofpl, err := olt.GetOnuFlowProfiles()
	if err != nil {
		return err
	}
	fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(ofpl.Entry), len(nofpl.Entry), len(nnofpl.Entry))
	return nil
}

func onuTcontProfileDemo(host string) error {
	var err error
	olt := gopon.NewLumiaOlt(host)

	if !olt.HostIsReachable() {
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	var otpl *gopon.OnuTcontProfileList

	otpl, err = olt.GetOnuTcontProfiles()
	if err != nil {
		return err
	}
	otpl.Tabwrite()
	// POST demo copies the first profile in the list, appends "_TEST" to the name and reposts it
	newName := "_TEST"
	var newProfile *gopon.OnuTcontProfile
	for _, otp := range otpl.Entry {
		newName = otp.Name + newName
		newProfile, err = otp.Copy(newName)
		if err != nil {
			fmt.Printf("Error copying Onu Tcont Profile: %v\n", err)
		} else {
			break
		}
	}
	err = olt.PostOnuTcontProfile(newProfile.GenerateJson())
	if err != nil {
		return err
	}

	// we will get a new service profile list to check that this profile was added
	notpl, err := olt.GetOnuTcontProfiles()
	if err != nil {
		return err
	}
	for _, otp := range notpl.Entry {
		if otp.Name == newName {
			otp.Tabwrite()
			break
		}
	}

	// then we will delete this profile by name
	err = olt.DeleteOnuTcontProfile(newName)
	if err != nil {
		return err
	}

	// and get the list to make sure it is gone
	nnotpl, err := olt.GetOnuTcontProfiles()
	if err != nil {
		return err
	}
	fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(otpl.Entry), len(notpl.Entry), len(nnotpl.Entry))
	return nil
}

func securityProfileDemo(host string) error {
	var err error
	olt := gopon.NewLumiaOlt(host)

	if !olt.HostIsReachable() {
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	var spl *gopon.SecurityProfileList

	spl, err = olt.GetSecurityProfiles()
	if err != nil {
		return err
	}
	spl.Tabwrite()
	// POST demo copies the first profile in the list, appends "_TEST" to the name and reposts it
	newName := "_TEST"
	var newProfile *gopon.SecurityProfile
	for _, sp := range spl.Entry {
		newName = sp.Name + newName
		newProfile, err = sp.Copy(newName)
		if err != nil {
			fmt.Printf("Error copying Security Profile: %v\n", err)
		} else {
			break
		}
	}
	err = olt.PostSecurityProfile(newProfile.GenerateJson())
	if err != nil {
		return err
	}

	// we will get a new service profile list to check that this profile was added
	nspl, err := olt.GetSecurityProfiles()
	if err != nil {
		return err
	}
	for _, sp := range nspl.Entry {
		if sp.Name == newName {
			sp.Tabwrite()
			break
		}
	}

	// then we will delete this profile by name
	err = olt.DeleteSecurityProfile(newName)
	if err != nil {
		return err
	}

	// and get the list to make sure it is gone
	nnspl, err := olt.GetSecurityProfiles()
	if err != nil {
		return err
	}
	fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(spl.Entry), len(nspl.Entry), len(nnspl.Entry))
	return nil
}

func multicastProfileDemo(host string) error {
	var err error
	olt := gopon.NewLumiaOlt(host)

	if !olt.HostIsReachable() {
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	var mpl *gopon.IgmpProfileList

	mpl, err = olt.GetMulticastProfiles()
	if err != nil {
		return err
	}
	mpl.Tabwrite()
	// POST demo copies the first profile in the list, appends "_TEST" to the name and reposts it
	newName := "_TEST"
	var newProfile *gopon.IgmpProfile
	for _, mp := range mpl.Entry {
		newName = mp.Name + newName
		newProfile, err = mp.Copy(newName)
		if err != nil {
			fmt.Printf("Error copying Multicast Profile: %v\n", err)
		} else {
			break
		}
	}
	err = olt.PostMulticastProfile(newProfile.GenerateJson())
	if err != nil {
		return err
	}

	// we will get a new service profile list to check that this profile was added
	nmpl, err := olt.GetMulticastProfiles()
	if err != nil {
		return err
	}
	for _, mp := range nmpl.Entry {
		if mp.Name == newName {
			mp.Tabwrite()
			break
		}
	}

	// then we will delete this profile by name
	err = olt.DeleteMulticastProfile(newName)
	if err != nil {
		return err
	}

	// and get the list to make sure it is gone
	nnmpl, err := olt.GetMulticastProfiles()
	if err != nil {
		return err
	}
	fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(mpl.Entry), len(nmpl.Entry), len(nnmpl.Entry))
	return nil
}

func onuMulticastProfileDemo(host string) error {
	var err error
	olt := gopon.NewLumiaOlt(host)

	if !olt.HostIsReachable() {
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}

	var ompl *gopon.OnuIgmpProfileList

	ompl, err = olt.GetOnuMulticastProfiles()
	if err != nil {
		return err
	}
	ompl.Tabwrite()
	// POST demo copies the first profile in the list, appends "_TEST" to the name and reposts it
	newName := "_TEST"
	var newProfile *gopon.OnuIgmpProfile
	for _, omp := range ompl.Entry {
		newName = omp.Name + newName
		newProfile, err = omp.Copy(newName)
		if err != nil {
			fmt.Printf("Error copying Onu Multicast Profile: %v\n", err)
		} else {
			break
		}
	}
	err = olt.PostOnuMulticastProfile(newProfile.GenerateJson())
	if err != nil {
		return err
	}

	// we will get a new service profile list to check that this profile was added
	nompl, err := olt.GetOnuMulticastProfiles()
	if err != nil {
		return err
	}
	for _, omp := range nompl.Entry {
		if omp.Name == newName {
			omp.Tabwrite()
			break
		}
	}

	// then we will delete this profile by name
	err = olt.DeleteOnuMulticastProfile(newName)
	if err != nil {
		return err
	}

	// and get the list to make sure it is gone
	nnompl, err := olt.GetOnuMulticastProfiles()
	if err != nil {
		return err
	}
	fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(ompl.Entry), len(nompl.Entry), len(nnompl.Entry))
	return nil
}

/*
func subProfileDemo(host string) error {
	var err error
	// the root object is a "LumiaOlt" which holds and active and cached copy of the "IskratelMsan" complete nested data structure
	olt := gopon.NewLumiaOlt(host)
	// the OLT object has a method to check whether it is reachable on port 443
	// failure here could be due to TCP/IP, or due to the HTTP service being disabled on the OLT
	if !olt.HostIsReachable() {
		//fmt.Printf("Host %s is not reachable\n", host)
		return gopon.ErrNotReachable
	} else {
		fmt.Printf("Host %s is reachable\n", host)
	}




		// for now, create each new object directly
		// serialize the data and POST it to the NE
		nfp := gopon.NewFlowProfile("RestDemo")
		nfp.SetMatchBothVlanProfile()
		err = olt.PostFlowProfile(nfp.GenerateJson())
		fmt.Printf("Error posting Flow Profile: %v\n", err)

		nvp := gopon.NewVlanProfile("RestDemo")
		err = nvp.SetCVid([]int{200})
		fmt.Printf("Error setting C-VID range on VLAN Profile: %v\n", err)
		err = olt.PostVlanProfile(nfp.GenerateJson())
		fmt.Printf("Error posting VLAN Profile: %v\n", err)

		nofp := gopon.NewOnuFlowProfile("RestDemo")
		err = nofp.SetMatchUsCVlanIDRange([]int{200})
		fmt.Printf("Error setting US C-VID Range on ONU Flow Profile: %v\n", err)
		err = olt.PostOnuFlowProfile(nofp.GenerateJson())
		fmt.Printf("Error posting ONU Flow Profile: %v\n", err)

		notp := gopon.NewOnuTcontProfile("RestDemo")
		notp.SetTcontID(5)
		notp.SetFAM(1000000, 10000, 0)
		err = olt.PostOnuTcontProfile(notp.GenerateJson())
		fmt.Printf("Error posting ONU T-CONT Profile: %v\n", err)

		nsp := gopon.NewServiceProfile("RestDemo")
		nsp.SetFlowProfile(nfp.GetName())
		nsp.SetVlanProfile(nvp.GetName())
		nsp.SetOnuFlowProfile(nofp.GetName())
		nsp.SetOnuTcontProfile(notp.GetName())
		err = olt.PostServiceProfile(nsp.GenerateJson())
		fmt.Printf("Error posting Service Profile: %v\n", err)

		check, err := olt.GetServiceProfileByName(nsp.GetName())
		fmt.Printf("Error getting the New Service Profile by Name: %v\n", err)
		check.Tabwrite()

		// now we can delete the newly created profile (later we will apply this to a ONU)
		// we cannot delete a Service Profile if it is assigned to an ONU
		// we cannot delete a sub-profile that is being used in a service profile
		// we will assert this operation to error
		err = olt.DeleteVlanProfile(nvp.GetName())
		fmt.Printf("Error deleting VLAN Profile applied in our Service Profile: %v\n", err)
		err = olt.DeleteServiceProfile(nsp.GetName())
		fmt.Printf("Error deleting our new Service Profile: %v\n", err)
		err = olt.DeleteFlowProfile(nfp.GetName())
		fmt.Printf("Error deleting our new Flow Profile: %v\n", err)
		err = olt.DeleteVlanProfile(nvp.GetName())
		fmt.Printf("Error deleting our new VLAN Profile: %v\n", err)
		err = olt.DeleteOnuFlowProfile(nofp.GetName())
		fmt.Printf("Error deleting our new ONU Flow Profile: %v\n", err)
		err = olt.DeleteOnuTcontProfile(notp.GetName())
		fmt.Printf("Error deleting our new ONU T-CONT Profile: %v\n", err)

		// now we will check that these profiles have been removed
		_, err = olt.GetServiceProfileByName(nsp.GetName())
		fmt.Printf("Error getting our recently deleted Service Profile: %v\n", err)

		sps, err = olt.GetServiceProfiles()
		fmt.Printf("Error getting all Service Profiles from OLT: %v\n", err)
		for _, sp := range sps {
			sp.Tabwrite()
			fmt.Print("\n\n")
		}


	return nil
}
*/

/* 02-02-2021
{14:58}/home/lbnp/go/src/iskratel/gopon/cmd> go run main.go 10.5.100.10
Host 10.5.100.10 is reachable

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceProfileTable
------------
Name          Flow Profile  VLAN Profile  ONU Flow Profile  ONU TCONT Profile  ONU VLAN Profile  Virtual GEM Port  ONU TP Type  Security Profile  IGMP Profile  ONU IGMP Profile  L2CP Profile  DHCP RA             PPPoE IA
----          ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------             --------
101_CWMP      MB            101           101               T5I6__M-20M        A101              2                 IPHOST                                                                       [map[RateLimit:5]]  []
102_DATA      MB            102           102               T5I2__M-100M                         10                VEIP                                                                         [map[RateLimit:5]]  []
102_DATA_Acc  MB            102           102               T5I2__M-100M       A102              11                UNI                                                                          [map[RateLimit:5]]  []
DEFAULT       DEFAULT       DEFAULT                                                              1                 VEIP                                                                         [map[RateLimit:5]]  []
----          ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------             --------

POST Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceProfileTable/msanServiceProfileEntry=101_CWMP_TEST
------------
PostServiceProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceProfileTable
------------
Name           Flow Profile  VLAN Profile  ONU Flow Profile  ONU TCONT Profile  ONU VLAN Profile  Virtual GEM Port  ONU TP Type  Security Profile  IGMP Profile  ONU IGMP Profile  L2CP Profile  DHCP RA             PPPoE IA
----           ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------             --------
101_CWMP_TEST  MB            101           101               T5I6__M-20M        A101              2                 IPHOST                                                                       [map[RateLimit:5]]  []
----           ------------  ------------  ----------------  -----------------  ----------------  ----------------  -----------  ----------------  ------------  ----------------  ------------  -------             --------

DELETE Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceProfileTable/msanServiceProfileEntry=101_CWMP_TEST
------------
DeleteServiceProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceProfileTable
------------
First the OLT had 4 entries, then we added one and it had 5, now it has 4
[21.264054671s] Error running demo: <nil>
Host 10.5.100.10 is reachable

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceFlowProfileTable
------------
Name     UsMatchVlanProfile  DsMatchVlanProfile  UsMatchOther         DsMatchOther         UsHandling  DsHandling  QueuingPriority  SchedulingMode
----     ------------------  ------------------  ------------         ------------         ----------  ----------  ---------------  --------------
DEFAULT  false               false               []                   []                   []          []          0                Weighted
MB       true                true                [map[MatchUsAny:2]]  [map[MatchDsAny:2]]  []          []          0                Weighted
----     ------------------  ------------------  ------------         ------------         ----------  ----------  ---------------  --------------

POST Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceFlowProfileTable/msanServiceFlowProfileEntry=DEFAULT_TEST
------------
PostFlowProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceFlowProfileTable
------------
Name          UsMatchVlanProfile  DsMatchVlanProfile  UsMatchOther  DsMatchOther  UsHandling  DsHandling  QueuingPriority  SchedulingMode
----          ------------------  ------------------  ------------  ------------  ----------  ----------  ---------------  --------------
DEFAULT_TEST  false               false               []            []            []          []          0                Weighted
----          ------------------  ------------------  ------------  ------------  ----------  ----------  ---------------  --------------

DELETE Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceFlowProfileTable/msanServiceFlowProfileEntry=DEFAULT_TEST
------------
DeleteFlowProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanServiceFlowProfileTable
------------
First the OLT had 2 entries, then we added one and it had 3, now it has 2
[2m52.207477906s] Error running demo: <nil>
Host 10.5.100.10 is reachable

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanVlanProfileTable
------------
Name     C-Vid  C-Vid Native  S-Vid  S-Ethertype
----     -----  ------------  -----  -----------
101      [101]  -1            -1     34984
102      [102]  -1            -1     34984
DEFAULT  [1]    1             -1     34984
----     -----  ------------  -----  -----------

POST Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanVlanProfileTable/msanVlanProfileEntry=101_TEST
------------
PostVlanProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanVlanProfileTable
------------
Name      C-Vid  C-Vid Native  S-Vid  S-Ethertype
----      -----  ------------  -----  -----------
101_TEST  [101]  -1            -1     34984
----      -----  ------------  -----  -----------

DELETE Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanVlanProfileTable/msanVlanProfileEntry=101_TEST
------------
DeleteVlanProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanVlanProfileTable
------------
First the OLT had 3 entries, then we added one and it had 4, now it has 3
[5m7.544485254s] Error running demo: <nil>
Host 10.5.100.10 is reachable

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuFlowProfileTable
------------
Name  MatchUsC-VidRange  MatchUsCPcp  UsFlowPriority
----  -----------------  -----------  --------------
101   [101]              -1           0
102   [102]              -1           0
----  -----------------  -----------  --------------

POST Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuFlowProfileTable/msanOnuFlowProfileEntry=101_TEST
------------
PostOnuFlowProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuFlowProfileTable
------------
Name      MatchUsC-VidRange  MatchUsCPcp  UsFlowPriority
----      -----------------  -----------  --------------
101_TEST  [101]              -1           0
----      -----------------  -----------  --------------

DELETE Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuFlowProfileTable/msanOnuFlowProfileEntry=101_TEST
------------
DeleteOnuFlowProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuFlowProfileTable
------------
First the OLT had 2 entries, then we added one and it had 3, now it has 2
[5m42.810176248s] Error running demo: <nil>
Host 10.5.100.10 is reachable

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuTcontProfileTable
------------
Name          Description    ID  Type  Fixed  Assured  Max
----          -----------    --  ----  -----  -------  ---
New           T5I1__M-256Kb  1   5     0Kb    0Kb      256Kb
T5I2__M-100M  T5I2__M-98Mb   2   5     0Kb    0Kb      98Mb
T5I6__M-20M   T5I6__M-20Mb   6   5     0Kb    0Kb      20Mb
----          -----------    --  ----  -----  -------  ---

POST Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuTcontProfileTable/msanOnuTcontProfileEntry=New_TEST
------------
PostOnuTcontProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuTcontProfileTable
------------
Name      Description    ID  Type  Fixed  Assured  Max
----      -----------    --  ----  -----  -------  ---
New_TEST  T5I1__M-256Kb  1   5     0Kb    0Kb      256Kb
----      -----------    --  ----  -----  -------  ---

DELETE Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuTcontProfileTable/msanOnuTcontProfileEntry=New_TEST
------------
DeleteOnuTcontProfile: 200 OK

GET Request: https://10.5.100.10/restconf/data/ISKRATEL-MSAN-MIB:ISKRATEL-MSAN-MIB/msanOnuTcontProfileTable
------------
First the OLT had 3 entries, then we added one and it had 4, now it has 3
[6m7.890997597s] Error running demo: <nil>

*/
