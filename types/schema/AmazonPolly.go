package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonPolly struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonPolly_Product
	Terms		map[string]map[string]map[string]AmazonPolly_Term
}
type AmazonPolly_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonPolly_Product_Attributes
}
type AmazonPolly_Product_Attributes struct {	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AmazonPolly_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonPolly_Term_PriceDimensions
	TermAttributes AmazonPolly_Term_TermAttributes
}

type AmazonPolly_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonPolly_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonPolly_Term_PricePerUnit struct {
	USD	string
}

type AmazonPolly_Term_TermAttributes struct {

}
func (a AmazonPolly) QueryProducts(q func(product AmazonPolly_Product) bool) []AmazonPolly_Product{
	ret := []AmazonPolly_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonPolly) QueryTerms(t string, q func(product AmazonPolly_Term) bool) []AmazonPolly_Term{
	ret := []AmazonPolly_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonPolly) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonPolly/current/index.json"
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