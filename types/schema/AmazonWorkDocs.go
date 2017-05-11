package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonWorkDocs struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWorkDocs_Product
	Terms		map[string]map[string]map[string]AmazonWorkDocs_Term
}
type AmazonWorkDocs_Product struct {	ProductFamily	string
	Attributes	AmazonWorkDocs_Product_Attributes
	Sku	string
}
type AmazonWorkDocs_Product_Attributes struct {	Description	string
	Location	string
	LocationType	string
	Storage	string
	Usagetype	string
	Operation	string
	Servicecode	string
	FreeTrial	string
	MaximumStorageVolume	string
	MinimumStorageVolume	string
}

type AmazonWorkDocs_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonWorkDocs_Term_PriceDimensions
	TermAttributes AmazonWorkDocs_Term_TermAttributes
}

type AmazonWorkDocs_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonWorkDocs_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonWorkDocs_Term_PricePerUnit struct {
	USD	string
}

type AmazonWorkDocs_Term_TermAttributes struct {

}
func (a AmazonWorkDocs) QueryProducts(q func(product AmazonWorkDocs_Product) bool) []AmazonWorkDocs_Product{
	ret := []AmazonWorkDocs_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonWorkDocs) QueryTerms(t string, q func(product AmazonWorkDocs_Term) bool) []AmazonWorkDocs_Term{
	ret := []AmazonWorkDocs_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonWorkDocs) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkDocs/current/index.json"
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