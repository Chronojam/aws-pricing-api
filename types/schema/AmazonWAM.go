package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonWAM struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWAM_Product
	Terms		map[string]map[string]AmazonWAM_Term
}
type AmazonWAM_Product struct {	Attributes	AmazonWAM_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AmazonWAM_Product_Attributes struct {	Usagetype	string
	Operation	string
	PlanType	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
}

type AmazonWAM_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonWAM_Term_PriceDimensions
	TermAttributes AmazonWAM_Term_TermAttributes
}

type AmazonWAM_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonWAM_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonWAM_Term_PricePerUnit struct {
	USD	string
}

type AmazonWAM_Term_TermAttributes struct {

}
func (a *AmazonWAM) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWAM/current/index.json"
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