package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type IngestionServiceSnowball struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]IngestionServiceSnowball_Product
	Terms		map[string]map[string]map[string]IngestionServiceSnowball_Term
}
type IngestionServiceSnowball_Product struct {	Sku	string
	ProductFamily	string
	Attributes	IngestionServiceSnowball_Product_Attributes
}
type IngestionServiceSnowball_Product_Attributes struct {	SnowballType	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
}

type IngestionServiceSnowball_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]IngestionServiceSnowball_Term_PriceDimensions
	TermAttributes IngestionServiceSnowball_Term_TermAttributes
}

type IngestionServiceSnowball_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	IngestionServiceSnowball_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type IngestionServiceSnowball_Term_PricePerUnit struct {
	USD	string
}

type IngestionServiceSnowball_Term_TermAttributes struct {

}
func (a IngestionServiceSnowball) QueryProducts(q func(product IngestionServiceSnowball_Product) bool) []IngestionServiceSnowball_Product{
	ret := []IngestionServiceSnowball_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a IngestionServiceSnowball) QueryTerms(t string, q func(product IngestionServiceSnowball_Term) bool) []IngestionServiceSnowball_Term{
	ret := []IngestionServiceSnowball_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *IngestionServiceSnowball) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/IngestionServiceSnowball/current/index.json"
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