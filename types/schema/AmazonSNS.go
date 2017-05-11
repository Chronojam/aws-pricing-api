package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonSNS struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonSNS_Product
	Terms		map[string]map[string]map[string]AmazonSNS_Term
}
type AmazonSNS_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonSNS_Product_Attributes
}
type AmazonSNS_Product_Attributes struct {	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
}

type AmazonSNS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonSNS_Term_PriceDimensions
	TermAttributes AmazonSNS_Term_TermAttributes
}

type AmazonSNS_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonSNS_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonSNS_Term_PricePerUnit struct {
	USD	string
}

type AmazonSNS_Term_TermAttributes struct {

}
func (a AmazonSNS) QueryProducts(q func(product AmazonSNS_Product) bool) []AmazonSNS_Product{
	ret := []AmazonSNS_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonSNS) QueryTerms(t string, q func(product AmazonSNS_Term) bool) []AmazonSNS_Term{
	ret := []AmazonSNS_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonSNS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonSNS/current/index.json"
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