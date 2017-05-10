package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSQueueService struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSQueueService_Product
	Terms		map[string]map[string]AWSQueueService_Term
}
type AWSQueueService_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSQueueService_Product_Attributes
}
type AWSQueueService_Product_Attributes struct {	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
}

type AWSQueueService_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSQueueService_Term_PriceDimensions
	TermAttributes AWSQueueService_Term_TermAttributes
}

type AWSQueueService_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSQueueService_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSQueueService_Term_PricePerUnit struct {
	USD	string
}

type AWSQueueService_Term_TermAttributes struct {

}
func (a *AWSQueueService) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSQueueService/current/index.json"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, a)
	if err != nil {
		return err
	}

	return nil
}