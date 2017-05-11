package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSDirectConnect struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDirectConnect_Product
	Terms		map[string]map[string]map[string]AWSDirectConnect_Term
}
type AWSDirectConnect_Product struct {	ProductFamily	string
	Attributes	AWSDirectConnect_Product_Attributes
	Sku	string
}
type AWSDirectConnect_Product_Attributes struct {	FromLocationType	string
	Servicecode	string
	FromLocation	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Version	string
	VirtualInterfaceType	string
	TransferType	string
}

type AWSDirectConnect_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSDirectConnect_Term_PriceDimensions
	TermAttributes AWSDirectConnect_Term_TermAttributes
}

type AWSDirectConnect_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDirectConnect_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSDirectConnect_Term_PricePerUnit struct {
	USD	string
}

type AWSDirectConnect_Term_TermAttributes struct {

}
func (a AWSDirectConnect) QueryProducts(q func(product AWSDirectConnect_Product) bool) []AWSDirectConnect_Product{
	ret := []AWSDirectConnect_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSDirectConnect) QueryTerms(t string, q func(product AWSDirectConnect_Term) bool) []AWSDirectConnect_Term{
	ret := []AWSDirectConnect_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AWSDirectConnect) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDirectConnect/current/index.json"
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