package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonCognito struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCognito_Product
	Terms		map[string]map[string]AmazonCognito_Term
}
type AmazonCognito_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonCognito_Product_Attributes
}
type AmazonCognito_Product_Attributes struct {	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
}

type AmazonCognito_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonCognito_Term_PriceDimensions
	TermAttributes AmazonCognito_Term_TermAttributes
}

type AmazonCognito_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCognito_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCognito_Term_PricePerUnit struct {
	USD	string
}

type AmazonCognito_Term_TermAttributes struct {

}
func (a *AmazonCognito) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCognito/current/index.json"
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