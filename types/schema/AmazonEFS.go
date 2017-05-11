package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonEFS struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonEFS_Product
	Terms		map[string]map[string]map[string]AmazonEFS_Term
}
type AmazonEFS_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonEFS_Product_Attributes
}
type AmazonEFS_Product_Attributes struct {	Location	string
	LocationType	string
	StorageClass	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AmazonEFS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonEFS_Term_PriceDimensions
	TermAttributes AmazonEFS_Term_TermAttributes
}

type AmazonEFS_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonEFS_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonEFS_Term_PricePerUnit struct {
	USD	string
}

type AmazonEFS_Term_TermAttributes struct {

}
func (a AmazonEFS) QueryProducts(q func(product AmazonEFS_Product) bool) []AmazonEFS_Product{
	ret := []AmazonEFS_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonEFS) QueryTerms(t string, q func(product AmazonEFS_Term) bool) []AmazonEFS_Term{
	ret := []AmazonEFS_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonEFS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEFS/current/index.json"
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