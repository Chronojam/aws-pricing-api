package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonElastiCache struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonElastiCache_Product
	Terms		map[string]map[string]map[string]AmazonElastiCache_Term
}
type AmazonElastiCache_Product struct {	ProductFamily	string
	Attributes	AmazonElastiCache_Product_Attributes
	Sku	string
}
type AmazonElastiCache_Product_Attributes struct {	Location	string
	InstanceType	string
	Vcpu	string
	CacheEngine	string
	Servicecode	string
	CurrentGeneration	string
	InstanceFamily	string
	Memory	string
	NetworkPerformance	string
	Usagetype	string
	Operation	string
	LocationType	string
}

type AmazonElastiCache_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonElastiCache_Term_PriceDimensions
	TermAttributes AmazonElastiCache_Term_TermAttributes
}

type AmazonElastiCache_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonElastiCache_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonElastiCache_Term_PricePerUnit struct {
	USD	string
}

type AmazonElastiCache_Term_TermAttributes struct {

}
func (a AmazonElastiCache) QueryProducts(q func(product AmazonElastiCache_Product) bool) []AmazonElastiCache_Product{
	ret := []AmazonElastiCache_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonElastiCache) QueryTerms(t string, q func(product AmazonElastiCache_Term) bool) []AmazonElastiCache_Term{
	ret := []AmazonElastiCache_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonElastiCache) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonElastiCache/current/index.json"
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