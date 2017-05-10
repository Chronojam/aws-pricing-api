package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonAthena struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonAthena_Product
	Terms		map[string]map[string]AmazonAthena_Term
}
type AmazonAthena_Product struct {	ProductFamily	string
	Attributes	AmazonAthena_Product_Attributes
	Sku	string
}
type AmazonAthena_Product_Attributes struct {	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	FreeQueryTypes	string
	Servicecode	string
	Description	string
}

type AmazonAthena_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonAthena_Term_PriceDimensions
	TermAttributes AmazonAthena_Term_TermAttributes
}

type AmazonAthena_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonAthena_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonAthena_Term_PricePerUnit struct {
	USD	string
}

type AmazonAthena_Term_TermAttributes struct {

}
func (a *AmazonAthena) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonAthena/current/index.json"
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