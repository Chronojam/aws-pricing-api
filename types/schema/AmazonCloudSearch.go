package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonCloudSearch struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCloudSearch_Product
	Terms		map[string]map[string]AmazonCloudSearch_Term
}
type AmazonCloudSearch_Product struct {	ProductFamily	string
	Attributes	AmazonCloudSearch_Product_Attributes
	Sku	string
}
type AmazonCloudSearch_Product_Attributes struct {	LocationType	string
	InstanceType	string
	Usagetype	string
	Operation	string
	CloudSearchVersion	string
	Servicecode	string
	Location	string
}

type AmazonCloudSearch_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonCloudSearch_Term_PriceDimensions
	TermAttributes AmazonCloudSearch_Term_TermAttributes
}

type AmazonCloudSearch_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCloudSearch_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCloudSearch_Term_PricePerUnit struct {
	USD	string
}

type AmazonCloudSearch_Term_TermAttributes struct {

}
func (a *AmazonCloudSearch) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudSearch/current/index.json"
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