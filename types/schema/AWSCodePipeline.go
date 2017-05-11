package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSCodePipeline struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSCodePipeline_Product
	Terms		map[string]map[string]map[string]AWSCodePipeline_Term
}
type AWSCodePipeline_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSCodePipeline_Product_Attributes
}
type AWSCodePipeline_Product_Attributes struct {	Description	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AWSCodePipeline_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSCodePipeline_Term_PriceDimensions
	TermAttributes AWSCodePipeline_Term_TermAttributes
}

type AWSCodePipeline_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSCodePipeline_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSCodePipeline_Term_PricePerUnit struct {
	USD	string
}

type AWSCodePipeline_Term_TermAttributes struct {

}
func (a AWSCodePipeline) QueryProducts(q func(product AWSCodePipeline_Product) bool) []AWSCodePipeline_Product{
	ret := []AWSCodePipeline_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSCodePipeline) QueryTerms(t string, q func(product AWSCodePipeline_Term) bool) []AWSCodePipeline_Term{
	ret := []AWSCodePipeline_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AWSCodePipeline) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCodePipeline/current/index.json"
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