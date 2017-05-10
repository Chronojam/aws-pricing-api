package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonES struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonES_Product
	Terms		map[string]map[string]AmazonES_Term
}
type AmazonES_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonES_Product_Attributes
}
type AmazonES_Product_Attributes struct {	Servicecode	string
	Location	string
	CurrentGeneration	string
	InstanceFamily	string
	Storage	string
	Usagetype	string
	Operation	string
	Ecu	string
	LocationType	string
	InstanceType	string
	Vcpu	string
	MemoryGib	string
}

type AmazonES_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonES_Term_PriceDimensions
	TermAttributes AmazonES_Term_TermAttributes
}

type AmazonES_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonES_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonES_Term_PricePerUnit struct {
	USD	string
}

type AmazonES_Term_TermAttributes struct {

}
func (a *AmazonES) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonES/current/index.json"
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