package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonStates struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonStates_Product
	Terms		map[string]map[string]map[string]AmazonStates_Term
}
type AmazonStates_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonStates_Product_Attributes
}
type AmazonStates_Product_Attributes struct {	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
}

type AmazonStates_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonStates_Term_PriceDimensions
	TermAttributes AmazonStates_Term_TermAttributes
}

type AmazonStates_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonStates_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonStates_Term_PricePerUnit struct {
	USD	string
}

type AmazonStates_Term_TermAttributes struct {

}
func (a AmazonStates) QueryProducts(q func(product AmazonStates_Product) bool) []AmazonStates_Product{
	ret := []AmazonStates_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonStates) QueryTerms(t string, q func(product AmazonStates_Term) bool) []AmazonStates_Term{
	ret := []AmazonStates_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonStates) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonStates/current/index.json"
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