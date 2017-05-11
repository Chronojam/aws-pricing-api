package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonGlacier struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonGlacier_Product
	Terms		map[string]map[string]map[string]AmazonGlacier_Term
}
type AmazonGlacier_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonGlacier_Product_Attributes
}
type AmazonGlacier_Product_Attributes struct {	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
}

type AmazonGlacier_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonGlacier_Term_PriceDimensions
	TermAttributes AmazonGlacier_Term_TermAttributes
}

type AmazonGlacier_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonGlacier_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonGlacier_Term_PricePerUnit struct {
	USD	string
}

type AmazonGlacier_Term_TermAttributes struct {

}
func (a AmazonGlacier) QueryProducts(q func(product AmazonGlacier_Product) bool) []AmazonGlacier_Product{
	ret := []AmazonGlacier_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonGlacier) QueryTerms(t string, q func(product AmazonGlacier_Term) bool) []AmazonGlacier_Term{
	ret := []AmazonGlacier_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonGlacier) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonGlacier/current/index.json"
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