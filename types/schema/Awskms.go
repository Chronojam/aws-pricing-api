package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type Awskms struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]Awskms_Product
	Terms		map[string]map[string]map[string]Awskms_Term
}
type Awskms_Product struct {	Sku	string
	ProductFamily	string
	Attributes	Awskms_Product_Attributes
}
type Awskms_Product_Attributes struct {	Operation	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
}

type Awskms_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]Awskms_Term_PriceDimensions
	TermAttributes Awskms_Term_TermAttributes
}

type Awskms_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	Awskms_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type Awskms_Term_PricePerUnit struct {
	USD	string
}

type Awskms_Term_TermAttributes struct {

}
func (a Awskms) QueryProducts(q func(product Awskms_Product) bool) []Awskms_Product{
	ret := []Awskms_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a Awskms) QueryTerms(t string, q func(product Awskms_Term) bool) []Awskms_Term{
	ret := []Awskms_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *Awskms) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/awskms/current/index.json"
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