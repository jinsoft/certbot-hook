package providers

type AddDomainRecordResponseBody struct {
	RequestId string
	RecordId  string
}

type DescribeDomainRecords struct {
	RecordId   string
	DomainName string
	Value      string
}