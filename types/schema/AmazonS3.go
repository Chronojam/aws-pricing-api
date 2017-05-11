package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonS3 struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonS3_Product
	Terms		map[string]map[string]map[string]AmazonS3_Term
}
type AmazonS3_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonS3_Product_Attributes
}
type AmazonS3_Product_Attributes struct {	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
}

type AmazonS3_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonS3_Term_PriceDimensions
	TermAttributes AmazonS3_Term_TermAttributes
}

type AmazonS3_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonS3_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonS3_Term_PricePerUnit struct {
	USD	string
}

type AmazonS3_Term_TermAttributes struct {

}
func (a AmazonS3) QueryProducts(q func(product AmazonS3_Product) bool) []AmazonS3_Product{
	ret := []AmazonS3_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonS3) QueryTerms(t string, q func(product AmazonS3_Term) bool) []AmazonS3_Term{
	ret := []AmazonS3_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonS3) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonS3/current/index.json"
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