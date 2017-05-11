package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonKinesis struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonKinesis_Product
	Terms		map[string]map[string]map[string]AmazonKinesis_Term
}
type AmazonKinesis_Product struct {	Attributes	AmazonKinesis_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AmazonKinesis_Product_Attributes struct {	LocationType	string
	StandardStorageRetentionIncluded	string
	Operation	string
	MaximumExtendedStorage	string
	Servicecode	string
	Location	string
	Group	string
	GroupDescription	string
	Usagetype	string
}

type AmazonKinesis_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonKinesis_Term_PriceDimensions
	TermAttributes AmazonKinesis_Term_TermAttributes
}

type AmazonKinesis_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonKinesis_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonKinesis_Term_PricePerUnit struct {
	USD	string
}

type AmazonKinesis_Term_TermAttributes struct {

}
func (a AmazonKinesis) QueryProducts(q func(product AmazonKinesis_Product) bool) []AmazonKinesis_Product{
	ret := []AmazonKinesis_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonKinesis) QueryTerms(t string, q func(product AmazonKinesis_Term) bool) []AmazonKinesis_Term{
	ret := []AmazonKinesis_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonKinesis) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesis/current/index.json"
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