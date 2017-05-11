package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonWorkSpaces struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWorkSpaces_Product
	Terms		map[string]map[string]map[string]AmazonWorkSpaces_Term
}
type AmazonWorkSpaces_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonWorkSpaces_Product_Attributes
}
type AmazonWorkSpaces_Product_Attributes struct {	License	string
	LocationType	string
	Group	string
	Bundle	string
	Vcpu	string
	Storage	string
	Operation	string
	Location	string
	GroupDescription	string
	ResourceType	string
	RunningMode	string
	SoftwareIncluded	string
	Servicecode	string
	Memory	string
	Usagetype	string
}

type AmazonWorkSpaces_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonWorkSpaces_Term_PriceDimensions
	TermAttributes AmazonWorkSpaces_Term_TermAttributes
}

type AmazonWorkSpaces_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonWorkSpaces_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonWorkSpaces_Term_PricePerUnit struct {
	USD	string
}

type AmazonWorkSpaces_Term_TermAttributes struct {

}
func (a AmazonWorkSpaces) QueryProducts(q func(product AmazonWorkSpaces_Product) bool) []AmazonWorkSpaces_Product{
	ret := []AmazonWorkSpaces_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonWorkSpaces) QueryTerms(t string, q func(product AmazonWorkSpaces_Term) bool) []AmazonWorkSpaces_Term{
	ret := []AmazonWorkSpaces_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonWorkSpaces) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkSpaces/current/index.json"
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