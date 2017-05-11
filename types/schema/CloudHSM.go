package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type CloudHSM struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]CloudHSM_Product
	Terms		map[string]map[string]map[string]CloudHSM_Term
}
type CloudHSM_Product struct {	Sku	string
	ProductFamily	string
	Attributes	CloudHSM_Product_Attributes
}
type CloudHSM_Product_Attributes struct {	Location	string
	LocationType	string
	InstanceFamily	string
	Usagetype	string
	Operation	string
	TrialProduct	string
	UpfrontCommitment	string
	Servicecode	string
}

type CloudHSM_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]CloudHSM_Term_PriceDimensions
	TermAttributes CloudHSM_Term_TermAttributes
}

type CloudHSM_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	CloudHSM_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type CloudHSM_Term_PricePerUnit struct {
	USD	string
}

type CloudHSM_Term_TermAttributes struct {

}
func (a CloudHSM) QueryProducts(q func(product CloudHSM_Product) bool) []CloudHSM_Product{
	ret := []CloudHSM_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a CloudHSM) QueryTerms(t string, q func(product CloudHSM_Term) bool) []CloudHSM_Term{
	ret := []CloudHSM_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *CloudHSM) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/CloudHSM/current/index.json"
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