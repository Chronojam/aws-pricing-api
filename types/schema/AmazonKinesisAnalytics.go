package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonKinesisAnalytics struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonKinesisAnalytics_Product
	Terms		map[string]map[string]map[string]AmazonKinesisAnalytics_Term
}
type AmazonKinesisAnalytics_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonKinesisAnalytics_Product_Attributes
}
type AmazonKinesisAnalytics_Product_Attributes struct {	Servicecode	string
	Description	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
}

type AmazonKinesisAnalytics_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonKinesisAnalytics_Term_PriceDimensions
	TermAttributes AmazonKinesisAnalytics_Term_TermAttributes
}

type AmazonKinesisAnalytics_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonKinesisAnalytics_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonKinesisAnalytics_Term_PricePerUnit struct {
	USD	string
}

type AmazonKinesisAnalytics_Term_TermAttributes struct {

}
func (a AmazonKinesisAnalytics) QueryProducts(q func(product AmazonKinesisAnalytics_Product) bool) []AmazonKinesisAnalytics_Product{
	ret := []AmazonKinesisAnalytics_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonKinesisAnalytics) QueryTerms(t string, q func(product AmazonKinesisAnalytics_Term) bool) []AmazonKinesisAnalytics_Term{
	ret := []AmazonKinesisAnalytics_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonKinesisAnalytics) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesisAnalytics/current/index.json"
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