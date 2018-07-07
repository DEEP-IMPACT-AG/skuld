// Copyright 2018 Deep Impact AG. All rights reserved.
// Use of this source code is governed by the Apache License Version 2.0
// that can be found in the LICENSE file.
package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/mitchellh/go-homedir"
	"strings"
	"fmt"
	"gopkg.in/ini.v1"
	"flag"
	"os"
	"path/filepath"
)

var (
	version = "dev"
	verbose *bool
)

const admDuration = 3600
const standardDuration = 36000

func main() {
	profile, region := parseArguments()
	vPrintln("Loading the AWS profile %s", profile)
	cfg, err := external.LoadDefaultAWSConfig(
		external.WithSharedConfigProfile(profile),
	)
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	if len(region) == 0 {
		region = cfg.Region
	}
	_, err = cfg.Credentials.Retrieve()
	if err != nil {
		panic("unable to retrieve credentials from profile")
	}
	vPrintln("Using the AWS Region %s", region)
	iamc := iam.New(cfg)
	stsc := sts.New(cfg)
	arnChan := mfaDeviceArnChan(stsc, iamc)
	tokenCode := tokenCode()
	mfaDeviceArn := pullMfaDeviceArn(arnChan)
	credentials := sessionCredentials(stsc, mfaDeviceArn, tokenCode, duration(profile))
	storeCredentials(profile, region, credentials)
	storeConfig(profile, region)
	fmt.Printf("Credentials valid until: %s\n", credentials.Expiration)
}

func parseArguments() (string, string) {
	flag.Usage = printUsage
	region := flag.String("r", "", "Override the AWS Region")
	version := flag.Bool("V", false, "Print the skuld version")
	verbose = flag.Bool("v", false, "Print verbose log")
	help := flag.Bool("h", false, "Print this help message")
	flag.Parse()
	tail := flag.Args()
	if *help {
		printUsage()
		os.Exit(0)
	}
	if *version {
		printVersion()
		os.Exit(0)
	}
	if len(tail) != 1 {
		printUsage()
		os.Exit(2)
	}
	return tail[0], *region
}

func printUsage() {
	println("skuld [-r region] <aws profile>")
	printVersion()
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Printf("Skuld version: %s\n", version)
}

type arn struct {
	arn   string
	error interface{}
}

func mfaDeviceArnChan(stsc *sts.STS, iamc *iam.IAM) chan arn {
	result := make(chan arn)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- arn{error: err}
			}
		}()
		deviceArn := mfaDeviceArn(stsc, iamc)
		result <- arn{arn: deviceArn}
	}()
	return result;
}

func mfaDeviceArn(stsc *sts.STS, iamc *iam.IAM) string {
	userArn := userArn(stsc)
	vPrintln("Fetching MFA Device Arn")
	mfaDevice, err := iamc.ListMFADevicesRequest(
		&iam.ListMFADevicesInput{UserName: &userArn},
	).Send()
	if err != nil {
		println(err.Error())
		panic("Unable to fetch the MFA device Arn.")
	}
	return *mfaDevice.MFADevices[0].SerialNumber
}

func userArn(stsc *sts.STS) string {
	vPrintln("Fetching IAM User Arn")
	callerIdResp, err := stsc.GetCallerIdentityRequest(nil).Send()
	if err != nil {
		panic("Unable to get the userArn.")
	}
	arn := callerIdResp.Arn
	return strings.Split(*arn, ":user/")[1]
}

func pullMfaDeviceArn(arn chan arn) string {
	mfaDeviceArn := <-arn
	if mfaDeviceArn.error != nil {
		panic(mfaDeviceArn.error)
	}
	return mfaDeviceArn.arn
}

func duration(profile string) int64 {
	if strings.HasSuffix(profile, "-adm") {
		return admDuration
	}
	return standardDuration
}

func tokenCode() string {
	var tokenCode string
	fmt.Print("Enter your token: ")
	fmt.Scanf("%s", &tokenCode)
	fmt.Println("Fetching temporary credentials...")
	return tokenCode
}

func sessionCredentials(stsc *sts.STS, mfaDevice string, tokenCode string, duration int64) *sts.Credentials {
	token, err := stsc.GetSessionTokenRequest(&sts.GetSessionTokenInput{
		SerialNumber:    &mfaDevice,
		DurationSeconds: &duration,
		TokenCode:       &tokenCode,
	}).Send()
	if err != nil {
		println(err.Error())
		panic("Unable to create a new session.")
	}
	return token.Credentials
}

func storeCredentials(profile string, region string, credentials *sts.Credentials) {
	credsFile := awsFile("credentials")
	creds, err := ini.Load(credsFile)
	if err != nil {
		panic("Unable to load the credential file.")
	}
	tokenProfile := skuldProfile(profile)
	section := creds.Section(tokenProfile)
	section.Key("aws_access_key_id").SetValue(*credentials.AccessKeyId)
	section.Key("aws_secret_access_key").SetValue(*credentials.SecretAccessKey)
	section.Key("aws_session_token").SetValue(*credentials.SessionToken)
	section.Key("region").SetValue(region)
	err = creds.SaveTo(credsFile)
	if err != nil {
		panic("Could not save the credential file.")
	}
}

func awsFile(fileName string) string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err.Error())
	}
	return filepath.Join(dir, ".aws", fileName)
}

func skuldProfile(profile string) string {
	return profile + "-skuld"
}

func storeConfig(profile string, region string) {
	configFile := awsFile("config")
	config, err := ini.Load(configFile)
	if err != nil {
		panic("Unable to load the config file.")
	}
	tokenProfile := skuldProfile(profile)
	section := config.Section(tokenProfile)
	section.Key("region").SetValue(region)
	err = config.SaveTo(configFile)
	if err != nil {
		panic("Unable to save the config file")
	}
}

func vPrintln(format string, a ... interface{}) {
	if *verbose {
		fmt.Printf(format, a...)
		fmt.Println("")
	}
}
