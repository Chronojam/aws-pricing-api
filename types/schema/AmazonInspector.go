package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonInspector struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonInspector_Product
	Terms		map[string]map[string]map[string]AmazonInspector_Term
}
type AmazonInspector_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonInspector_Product_Attributes
}
type AmazonInspector_Product_Attributes struct {	Usagetype	string
	Operation	string
	FreeUsageIncluded	string
	Servicecode	string
	Description	string
	Location	string
	LocationType	string
}

type AmazonInspector_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonInspector_Term_PriceDimensions
	TermAttributes AmazonInspector_Term_TermAttributes
}

type AmazonInspector_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonInspector_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonInspector_Term_PricePerUnit struct {
	USD	string
}

type AmazonInspector_Term_TermAttributes struct {

}
func (a AmazonInspector) QueryProducts(q func(product AmazonInspector_Product) bool) []AmazonInspector_Product{
	ret := []AmazonInspector_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonInspector) QueryTerms(t string, q func(product AmazonInspector_Term) bool) []AmazonInspector_Term{
	ret := []AmazonInspector_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonInspector) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonInspector/current/index.json"
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