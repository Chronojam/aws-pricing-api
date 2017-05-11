package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSServiceCatalog struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSServiceCatalog_Product
	Terms		map[string]map[string]map[string]AWSServiceCatalog_Term
}
type AWSServiceCatalog_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSServiceCatalog_Product_Attributes
}
type AWSServiceCatalog_Product_Attributes struct {	Operation	string
	WithActiveUsers	string
	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
}

type AWSServiceCatalog_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSServiceCatalog_Term_PriceDimensions
	TermAttributes AWSServiceCatalog_Term_TermAttributes
}

type AWSServiceCatalog_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSServiceCatalog_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSServiceCatalog_Term_PricePerUnit struct {
	USD	string
}

type AWSServiceCatalog_Term_TermAttributes struct {

}
func (a AWSServiceCatalog) QueryProducts(q func(product AWSServiceCatalog_Product) bool) []AWSServiceCatalog_Product{
	ret := []AWSServiceCatalog_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSServiceCatalog) QueryTerms(t string, q func(product AWSServiceCatalog_Term) bool) []AWSServiceCatalog_Term{
	ret := []AWSServiceCatalog_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AWSServiceCatalog) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSServiceCatalog/current/index.json"
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