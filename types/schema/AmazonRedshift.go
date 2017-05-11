package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonRedshift struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonRedshift_Product
	Terms		map[string]map[string]map[string]AmazonRedshift_Term
}
type AmazonRedshift_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonRedshift_Product_Attributes
}
type AmazonRedshift_Product_Attributes struct {	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
}

type AmazonRedshift_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonRedshift_Term_PriceDimensions
	TermAttributes AmazonRedshift_Term_TermAttributes
}

type AmazonRedshift_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonRedshift_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonRedshift_Term_PricePerUnit struct {
	USD	string
}

type AmazonRedshift_Term_TermAttributes struct {

}
func (a AmazonRedshift) QueryProducts(q func(product AmazonRedshift_Product) bool) []AmazonRedshift_Product{
	ret := []AmazonRedshift_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonRedshift) QueryTerms(t string, q func(product AmazonRedshift_Term) bool) []AmazonRedshift_Term{
	ret := []AmazonRedshift_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonRedshift) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRedshift/current/index.json"
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