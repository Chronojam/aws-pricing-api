package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonVPC struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonVPC_Product
	Terms		map[string]map[string]map[string]AmazonVPC_Term
}
type AmazonVPC_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonVPC_Product_Attributes
}
type AmazonVPC_Product_Attributes struct {	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
}

type AmazonVPC_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonVPC_Term_PriceDimensions
	TermAttributes AmazonVPC_Term_TermAttributes
}

type AmazonVPC_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonVPC_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonVPC_Term_PricePerUnit struct {
	USD	string
}

type AmazonVPC_Term_TermAttributes struct {

}
func (a AmazonVPC) QueryProducts(q func(product AmazonVPC_Product) bool) []AmazonVPC_Product{
	ret := []AmazonVPC_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonVPC) QueryTerms(t string, q func(product AmazonVPC_Term) bool) []AmazonVPC_Term{
	ret := []AmazonVPC_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonVPC) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonVPC/current/index.json"
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