package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonApiGateway struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonApiGateway_Product
	Terms		map[string]map[string]map[string]AmazonApiGateway_Term
}
type AmazonApiGateway_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonApiGateway_Product_Attributes
}
type AmazonApiGateway_Product_Attributes struct {	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
}

type AmazonApiGateway_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonApiGateway_Term_PriceDimensions
	TermAttributes AmazonApiGateway_Term_TermAttributes
}

type AmazonApiGateway_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonApiGateway_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonApiGateway_Term_PricePerUnit struct {
	USD	string
}

type AmazonApiGateway_Term_TermAttributes struct {

}
func (a AmazonApiGateway) QueryProducts(q func(product AmazonApiGateway_Product) bool) []AmazonApiGateway_Product{
	ret := []AmazonApiGateway_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonApiGateway) QueryTerms(t string, q func(product AmazonApiGateway_Term) bool) []AmazonApiGateway_Term{
	ret := []AmazonApiGateway_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonApiGateway) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonApiGateway/current/index.json"
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