package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonChime struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonChime_Product
	Terms		map[string]map[string]AmazonChime_Term
}
type AmazonChime_Product struct {	Attributes	AmazonChime_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AmazonChime_Product_Attributes struct {	Usagetype	string
	Operation	string
	LicenseType	string
	Servicecode	string
	Location	string
	LocationType	string
}

type AmazonChime_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonChime_Term_PriceDimensions
	TermAttributes AmazonChime_Term_TermAttributes
}

type AmazonChime_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonChime_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonChime_Term_PricePerUnit struct {
	USD	string
}

type AmazonChime_Term_TermAttributes struct {

}
func (a *AmazonChime) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonChime/current/index.json"
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