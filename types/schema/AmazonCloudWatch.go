package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonCloudWatch struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCloudWatch_Product
	Terms		map[string]map[string]AmazonCloudWatch_Term
}
type AmazonCloudWatch_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonCloudWatch_Product_Attributes
}
type AmazonCloudWatch_Product_Attributes struct {	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
}

type AmazonCloudWatch_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonCloudWatch_Term_PriceDimensions
	TermAttributes AmazonCloudWatch_Term_TermAttributes
}

type AmazonCloudWatch_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCloudWatch_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCloudWatch_Term_PricePerUnit struct {
	USD	string
}

type AmazonCloudWatch_Term_TermAttributes struct {

}
func (a *AmazonCloudWatch) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudWatch/current/index.json"
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