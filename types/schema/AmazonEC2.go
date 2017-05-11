package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonEC2 struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonEC2_Product
	Terms		map[string]map[string]map[string]AmazonEC2_Term
}
type AmazonEC2_Product struct {	ProductFamily	string
	Attributes	AmazonEC2_Product_Attributes
	Sku	string
}
type AmazonEC2_Product_Attributes struct {	LicenseModel	string
	Usagetype	string
	EnhancedNetworkingSupported	string
	LocationType	string
	PhysicalProcessor	string
	ProcessorArchitecture	string
	ProcessorFeatures	string
	Servicecode	string
	CurrentGeneration	string
	InstanceFamily	string
	Tenancy	string
	Operation	string
	PreInstalledSw	string
	InstanceType	string
	ClockSpeed	string
	Memory	string
	NetworkPerformance	string
	OperatingSystem	string
	Location	string
	Vcpu	string
	Storage	string
}

type AmazonEC2_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonEC2_Term_PriceDimensions
	TermAttributes AmazonEC2_Term_TermAttributes
}

type AmazonEC2_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonEC2_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonEC2_Term_PricePerUnit struct {
	USD	string
}

type AmazonEC2_Term_TermAttributes struct {

}
func (a AmazonEC2) QueryProducts(q func(product AmazonEC2_Product) bool) []AmazonEC2_Product{
	ret := []AmazonEC2_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonEC2) QueryTerms(t string, q func(product AmazonEC2_Term) bool) []AmazonEC2_Term{
	ret := []AmazonEC2_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonEC2) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json"
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