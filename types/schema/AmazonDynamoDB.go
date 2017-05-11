package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonDynamoDB struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonDynamoDB_Product
	Terms		map[string]map[string]map[string]AmazonDynamoDB_Term
}
type AmazonDynamoDB_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonDynamoDB_Product_Attributes
}
type AmazonDynamoDB_Product_Attributes struct {	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AmazonDynamoDB_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonDynamoDB_Term_PriceDimensions
	TermAttributes AmazonDynamoDB_Term_TermAttributes
}

type AmazonDynamoDB_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonDynamoDB_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonDynamoDB_Term_PricePerUnit struct {
	USD	string
}

type AmazonDynamoDB_Term_TermAttributes struct {

}
func (a AmazonDynamoDB) QueryProducts(q func(product AmazonDynamoDB_Product) bool) []AmazonDynamoDB_Product{
	ret := []AmazonDynamoDB_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonDynamoDB) QueryTerms(t string, q func(product AmazonDynamoDB_Term) bool) []AmazonDynamoDB_Term{
	ret := []AmazonDynamoDB_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonDynamoDB) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonDynamoDB/current/index.json"
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