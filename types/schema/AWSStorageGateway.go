package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSStorageGateway struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSStorageGateway_Product
	Terms		map[string]map[string]map[string]AWSStorageGateway_Term
}
type AWSStorageGateway_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSStorageGateway_Product_Attributes
}
type AWSStorageGateway_Product_Attributes struct {	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
}

type AWSStorageGateway_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSStorageGateway_Term_PriceDimensions
	TermAttributes AWSStorageGateway_Term_TermAttributes
}

type AWSStorageGateway_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSStorageGateway_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSStorageGateway_Term_PricePerUnit struct {
	USD	string
}

type AWSStorageGateway_Term_TermAttributes struct {

}
func (a AWSStorageGateway) QueryProducts(q func(product AWSStorageGateway_Product) bool) []AWSStorageGateway_Product{
	ret := []AWSStorageGateway_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSStorageGateway) QueryTerms(t string, q func(product AWSStorageGateway_Term) bool) []AWSStorageGateway_Term{
	ret := []AWSStorageGateway_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AWSStorageGateway) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSStorageGateway/current/index.json"
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