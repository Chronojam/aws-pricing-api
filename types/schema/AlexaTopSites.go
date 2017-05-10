package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AlexaTopSites struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AlexaTopSites_Product
	Terms		map[string]map[string]AlexaTopSites_Term
}
type AlexaTopSites_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AlexaTopSites_Product_Attributes
}
type AlexaTopSites_Product_Attributes struct {	Operation	string
	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
}

type AlexaTopSites_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AlexaTopSites_Term_PriceDimensions
	TermAttributes AlexaTopSites_Term_TermAttributes
}

type AlexaTopSites_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AlexaTopSites_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AlexaTopSites_Term_PricePerUnit struct {
	USD	string
}

type AlexaTopSites_Term_TermAttributes struct {

}
func (a *AlexaTopSites) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AlexaTopSites/current/index.json"
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