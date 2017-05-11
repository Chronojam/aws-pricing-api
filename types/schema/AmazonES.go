package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonES struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonES_Product
	Terms		map[string]map[string]map[string]AmazonES_Term
}
type AmazonES_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonES_Product_Attributes
}
type AmazonES_Product_Attributes struct {	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
}

type AmazonES_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonES_Term_PriceDimensions
	TermAttributes AmazonES_Term_TermAttributes
}

type AmazonES_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonES_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonES_Term_PricePerUnit struct {
	USD	string
}

type AmazonES_Term_TermAttributes struct {

}
func (a AmazonES) QueryProducts(q func(product AmazonES_Product) bool) []AmazonES_Product{
	ret := []AmazonES_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonES) QueryTerms(t string, q func(product AmazonES_Term) bool) []AmazonES_Term{
	ret := []AmazonES_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonES) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonES/current/index.json"
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