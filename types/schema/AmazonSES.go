package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonSES struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonSES_Product
	Terms		map[string]map[string]AmazonSES_Term
}
type AmazonSES_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonSES_Product_Attributes
}
type AmazonSES_Product_Attributes struct {	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
}

type AmazonSES_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonSES_Term_PriceDimensions
	TermAttributes AmazonSES_Term_TermAttributes
}

type AmazonSES_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonSES_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonSES_Term_PricePerUnit struct {
	USD	string
}

type AmazonSES_Term_TermAttributes struct {

}
func (a *AmazonSES) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonSES/current/index.json"
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