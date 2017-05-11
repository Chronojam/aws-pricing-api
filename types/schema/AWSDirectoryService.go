package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSDirectoryService struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDirectoryService_Product
	Terms		map[string]map[string]map[string]AWSDirectoryService_Term
}
type AWSDirectoryService_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSDirectoryService_Product_Attributes
}
type AWSDirectoryService_Product_Attributes struct {	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	DirectorySize	string
	DirectoryType	string
	DirectoryTypeDescription	string
	Servicecode	string
}

type AWSDirectoryService_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSDirectoryService_Term_PriceDimensions
	TermAttributes AWSDirectoryService_Term_TermAttributes
}

type AWSDirectoryService_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDirectoryService_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSDirectoryService_Term_PricePerUnit struct {
	USD	string
}

type AWSDirectoryService_Term_TermAttributes struct {

}
func (a AWSDirectoryService) QueryProducts(q func(product AWSDirectoryService_Product) bool) []AWSDirectoryService_Product{
	ret := []AWSDirectoryService_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSDirectoryService) QueryTerms(t string, q func(product AWSDirectoryService_Term) bool) []AWSDirectoryService_Term{
	ret := []AWSDirectoryService_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AWSDirectoryService) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDirectoryService/current/index.json"
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