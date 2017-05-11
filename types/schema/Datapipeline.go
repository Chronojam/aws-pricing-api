package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type Datapipeline struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]Datapipeline_Product
	Terms		map[string]map[string]map[string]Datapipeline_Term
}
type Datapipeline_Product struct {	Sku	string
	ProductFamily	string
	Attributes	Datapipeline_Product_Attributes
}
type Datapipeline_Product_Attributes struct {	Operation	string
	ExecutionFrequency	string
	ExecutionLocation	string
	FrequencyMode	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	Usagetype	string
}

type Datapipeline_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]Datapipeline_Term_PriceDimensions
	TermAttributes Datapipeline_Term_TermAttributes
}

type Datapipeline_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	Datapipeline_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type Datapipeline_Term_PricePerUnit struct {
	USD	string
}

type Datapipeline_Term_TermAttributes struct {

}
func (a Datapipeline) QueryProducts(q func(product Datapipeline_Product) bool) []Datapipeline_Product{
	ret := []Datapipeline_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a Datapipeline) QueryTerms(t string, q func(product Datapipeline_Term) bool) []Datapipeline_Term{
	ret := []Datapipeline_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *Datapipeline) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/datapipeline/current/index.json"
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