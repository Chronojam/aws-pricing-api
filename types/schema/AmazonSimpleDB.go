package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonSimpleDB struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonSimpleDB_Product
	Terms		map[string]map[string]map[string]AmazonSimpleDB_Term
}
type AmazonSimpleDB_Product struct {	ProductFamily	string
	Attributes	AmazonSimpleDB_Product_Attributes
	Sku	string
}
type AmazonSimpleDB_Product_Attributes struct {	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
}

type AmazonSimpleDB_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonSimpleDB_Term_PriceDimensions
	TermAttributes AmazonSimpleDB_Term_TermAttributes
}

type AmazonSimpleDB_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonSimpleDB_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonSimpleDB_Term_PricePerUnit struct {
	USD	string
}

type AmazonSimpleDB_Term_TermAttributes struct {

}
func (a AmazonSimpleDB) QueryProducts(q func(product AmazonSimpleDB_Product) bool) []AmazonSimpleDB_Product{
	ret := []AmazonSimpleDB_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonSimpleDB) QueryTerms(t string, q func(product AmazonSimpleDB_Term) bool) []AmazonSimpleDB_Term{
	ret := []AmazonSimpleDB_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonSimpleDB) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonSimpleDB/current/index.json"
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