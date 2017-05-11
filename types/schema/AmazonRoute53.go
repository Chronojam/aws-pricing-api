package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonRoute53 struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonRoute53_Product
	Terms		map[string]map[string]map[string]AmazonRoute53_Term
}
type AmazonRoute53_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonRoute53_Product_Attributes
}
type AmazonRoute53_Product_Attributes struct {	Servicecode	string
	RoutingType	string
	RoutingTarget	string
	Usagetype	string
	Operation	string
}

type AmazonRoute53_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonRoute53_Term_PriceDimensions
	TermAttributes AmazonRoute53_Term_TermAttributes
}

type AmazonRoute53_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonRoute53_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonRoute53_Term_PricePerUnit struct {
	USD	string
}

type AmazonRoute53_Term_TermAttributes struct {

}
func (a AmazonRoute53) QueryProducts(q func(product AmazonRoute53_Product) bool) []AmazonRoute53_Product{
	ret := []AmazonRoute53_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonRoute53) QueryTerms(t string, q func(product AmazonRoute53_Term) bool) []AmazonRoute53_Term{
	ret := []AmazonRoute53_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonRoute53) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRoute53/current/index.json"
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