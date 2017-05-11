package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonCloudFront struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCloudFront_Product
	Terms		map[string]map[string]map[string]AmazonCloudFront_Term
}
type AmazonCloudFront_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonCloudFront_Product_Attributes
}
type AmazonCloudFront_Product_Attributes struct {	Servicecode	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
}

type AmazonCloudFront_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonCloudFront_Term_PriceDimensions
	TermAttributes AmazonCloudFront_Term_TermAttributes
}

type AmazonCloudFront_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCloudFront_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCloudFront_Term_PricePerUnit struct {
	USD	string
}

type AmazonCloudFront_Term_TermAttributes struct {

}
func (a AmazonCloudFront) QueryProducts(q func(product AmazonCloudFront_Product) bool) []AmazonCloudFront_Product{
	ret := []AmazonCloudFront_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonCloudFront) QueryTerms(t string, q func(product AmazonCloudFront_Term) bool) []AmazonCloudFront_Term{
	ret := []AmazonCloudFront_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonCloudFront) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudFront/current/index.json"
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