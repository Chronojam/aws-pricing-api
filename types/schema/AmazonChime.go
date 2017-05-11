package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonChime struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonChime_Product
	Terms		map[string]map[string]map[string]AmazonChime_Term
}
type AmazonChime_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonChime_Product_Attributes
}
type AmazonChime_Product_Attributes struct {	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	LicenseType	string
	Servicecode	string
}

type AmazonChime_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonChime_Term_PriceDimensions
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
func (a AmazonChime) QueryProducts(q func(product AmazonChime_Product) bool) []AmazonChime_Product{
	ret := []AmazonChime_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonChime) QueryTerms(t string, q func(product AmazonChime_Term) bool) []AmazonChime_Term{
	ret := []AmazonChime_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
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