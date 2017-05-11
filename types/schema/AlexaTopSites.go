package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AlexaTopSites struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AlexaTopSites_Product
	Terms		map[string]map[string]map[string]AlexaTopSites_Term
}
type AlexaTopSites_Product struct {	ProductFamily	string
	Attributes	AlexaTopSites_Product_Attributes
	Sku	string
}
type AlexaTopSites_Product_Attributes struct {	LocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	Location	string
}

type AlexaTopSites_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AlexaTopSites_Term_PriceDimensions
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
func (a AlexaTopSites) QueryProducts(q func(product AlexaTopSites_Product) bool) []AlexaTopSites_Product{
	ret := []AlexaTopSites_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AlexaTopSites) QueryTerms(t string, q func(product AlexaTopSites_Term) bool) []AlexaTopSites_Term{
	ret := []AlexaTopSites_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
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