package service

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/chujieyang/commonops/ops/opslog"
	"github.com/pkg/errors"
)

type aliDnsService struct {
	dns *alidns.Client
}

func GetAliDnsService(region string, accessKey string, keySecret string) (client *aliDnsService, err error) {
	if c, err := alidns.NewClientWithAccessKey(region, accessKey, keySecret); err == nil {
		client = &aliDnsService{
			dns: c,
		}
	}
	return
}

func (client *aliDnsService) GetCloudAccountDomainList(domainName string, pageNumber string, pageSize string) (data string, err error) {
	request := alidns.CreateDescribeDomainsRequest()
	request.Scheme = "https"
	request.KeyWord = domainName
	request.PageNumber = requests.Integer(pageNumber)
	request.PageSize = requests.Integer(pageSize)
	response, err := client.dns.DescribeDomains(request)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	if response != nil && response.IsSuccess() {
		data = response.GetHttpContentString()
	} else {
		err = errors.New(response.GetHttpContentString())
	}
	return
}

func (client *aliDnsService) GetCloudAccountDomainHistoryList(domainName string, pageNumber string, pageSize string) (data string, err error) {
	request := alidns.CreateDescribeRecordLogsRequest()
	request.Scheme = "https"
	request.DomainName = domainName
	request.PageNumber = requests.Integer(pageNumber)
	request.PageSize = requests.Integer(pageSize)
	response, err := client.dns.DescribeRecordLogs(request)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	if response != nil && response.IsSuccess() {
		data = response.GetHttpContentString()
	} else {
		err = errors.New(response.GetHttpContentString())
	}
	return
}

func (client *aliDnsService) GetCloudAccountDomainRecordsList(domainName string, domainRecordRR string,
	pageNumber string, pageSize string) (data string, err error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = domainName
	request.RRKeyWord = domainRecordRR
	request.PageNumber = requests.Integer(pageNumber)
	request.PageSize = requests.Integer(pageSize)
	response, err := client.dns.DescribeDomainRecords(request)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	if response != nil && response.IsSuccess() {
		data = response.GetHttpContentString()
	} else {
		err = errors.New(response.GetHttpContentString())
	}
	return
}

func (client *aliDnsService) AddCloudAccountDomain(domainName string) (data string, err error) {
	request := alidns.CreateAddDomainRequest()
	request.Scheme = "https"
	request.DomainName = domainName
	response, err := client.dns.AddDomain(request)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	if response != nil && response.IsSuccess() {
		data = response.GetHttpContentString()
	} else {
		err = errors.New(response.GetHttpContentString())
	}
	return
}

func (client *aliDnsService) AddCloudAccountDomainRecord(domainName string, rr string, rrType string, value string) (data string, err error) {
	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"
	request.DomainName = domainName
	request.RR = rr
	request.Type = rrType
	request.Value = value
	response, err := client.dns.AddDomainRecord(request)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	if response != nil && response.IsSuccess() {
		data = response.GetHttpContentString()
	} else {
		err = errors.New(response.GetHttpContentString())
	}
	return
}

func (client *aliDnsService) UpdateCloudAccountDomainRecord(recordId string, rr string, rrType string,
	value string) (data string, err error) {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = recordId
	request.RR = rr
	request.Type = rrType
	request.Value = value
	response, err := client.dns.UpdateDomainRecord(request)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	if response != nil && response.IsSuccess() {
		data = response.GetHttpContentString()
	} else {
		err = errors.New(response.GetHttpContentString())
	}
	return
}

func (client *aliDnsService) DeleteCloudAccountDomainRecord(recordId string) (data string, err error) {
	request := alidns.CreateDeleteDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = recordId
	response, err := client.dns.DeleteDomainRecord(request)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	if response != nil && response.IsSuccess() {
		data = response.GetHttpContentString()
	} else {
		err = errors.New(response.GetHttpContentString())
	}
	return
}

func (client *aliDnsService) SetCloudAccountDomainRecordStatus(recordId string, status string) (data string, err error) {
	request := alidns.CreateSetDomainRecordStatusRequest()
	request.Scheme = "https"
	request.RecordId = recordId
	request.Status = status
	response, err := client.dns.SetDomainRecordStatus(request)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	if response != nil && response.IsSuccess() {
		data = response.GetHttpContentString()
	} else {
		err = errors.New(response.GetHttpContentString())
	}
	return
}
