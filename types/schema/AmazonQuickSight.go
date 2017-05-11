package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonQuickSight struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonQuickSight_Product
	Terms		map[string]map[string]map[string]AmazonQuickSight_Term
}
type AmazonQuickSight_Product struct {	Attributes	AmazonQuickSight_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AmazonQuickSight_Product_Attributes struct {	LocationType	string
	Group	string
	Usagetype	string
	Servicecode	string
	Location	string
	Edition	string
	SubscriptionType	string
	GroupDescription	string
	Operation	string
}

type AmazonQuickSight_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonQuickSight_Term_PriceDimensions
	TermAttributes AmazonQuickSight_Term_TermAttributes
}

type AmazonQuickSight_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonQuickSight_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonQuickSight_Term_PricePerUnit struct {
	USD	string
}

type AmazonQuickSight_Term_TermAttributes struct {

}
func (a AmazonQuickSight) QueryProducts(q func(product AmazonQuickSight_Product) bool) []AmazonQuickSight_Product{
	ret := []AmazonQuickSight_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonQuickSight) QueryTerms(t string, q func(product AmazonQuickSight_Term) bool) []AmazonQuickSight_Term{
	ret := []AmazonQuickSight_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonQuickSight) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonQuickSight/current/index.json"
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