package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonCloudSearch struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCloudSearch_Product
	Terms		map[string]map[string]map[string]AmazonCloudSearch_Term
}
type AmazonCloudSearch_Product struct {	Attributes	AmazonCloudSearch_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AmazonCloudSearch_Product_Attributes struct {	Usagetype	string
	Operation	string
	CloudSearchVersion	string
	Servicecode	string
	Location	string
	LocationType	string
	InstanceType	string
}

type AmazonCloudSearch_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonCloudSearch_Term_PriceDimensions
	TermAttributes AmazonCloudSearch_Term_TermAttributes
}

type AmazonCloudSearch_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCloudSearch_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCloudSearch_Term_PricePerUnit struct {
	USD	string
}

type AmazonCloudSearch_Term_TermAttributes struct {

}
func (a AmazonCloudSearch) QueryProducts(q func(product AmazonCloudSearch_Product) bool) []AmazonCloudSearch_Product{
	ret := []AmazonCloudSearch_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonCloudSearch) QueryTerms(t string, q func(product AmazonCloudSearch_Term) bool) []AmazonCloudSearch_Term{
	ret := []AmazonCloudSearch_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonCloudSearch) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudSearch/current/index.json"
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