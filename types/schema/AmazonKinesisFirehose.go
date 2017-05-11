package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonKinesisFirehose struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonKinesisFirehose_Product
	Terms		map[string]map[string]map[string]AmazonKinesisFirehose_Term
}
type AmazonKinesisFirehose_Product struct {	Attributes	AmazonKinesisFirehose_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AmazonKinesisFirehose_Product_Attributes struct {	Group	string
	Usagetype	string
	Operation	string
	Servicecode	string
	Description	string
	Location	string
	LocationType	string
}

type AmazonKinesisFirehose_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonKinesisFirehose_Term_PriceDimensions
	TermAttributes AmazonKinesisFirehose_Term_TermAttributes
}

type AmazonKinesisFirehose_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonKinesisFirehose_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonKinesisFirehose_Term_PricePerUnit struct {
	USD	string
}

type AmazonKinesisFirehose_Term_TermAttributes struct {

}
func (a AmazonKinesisFirehose) QueryProducts(q func(product AmazonKinesisFirehose_Product) bool) []AmazonKinesisFirehose_Product{
	ret := []AmazonKinesisFirehose_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonKinesisFirehose) QueryTerms(t string, q func(product AmazonKinesisFirehose_Term) bool) []AmazonKinesisFirehose_Term{
	ret := []AmazonKinesisFirehose_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonKinesisFirehose) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesisFirehose/current/index.json"
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