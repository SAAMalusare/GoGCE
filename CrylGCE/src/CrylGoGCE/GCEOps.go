package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v0.beta"
	"log"
	"strings"
)

func main() {

	/*
		Suhas:- Get all parameters from commandline
		-----------------------------------------------------------------------------------------------

		go run GCEOps.go -pn=crygce-test -zn=us-central1-a -in=mynewinstance001 -mt=f1-micro -si=centos-cloud/global/images/centos-6-v20170620 -ops=I
		go run GCEOps.go -pn=crygce-test -ops= -hin=healthcheck007
		go run GCEOps.go -pn=crygce-test -zn=us-central1-a -ops=H -hin=myhc001

	*/

	pname := flag.String("pn", "pn", "Provide Product name")
	zone := flag.String("zn", "zn", "Provide zone, e.g. europe-west1-c")
	instanceName := flag.String("in", "itn", "Provide new VM name")
	machineType := flag.String("mt", "mt", "Provide type of machine, e.g. f1-micro")
	scrimage := flag.String("si", "srcimg", "Provide source Image, e.g. centos-cloud/global/images/centos-6-v20170620")
	hinstanceName := flag.String("hin", "hitn", "Provide new Healthcheck name")
	ops := flag.String("ops", "Options", "Select operations to perform, I.. for new instance, H for creating new Healthcheck, HC for Healthcheck response of an existing Healthcheck")
	//---------------------------------------------------------------------------------------------------
	flag.Parse()
	/*
		fmt.Println(*ops)
		fmt.Println(*pname)
		fmt.Println(*zone)
		fmt.Println(*instanceName)
		fmt.Println(*machineType)
		fmt.Println(*scrimage)
		//fmt.Println(myprojURL) //projects/centos-cloud/global/images/centos-6-v20170620
	*/

	myprojURL := "https://www.googleapis.com/compute/v1/projects/" + *pname

	if *ops == "" {
		*ops = "x"
	}

	if *ops == "I" {

		if (*pname == "") || (*zone == "") || (*instanceName == "") || (*machineType == "") || (*scrimage == "") {
			fmt.Println("Error:Please provide mandatory argumnts for creating new VM instance..")
			Helpme()
		} else {
			fmt.Println("Creating New VM Instance.... using Insert command")
			createVMinstance(strings.ToLower(*pname), *zone, strings.ToLower(*instanceName), *machineType, *scrimage)
		}

	} else if *ops == "H" {

		if (*pname == "") || (*zone == "") || (*hinstanceName == "") {
			fmt.Println("Error:Please provide mandatory argumnts for creating new Healthcheck instance..")
			Helpme()
		} else {
			fmt.Println("Creating New HealthCheck Instance.... using Insert command")
			createHealthchk(strings.ToLower(*pname), *zone, myprojURL, strings.ToLower(*hinstanceName))
		}

	} else if *ops == "B" {
		if (*pname == "") || (*zone == "") || (*instanceName == "") || (*machineType == "") || (*scrimage == "") || (*hinstanceName == "") {
			fmt.Println("Error:Please provide mandatory argumnts for creating new VM & Healthcheck instance..")
			Helpme()
		} else {
			fmt.Println("Creating New VM Instance and Healthcheck instance")
			createVMinstance(strings.ToLower(*pname), *zone, strings.ToLower(*instanceName), *machineType, *scrimage)
			createHealthchk(strings.ToLower(*pname), *zone, myprojURL, strings.ToLower(*hinstanceName))
		}
	} else {
		if *hinstanceName == "" {
			fmt.Println("Error:Please provide Healthcheck instance name..")
			Helpme()
		} else {
			fmt.Println("Checking HealthCheck response for " + *hinstanceName)
			Hcresponse := HealthStatusGet(strings.ToLower(*pname), strings.ToLower(*hinstanceName))
			fmt.Println(Hcresponse)
		}
	}

}

func createVMinstance(pname string, zone string, instanceName string, machineType string, scrimage string) {

	ctx := context.Background()
	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	myprojURL := "https://www.googleapis.com/compute/v1/projects/" + pname
	imageURL := "https://www.googleapis.com/compute/v1/projects/" + scrimage //centos-cloud/global/images/centos-6-v20170620"

	fmt.Println(myprojURL)
	fmt.Println(imageURL)

	//os.Exit(0)

	rb := &compute.Instance{

		Name:        instanceName,
		Description: "created by Golang",
		MachineType: myprojURL + "/zones/" + zone + "/machineTypes/" + machineType,

		Disks: []*compute.AttachedDisk{
			{
				AutoDelete: true,
				Boot:       true,
				Type:       "PERSISTENT",
				InitializeParams: &compute.AttachedDiskInitializeParams{
					DiskName:    "suhasdisc",
					SourceImage: imageURL,
				},
			},
		},

		NetworkInterfaces: []*compute.NetworkInterface{
			&compute.NetworkInterface{
				AccessConfigs: []*compute.AccessConfig{
					&compute.AccessConfig{
						Type: "ONE_TO_ONE_NAT",
						Name: "External NAT",
					},
				},
				Network: myprojURL + "/global/networks/default",
			},
		},

		ServiceAccounts: []*compute.ServiceAccount{
			{
				Email: "default",
				Scopes: []string{
					compute.DevstorageFullControlScope,
					compute.ComputeScope,
				},
			},
		},

		Tags: &compute.Tags{
			Items: []string{"http-server"},
		},
	}

	resp, err := computeService.Instances.Insert(pname, zone, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	if resp.HTTPStatusCode == 200 {

		fmt.Print("Request placed successfully....")

		for {

			currstatus := getresponse(pname, zone, instanceName)

			fmt.Println("Current Status:- " + currstatus)

			if currstatus == "RUNNING" {

				resp, err := computeService.Instances.Get(pname, zone, instanceName).Context(ctx).Do()
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("External IP = " + resp.NetworkInterfaces[0].AccessConfigs[0].NatIP)
				fmt.Println("Internal IP = " + resp.NetworkInterfaces[0].NetworkIP)

				break
			}
		}
	} else {
		fmt.Print("Request failed....")
		fmt.Println(resp.HTTPStatusCode)
	}

}

func getresponse(pname string, zone string, instance string) string {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := computeService.Instances.Get(pname, zone, instance).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	currstatus := resp.Status
	return currstatus
}

func createHealthchk(project string, zone string, myprojURL string, hcName string) {
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)

	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	rb := &compute.HealthCheck{
		Name:               hcName,
		TimeoutSec:         5,
		UnhealthyThreshold: 2,
		HealthyThreshold:   2,
		SelfLink:           myprojURL + "/global/httpHealthChecks",
		Type:               "http",

		HttpHealthCheck: &compute.HTTPHealthCheck{
			Port:        80,
			RequestPath: "/",
			Host:        "",
		},
	}

	resp, err := computeService.HealthChecks.Insert(project, rb).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Change code below to process the `resp` object:
	responsecode := resp.HttpErrorStatusCode

	if responsecode == 0 {
		fmt.Println("HealthCheck: " + resp.Name + " has been created at " + resp.InsertTime)
		fmt.Print("HealthChkResponse: ")
		fmt.Println(HealthStatusGet(project, hcName))
	}
}

func HealthStatusGet(project string, hcName string) int {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := computeService.HealthChecks.Get(project, hcName).Context(ctx).Do()

	if err != nil {
		log.Fatal(err)
	}

	return resp.HTTPStatusCode
}

func Helpme() {
	fmt.Println("-----------------------------------------------------------------------")
	fmt.Println("Usage:")
	fmt.Println("-----------------------------------------------------------------------")

	fmt.Println("GCEOps.go -pn=PROJECT_NAME -zn=ZONE -in=VMNAME -mt=MACHINE_TYPE -si=SOURCE_IMAGE_PATH -ops=TYPE_OF_OPERATION")
	fmt.Println("Where:")
	fmt.Println("-pn (Type string): PROJECT NAME")
	fmt.Println("-zn (Type string): ZONE, e.g. us-central1-a")
	fmt.Println("-in (Type string): INSTANCE NAME (VM NAME), e.g. mynewvm")
	fmt.Println("-mt (Type string): MACHINE TYPE, e.g. f1-micro")
	fmt.Println("-si (Type string): PATH TO SOURCE IMAGE, e.g. centos-cloud/global/images/centos-6-v20170620")
	fmt.Println("-ops (Type string): TYPE OF OPERATIONS e.g. I for create new instance, H for create new HealthCheck, B for creating new VM & HealthCheck instance and default is HealthCheck response.")
	fmt.Println("-hin (Type string): INSTANCE NAME (VM NAME), e.g. mynewhealthcheck")
	fmt.Println("-----------------------------------------------------------------------")
}

func ConvertToLower(str2conv string) string {
	return strings.ToLower(str2conv)
}
