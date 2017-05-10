package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonPolly struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonPolly_Product
	Terms		map[string]map[string]AmazonPolly_Term
}
type AmazonPolly_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonPolly_Product_Attributes
}
type AmazonPolly_Product_Attributes struct {	Operation	string
	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
}

type AmazonPolly_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonPolly_Term_PriceDimensions
	TermAttributes AmazonPolly_Term_TermAttributes
}

type AmazonPolly_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonPolly_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonPolly_Term_PricePerUnit struct {
	USD	string
}

type AmazonPolly_Term_TermAttributes struct {

}
func (a *AmazonPolly) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonPolly/current/index.json"
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