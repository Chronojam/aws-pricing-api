package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSIoT struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSIoT_Product
	Terms		map[string]map[string]map[string]AWSIoT_Term
}
type AWSIoT_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSIoT_Product_Attributes
}
type AWSIoT_Product_Attributes struct {	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	Isshadow	string
	Iswebsocket	string
	Protocol	string
}

type AWSIoT_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSIoT_Term_PriceDimensions
	TermAttributes AWSIoT_Term_TermAttributes
}

type AWSIoT_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSIoT_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSIoT_Term_PricePerUnit struct {
	USD	string
}

type AWSIoT_Term_TermAttributes struct {

}
func (a AWSIoT) QueryProducts(q func(product AWSIoT_Product) bool) []AWSIoT_Product{
	ret := []AWSIoT_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSIoT) QueryTerms(t string, q func(product AWSIoT_Term) bool) []AWSIoT_Term{
	ret := []AWSIoT_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AWSIoT) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSIoT/current/index.json"
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