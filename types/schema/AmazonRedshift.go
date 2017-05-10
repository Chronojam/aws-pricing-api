package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonRedshift struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonRedshift_Product
	Terms		map[string]map[string]AmazonRedshift_Term
}
type AmazonRedshift_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonRedshift_Product_Attributes
}
type AmazonRedshift_Product_Attributes struct {	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
}

type AmazonRedshift_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonRedshift_Term_PriceDimensions
	TermAttributes AmazonRedshift_Term_TermAttributes
}

type AmazonRedshift_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonRedshift_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonRedshift_Term_PricePerUnit struct {
	USD	string
}

type AmazonRedshift_Term_TermAttributes struct {

}
func (a *AmazonRedshift) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRedshift/current/index.json"
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