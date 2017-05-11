package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonSWF struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonSWF_Product
	Terms		map[string]map[string]map[string]AmazonSWF_Term
}
type AmazonSWF_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonSWF_Product_Attributes
}
type AmazonSWF_Product_Attributes struct {	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AmazonSWF_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonSWF_Term_PriceDimensions
	TermAttributes AmazonSWF_Term_TermAttributes
}

type AmazonSWF_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonSWF_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonSWF_Term_PricePerUnit struct {
	USD	string
}

type AmazonSWF_Term_TermAttributes struct {

}
func (a AmazonSWF) QueryProducts(q func(product AmazonSWF_Product) bool) []AmazonSWF_Product{
	ret := []AmazonSWF_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonSWF) QueryTerms(t string, q func(product AmazonSWF_Term) bool) []AmazonSWF_Term{
	ret := []AmazonSWF_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonSWF) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonSWF/current/index.json"
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