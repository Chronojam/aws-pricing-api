package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonCloudFront struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCloudFront_Product
	Terms		map[string]map[string]AmazonCloudFront_Term
}
type AmazonCloudFront_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonCloudFront_Product_Attributes
}
type AmazonCloudFront_Product_Attributes struct {	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AmazonCloudFront_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonCloudFront_Term_PriceDimensions
	TermAttributes AmazonCloudFront_Term_TermAttributes
}

type AmazonCloudFront_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCloudFront_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCloudFront_Term_PricePerUnit struct {
	USD	string
}

type AmazonCloudFront_Term_TermAttributes struct {

}
func (a *AmazonCloudFront) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudFront/current/index.json"
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