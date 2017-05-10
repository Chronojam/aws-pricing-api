package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonApiGateway struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonApiGateway_Product
	Terms		map[string]map[string]AmazonApiGateway_Term
}
type AmazonApiGateway_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonApiGateway_Product_Attributes
}
type AmazonApiGateway_Product_Attributes struct {	Description	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AmazonApiGateway_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonApiGateway_Term_PriceDimensions
	TermAttributes AmazonApiGateway_Term_TermAttributes
}

type AmazonApiGateway_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonApiGateway_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonApiGateway_Term_PricePerUnit struct {
	USD	string
}

type AmazonApiGateway_Term_TermAttributes struct {

}
func (a *AmazonApiGateway) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonApiGateway/current/index.json"
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