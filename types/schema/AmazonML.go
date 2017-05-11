package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonML struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonML_Product
	Terms		map[string]map[string]map[string]AmazonML_Term
}
type AmazonML_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonML_Product_Attributes
}
type AmazonML_Product_Attributes struct {	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
	MachineLearningProcess	string
	Servicecode	string
	Location	string
}

type AmazonML_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonML_Term_PriceDimensions
	TermAttributes AmazonML_Term_TermAttributes
}

type AmazonML_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonML_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonML_Term_PricePerUnit struct {
	USD	string
}

type AmazonML_Term_TermAttributes struct {

}
func (a AmazonML) QueryProducts(q func(product AmazonML_Product) bool) []AmazonML_Product{
	ret := []AmazonML_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonML) QueryTerms(t string, q func(product AmazonML_Term) bool) []AmazonML_Term{
	ret := []AmazonML_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonML) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonML/current/index.json"
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