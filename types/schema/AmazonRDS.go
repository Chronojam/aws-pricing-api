package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonRDS struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonRDS_Product
	Terms		map[string]map[string]map[string]AmazonRDS_Term
}
type AmazonRDS_Product struct {	Attributes	AmazonRDS_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AmazonRDS_Product_Attributes struct {	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
}

type AmazonRDS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonRDS_Term_PriceDimensions
	TermAttributes AmazonRDS_Term_TermAttributes
}

type AmazonRDS_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonRDS_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonRDS_Term_PricePerUnit struct {
	USD	string
}

type AmazonRDS_Term_TermAttributes struct {

}
func (a AmazonRDS) QueryProducts(q func(product AmazonRDS_Product) bool) []AmazonRDS_Product{
	ret := []AmazonRDS_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonRDS) QueryTerms(t string, q func(product AmazonRDS_Term) bool) []AmazonRDS_Term{
	ret := []AmazonRDS_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonRDS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRDS/current/index.json"
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