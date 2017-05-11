package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonCognitoSync struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCognitoSync_Product
	Terms		map[string]map[string]map[string]AmazonCognitoSync_Term
}
type AmazonCognitoSync_Product struct {	ProductFamily	string
	Attributes	AmazonCognitoSync_Product_Attributes
	Sku	string
}
type AmazonCognitoSync_Product_Attributes struct {	LocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	Location	string
}

type AmazonCognitoSync_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonCognitoSync_Term_PriceDimensions
	TermAttributes AmazonCognitoSync_Term_TermAttributes
}

type AmazonCognitoSync_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCognitoSync_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCognitoSync_Term_PricePerUnit struct {
	USD	string
}

type AmazonCognitoSync_Term_TermAttributes struct {

}
func (a AmazonCognitoSync) QueryProducts(q func(product AmazonCognitoSync_Product) bool) []AmazonCognitoSync_Product{
	ret := []AmazonCognitoSync_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonCognitoSync) QueryTerms(t string, q func(product AmazonCognitoSync_Term) bool) []AmazonCognitoSync_Term{
	ret := []AmazonCognitoSync_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonCognitoSync) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCognitoSync/current/index.json"
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