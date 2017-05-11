package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type IngestionService struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]IngestionService_Product
	Terms		map[string]map[string]map[string]IngestionService_Term
}
type IngestionService_Product struct {	Sku	string
	ProductFamily	string
	Attributes	IngestionService_Product_Attributes
}
type IngestionService_Product_Attributes struct {	DataAction	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
}

type IngestionService_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]IngestionService_Term_PriceDimensions
	TermAttributes IngestionService_Term_TermAttributes
}

type IngestionService_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	IngestionService_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type IngestionService_Term_PricePerUnit struct {
	USD	string
}

type IngestionService_Term_TermAttributes struct {

}
func (a IngestionService) QueryProducts(q func(product IngestionService_Product) bool) []IngestionService_Product{
	ret := []IngestionService_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a IngestionService) QueryTerms(t string, q func(product IngestionService_Term) bool) []IngestionService_Term{
	ret := []IngestionService_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *IngestionService) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/IngestionService/current/index.json"
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