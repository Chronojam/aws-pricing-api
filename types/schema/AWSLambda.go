package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSLambda struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSLambda_Product
	Terms		map[string]map[string]map[string]AWSLambda_Term
}
type AWSLambda_Product struct {	ProductFamily	string
	Attributes	AWSLambda_Product_Attributes
	Sku	string
}
type AWSLambda_Product_Attributes struct {	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
}

type AWSLambda_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSLambda_Term_PriceDimensions
	TermAttributes AWSLambda_Term_TermAttributes
}

type AWSLambda_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSLambda_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSLambda_Term_PricePerUnit struct {
	USD	string
}

type AWSLambda_Term_TermAttributes struct {

}
func (a AWSLambda) QueryProducts(q func(product AWSLambda_Product) bool) []AWSLambda_Product{
	ret := []AWSLambda_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSLambda) QueryTerms(t string, q func(product AWSLambda_Term) bool) []AWSLambda_Term{
	ret := []AWSLambda_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AWSLambda) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSLambda/current/index.json"
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