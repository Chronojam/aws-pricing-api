package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type ElasticMapReduce struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]ElasticMapReduce_Product
	Terms		map[string]map[string]map[string]ElasticMapReduce_Term
}
type ElasticMapReduce_Product struct {	Sku	string
	ProductFamily	string
	Attributes	ElasticMapReduce_Product_Attributes
}
type ElasticMapReduce_Product_Attributes struct {	Location	string
	LocationType	string
	InstanceType	string
	InstanceFamily	string
	Usagetype	string
	Operation	string
	SoftwareType	string
	Servicecode	string
}

type ElasticMapReduce_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]ElasticMapReduce_Term_PriceDimensions
	TermAttributes ElasticMapReduce_Term_TermAttributes
}

type ElasticMapReduce_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	ElasticMapReduce_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type ElasticMapReduce_Term_PricePerUnit struct {
	USD	string
}

type ElasticMapReduce_Term_TermAttributes struct {

}
func (a ElasticMapReduce) QueryProducts(q func(product ElasticMapReduce_Product) bool) []ElasticMapReduce_Product{
	ret := []ElasticMapReduce_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a ElasticMapReduce) QueryTerms(t string, q func(product ElasticMapReduce_Term) bool) []ElasticMapReduce_Term{
	ret := []ElasticMapReduce_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *ElasticMapReduce) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/ElasticMapReduce/current/index.json"
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