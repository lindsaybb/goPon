package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/lindsaybb/gopon"
)

var (
	helpFlag             = flag.Bool("h", false, "Show this help")
	allDemos             = flag.Bool("A", false, "Run all Profile demos")
	getOnly              = flag.Bool("G", false, "Only run the GET Operation on each Demo")
	serviceProfiles      = flag.Bool("sp", false, "Run the Service Profile Demo")
	flowProfiles         = flag.Bool("fp", false, "Run the Flow Profile Demo")
	vlanProfiles         = flag.Bool("vp", false, "Run the VLAN Profile Demo")
	onuFlowProfiles      = flag.Bool("ofp", false, "Run the ONU Flow Profile Demo")
	onuTcontProfiles     = flag.Bool("otp", false, "Run the ONU T-CONT Profile Demo")
	secuirtyProfiles     = flag.Bool("secp", false, "Run the Security Profile Demo")
	multicastProfiles    = flag.Bool("mp", false, "Run the Multicast Profile Demo")
	onuMulticastProfiles = flag.Bool("omp", false, "Run the Onu Multicast Profile Demo")
)

const usage = "`gopon service profile demo` [options] <olt_ip>"

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
	// the root object is a "LumiaOlt" which holds and active and cached copy of the "IskratelMsan" complete nested data structure
	olt := gopon.NewLumiaOlt(host)
	// the OLT object has a method to check whether it is reachable on port 443
	// failure here could be due to TCP/IP, or due to the HTTP service being disabled on the OLT
	if !olt.HostIsReachable() {
		fmt.Printf("Host %s is not reachable\n", host)
		return
	}

	if *allDemos {
		err = serviceProfileDemo(olt)
		fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		err = flowProfileDemo(olt)
		fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		err = vlanProfileDemo(olt)
		fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		err = onuFlowProfileDemo(olt)
		fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		err = onuTcontProfileDemo(olt)
		fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		err = securityProfileDemo(olt)
		fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		err = multicastProfileDemo(olt)
		fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		err = onuMulticastProfileDemo(olt)
		fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
	} else {
		if *serviceProfiles {
			err = serviceProfileDemo(olt)
			fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		}
		if *flowProfiles {
			err = flowProfileDemo(olt)
			fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		}
		if *vlanProfiles {
			err = vlanProfileDemo(olt)
			fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		}
		if *onuFlowProfiles {
			err = onuFlowProfileDemo(olt)
			fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		}
		if *onuTcontProfiles {
			err = onuTcontProfileDemo(olt)
			fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		}
		if *secuirtyProfiles {
			err = securityProfileDemo(olt)
			fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		}
		if *multicastProfiles {
			err = multicastProfileDemo(olt)
			fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		}
		if *onuMulticastProfiles {
			err = onuMulticastProfileDemo(olt)
			fmt.Printf("[%v] Error running demo: %v\n", time.Since(start), err)
		}
	}
}

func serviceProfileDemo(olt *gopon.LumiaOlt) error {
	// the top level data structure for provisioning services on an OLT is represented by a "Service Profile"
	// we will perform a GET Request to retrieve all currently configured Service Profiles on the OLT
	// a separate object holds a list of the individual profile objects to allow group tabwrite methods
	var err error
	var spl *gopon.ServiceProfileList
	spl, err = olt.GetServiceProfiles()
	if err != nil {
		return err
	}
	spl.Tabwrite()
	if *getOnly {
		return nil
	}

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
		return err
	}
	// we will get a new service profile list to check that this profile was added
	nspl, err := olt.GetServiceProfiles()
	if err != nil {
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
		return err
	}
	// and get the list to make sure it is gone
	nnspl, err := olt.GetServiceProfiles()
	if err != nil {
		return err
	}

	fmt.Printf("First the OLT had %d entries, then we added one and it had %d, now it has %d\n", len(spl.Entry), len(nspl.Entry), len(nnspl.Entry))
	return nil
}

func flowProfileDemo(olt *gopon.LumiaOlt) error {
	var err error
	var fpl *gopon.FlowProfileList

	fpl, err = olt.GetFlowProfiles()
	if err != nil {
		return err
	}
	fpl.Tabwrite()
	if *getOnly {
		return nil
	}
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

func vlanProfileDemo(olt *gopon.LumiaOlt) error {
	var err error
	var vpl *gopon.VlanProfileList

	vpl, err = olt.GetVlanProfiles()
	if err != nil {
		return err
	}
	vpl.Tabwrite()
	if *getOnly {
		return nil
	}
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

func onuFlowProfileDemo(olt *gopon.LumiaOlt) error {
	var err error
	var ofpl *gopon.OnuFlowProfileList

	ofpl, err = olt.GetOnuFlowProfiles()
	if err != nil {
		return err
	}
	ofpl.Tabwrite()
	if *getOnly {
		return nil
	}
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

func onuTcontProfileDemo(olt *gopon.LumiaOlt) error {
	var err error
	var otpl *gopon.OnuTcontProfileList

	otpl, err = olt.GetOnuTcontProfiles()
	if err != nil {
		return err
	}
	otpl.Tabwrite()
	if *getOnly {
		return nil
	}
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

func securityProfileDemo(olt *gopon.LumiaOlt) error {
	var err error
	var spl *gopon.SecurityProfileList

	spl, err = olt.GetSecurityProfiles()
	if err != nil {
		return err
	}
	spl.Tabwrite()
	if *getOnly {
		return nil
	}
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

func multicastProfileDemo(olt *gopon.LumiaOlt) error {
	var err error
	var mpl *gopon.IgmpProfileList

	mpl, err = olt.GetMulticastProfiles()
	if err != nil {
		return err
	}
	mpl.Tabwrite()
	if *getOnly {
		return nil
	}
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

func onuMulticastProfileDemo(olt *gopon.LumiaOlt) error {
	var err error
	var ompl *gopon.OnuIgmpProfileList

	ompl, err = olt.GetOnuMulticastProfiles()
	if err != nil {
		return err
	}
	ompl.Tabwrite()
	if *getOnly {
		return nil
	}
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
