package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonCloudDirectory struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCloudDirectory_Product
	Terms		map[string]map[string]AmazonCloudDirectory_Term
}
type AmazonCloudDirectory_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonCloudDirectory_Product_Attributes
}
type AmazonCloudDirectory_Product_Attributes struct {	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
	Servicecode	string
	Location	string
}

type AmazonCloudDirectory_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonCloudDirectory_Term_PriceDimensions
	TermAttributes AmazonCloudDirectory_Term_TermAttributes
}

type AmazonCloudDirectory_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCloudDirectory_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCloudDirectory_Term_PricePerUnit struct {
	USD	string
}

type AmazonCloudDirectory_Term_TermAttributes struct {

}
func (a *AmazonCloudDirectory) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudDirectory/current/index.json"
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