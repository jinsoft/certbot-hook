package providers

import (
	"certbot-hook/utils"
	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
)

type Aliyun struct {
	Client *alidns20150109.Client
	Domain string
	LevelsDomain string
}

func NewAliyun(domain, key, keySecret string) *Aliyun {
	if key == "" {
		panic("accessKey Id must be passed")
	}
	if keySecret == "" {
		panic("accessKey Secret must be passed")
	}
	client, err := createClient(&key, &keySecret)
	if err != nil {
		panic("init aliyun client error")
	}
	domain, levelsDomain := utils.ParseDomain(domain)
	return &Aliyun{
		Client:       client,
		Domain:       domain,
		LevelsDomain: levelsDomain,
	}
}

func createClient(accessKeyId, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dns.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

func (a *Aliyun) ResolveDomainName(dnsType, RR, value string) (recordId *string, _err error) {
	// 先判断有没有解析的
	record, err := a.DescribeDomainRecords(dnsType, RR)
	if err != nil {
		return
	}
	if record.RecordId != "" {
		return a.UpdateDomainRecord(record.RecordId, dnsType, RR, value) // 更新
	} else {
		return a.AddDomainRecord(dnsType, RR, value) // 新增
	}
}

func (a *Aliyun) DeleteResolveDomainName(dnsType, RR string) (recordId *string, _err error) {
	record, _err := a.DescribeDomainRecords(dnsType, RR)
	if _err != nil {
		return
	}
	if record.RecordId == "" {
		return
	}
	return a.DeleteDomainRecord(record.RecordId)
}

// 增加解析记录
func (a *Aliyun) AddDomainRecord(dnsType, RR, value string) (recordId *string, _err error) {
	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(a.Domain),
		RR:         tea.String(RR),
		Type:       tea.String(dnsType),
		Value:      tea.String(value),
	}
	res, _err := a.Client.AddDomainRecord(addDomainRecordRequest)
	if _err != nil {
		return recordId, _err
	}
	return res.Body.RecordId, _err
}

// 修改解析记录
func (a *Aliyun) UpdateDomainRecord(recordId, dnsType, RR, value string) (newRecordId *string, _err error) {
	request := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: tea.String(recordId),
		RR:       tea.String(RR),
		Type:     tea.String(dnsType),
		Value:    tea.String(value),
	}
	resp, _err := a.Client.UpdateDomainRecord(request)
	if _err != nil {
		return
	}
	return resp.Body.RecordId, _err
}

// 根据解析记录id删除解析
func (a *Aliyun) DeleteDomainRecord(recordId string) (newRecordId *string, _err error) {
	request := &alidns20150109.DeleteDomainRecordRequest{
		RecordId: tea.String(recordId),
	}
	resp, _err := a.Client.DeleteDomainRecord(request)
	if _err != nil {
		return
	}
	return resp.Body.RecordId, _err
}

// 获取解析列表
func (a *Aliyun) DescribeDomainRecords(dnsType, RR string ) (record DescribeDomainRecords, _err error) {
	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(a.Domain),
		Type:       tea.String(dnsType),
		PageSize:   tea.Int64(100),
	}
	resp, _err := a.Client.DescribeDomainRecords(describeDomainRecordsRequest)
	if _err != nil {
		return
	}

	if len(resp.Body.DomainRecords.Record) > 0 {
		for _, item := range resp.Body.DomainRecords.Record {
			if RR == *item.RR {
				record.DomainName = *item.DomainName
				record.RecordId = *item.RecordId
				record.Value = *item.Value
				break
			}
		}
	}
	return
}