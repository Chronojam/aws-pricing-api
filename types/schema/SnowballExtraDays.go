package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type SnowballExtraDays struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]SnowballExtraDays_Product
	Terms		map[string]map[string]map[string]SnowballExtraDays_Term
}
type SnowballExtraDays_Product struct {	Sku	string
	ProductFamily	string
	Attributes	SnowballExtraDays_Product_Attributes
}
type SnowballExtraDays_Product_Attributes struct {	Location	string
	LocationType	string
	FeeCode	string
	FeeDescription	string
	Usagetype	string
	Operation	string
	SnowballType	string
	Servicecode	string
}

type SnowballExtraDays_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]SnowballExtraDays_Term_PriceDimensions
	TermAttributes SnowballExtraDays_Term_TermAttributes
}

type SnowballExtraDays_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	SnowballExtraDays_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type SnowballExtraDays_Term_PricePerUnit struct {
	USD	string
}

type SnowballExtraDays_Term_TermAttributes struct {

}
func (a SnowballExtraDays) QueryProducts(q func(product SnowballExtraDays_Product) bool) []SnowballExtraDays_Product{
	ret := []SnowballExtraDays_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a SnowballExtraDays) QueryTerms(t string, q func(product SnowballExtraDays_Term) bool) []SnowballExtraDays_Term{
	ret := []SnowballExtraDays_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *SnowballExtraDays) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/SnowballExtraDays/current/index.json"
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