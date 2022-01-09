package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

var server = "dnspod.tencentcloudapi.com"
var recordLine = "默认"
var recordType = "TXT"
var configFile = "/etc/dnspod/config.yml"

var secretId string
var secretKey string
var domain string

type conf struct {
	SecretId  string `yaml:"secretId"`
	SecretKey string `yaml:"secretKey"`
	Domain    string `yaml:"domain"`
}

func (c *conf) getConf(configPath string) *conf {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fatalf("yamlFile.Get err #%v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fatalf("Unmarshal: %v", err)
	}
	return c
}

func fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func extractSubdomain(subDomain string, domain string) string {
	if strings.HasSuffix(subDomain, domain) {
		return subDomain[:len(subDomain)-len(domain)-1]
	}
	return subDomain
}

func newClient() *dnspod.Client {
	credential := common.NewCredential(secretId, secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = server
	client, _ := dnspod.NewClient(credential, "", cpf)
	return client
}

func createTXTRecord(subDomain string, value string) error {

	client := newClient()

	request := dnspod.NewCreateRecordRequest()
	request.Domain = &domain
	request.RecordType = &recordType
	request.RecordLine = &recordLine
	request.SubDomain = &subDomain
	request.Value = &value

	_, err := client.CreateRecord(request)

	if err != nil {
		return err
	}

	fmt.Printf("[CREATE] %s %s %s\n", subDomain, recordType, value)
	return nil
}

func fetchTXTRecords() (map[string]uint64, error) {

	client := newClient()

	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = &domain

	response, err := client.DescribeRecordList(request)
	recordList := response.Response.RecordList

	records := make(map[string]uint64)

	if err != nil {
		return nil, err
	}

	for _, record := range recordList {
		if *record.Type == "TXT" {
			records[*record.Name] = *record.RecordId
		}
	}
	return records, nil
}

func deleteTXTRecord(subDomain string) {
	recordMap, err := fetchTXTRecords()
	if err != nil {
		fatalf("[ERROR] for fetching records: %s\n", err)
	}
	if recordId, ok := recordMap[subDomain]; ok {
		client := newClient()
		request := dnspod.NewDeleteRecordRequest()
		request.Domain = &domain
		request.RecordId = &recordId

		_, err = client.DeleteRecord(request)
		if err != nil {
			fatalf("[ERROR] for deleting %s: %s\n", subDomain, err)
		}
		fmt.Printf("[DELETE] %s\n", subDomain)
	}
}

func main() {
	var dnspodConfig conf
	dnspodConfig.getConf(configFile)
	secretId = dnspodConfig.SecretId
	secretKey = dnspodConfig.SecretKey
	domain = dnspodConfig.Domain

	if len(os.Args) < 2 {
		fatalf("Usage: dnspod-ycli {add|del} ...\n")
	}
	opr := os.Args[1]
	if opr == "del" {
		if len(os.Args) != 3 {
			fatalf("Usage: dnspod-ycli del subdomain\n")
		}
		subDomain := extractSubdomain(os.Args[2], domain)
		deleteTXTRecord(subDomain)
	} else if opr == "add" {
		if len(os.Args) != 4 {
			fatalf("Usage: dnspod-ycli add subdomain value\n")
		}
		subDomain := extractSubdomain(os.Args[2], domain)
		value := os.Args[3]

		createTXTRecord(subDomain, value)
	}
}
