package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonWorkSpaces struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWorkSpaces_Product
	Terms		map[string]map[string]AmazonWorkSpaces_Term
}
type AmazonWorkSpaces_Product struct {	ProductFamily	string
	Attributes	AmazonWorkSpaces_Product_Attributes
	Sku	string
}
type AmazonWorkSpaces_Product_Attributes struct {	Storage	string
	Usagetype	string
	ResourceType	string
	RunningMode	string
	GroupDescription	string
	Bundle	string
	License	string
	Location	string
	Vcpu	string
	Memory	string
	Group	string
	Servicecode	string
	LocationType	string
	Operation	string
	SoftwareIncluded	string
}

type AmazonWorkSpaces_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonWorkSpaces_Term_PriceDimensions
	TermAttributes AmazonWorkSpaces_Term_TermAttributes
}

type AmazonWorkSpaces_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonWorkSpaces_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonWorkSpaces_Term_PricePerUnit struct {
	USD	string
}

type AmazonWorkSpaces_Term_TermAttributes struct {

}
func (a *AmazonWorkSpaces) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkSpaces/current/index.json"
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