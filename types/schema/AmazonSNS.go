package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonSNS struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonSNS_Product
	Terms		map[string]map[string]AmazonSNS_Term
}
type AmazonSNS_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonSNS_Product_Attributes
}
type AmazonSNS_Product_Attributes struct {	Usagetype	string
	Operation	string
	Servicecode	string
	Description	string
	Location	string
	LocationType	string
	EndpointType	string
}

type AmazonSNS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonSNS_Term_PriceDimensions
	TermAttributes AmazonSNS_Term_TermAttributes
}

type AmazonSNS_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonSNS_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonSNS_Term_PricePerUnit struct {
	USD	string
}

type AmazonSNS_Term_TermAttributes struct {

}
func (a *AmazonSNS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonSNS/current/index.json"
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