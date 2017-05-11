package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonWAM struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWAM_Product
	Terms		map[string]map[string]map[string]AmazonWAM_Term
}
type AmazonWAM_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonWAM_Product_Attributes
}
type AmazonWAM_Product_Attributes struct {	PlanType	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	Usagetype	string
	Operation	string
}

type AmazonWAM_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonWAM_Term_PriceDimensions
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
func (a AmazonWAM) QueryProducts(q func(product AmazonWAM_Product) bool) []AmazonWAM_Product{
	ret := []AmazonWAM_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonWAM) QueryTerms(t string, q func(product AmazonWAM_Term) bool) []AmazonWAM_Term{
	ret := []AmazonWAM_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
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