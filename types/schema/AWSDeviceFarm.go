package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSDeviceFarm struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDeviceFarm_Product
	Terms		map[string]map[string]map[string]AWSDeviceFarm_Term
}
type AWSDeviceFarm_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSDeviceFarm_Product_Attributes
}
type AWSDeviceFarm_Product_Attributes struct {	Description	string
	Location	string
	Usagetype	string
	Operation	string
	DeviceOs	string
	Servicecode	string
	LocationType	string
	ExecutionMode	string
	MeterMode	string
}

type AWSDeviceFarm_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSDeviceFarm_Term_PriceDimensions
	TermAttributes AWSDeviceFarm_Term_TermAttributes
}

type AWSDeviceFarm_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDeviceFarm_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSDeviceFarm_Term_PricePerUnit struct {
	USD	string
}

type AWSDeviceFarm_Term_TermAttributes struct {

}
func (a AWSDeviceFarm) QueryProducts(q func(product AWSDeviceFarm_Product) bool) []AWSDeviceFarm_Product{
	ret := []AWSDeviceFarm_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSDeviceFarm) QueryTerms(t string, q func(product AWSDeviceFarm_Term) bool) []AWSDeviceFarm_Term{
	ret := []AWSDeviceFarm_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AWSDeviceFarm) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDeviceFarm/current/index.json"
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