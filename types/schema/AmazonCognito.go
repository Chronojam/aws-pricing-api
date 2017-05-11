package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonCognito struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCognito_Product
	Terms		map[string]map[string]map[string]AmazonCognito_Term
}
type AmazonCognito_Product struct {	Attributes	AmazonCognito_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AmazonCognito_Product_Attributes struct {	Operation	string
	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
}

type AmazonCognito_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonCognito_Term_PriceDimensions
	TermAttributes AmazonCognito_Term_TermAttributes
}

type AmazonCognito_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCognito_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCognito_Term_PricePerUnit struct {
	USD	string
}

type AmazonCognito_Term_TermAttributes struct {

}
func (a AmazonCognito) QueryProducts(q func(product AmazonCognito_Product) bool) []AmazonCognito_Product{
	ret := []AmazonCognito_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonCognito) QueryTerms(t string, q func(product AmazonCognito_Term) bool) []AmazonCognito_Term{
	ret := []AmazonCognito_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonCognito) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCognito/current/index.json"
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