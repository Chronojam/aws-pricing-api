package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSCodeDeploy struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSCodeDeploy_Product
	Terms		map[string]map[string]map[string]AWSCodeDeploy_Term
}
type AWSCodeDeploy_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSCodeDeploy_Product_Attributes
}
type AWSCodeDeploy_Product_Attributes struct {	LocationType	string
	Usagetype	string
	Operation	string
	DeploymentLocation	string
	Servicecode	string
	Location	string
}

type AWSCodeDeploy_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSCodeDeploy_Term_PriceDimensions
	TermAttributes AWSCodeDeploy_Term_TermAttributes
}

type AWSCodeDeploy_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSCodeDeploy_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSCodeDeploy_Term_PricePerUnit struct {
	USD	string
}

type AWSCodeDeploy_Term_TermAttributes struct {

}
func (a AWSCodeDeploy) QueryProducts(q func(product AWSCodeDeploy_Product) bool) []AWSCodeDeploy_Product{
	ret := []AWSCodeDeploy_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSCodeDeploy) QueryTerms(t string, q func(product AWSCodeDeploy_Term) bool) []AWSCodeDeploy_Term{
	ret := []AWSCodeDeploy_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AWSCodeDeploy) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCodeDeploy/current/index.json"
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