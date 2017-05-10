package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSCodeCommit struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSCodeCommit_Product
	Terms		map[string]map[string]AWSCodeCommit_Term
}
type AWSCodeCommit_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSCodeCommit_Product_Attributes
}
type AWSCodeCommit_Product_Attributes struct {	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	Usagetype	string
	Operation	string
}

type AWSCodeCommit_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSCodeCommit_Term_PriceDimensions
	TermAttributes AWSCodeCommit_Term_TermAttributes
}

type AWSCodeCommit_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSCodeCommit_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSCodeCommit_Term_PricePerUnit struct {
	USD	string
}

type AWSCodeCommit_Term_TermAttributes struct {

}
func (a *AWSCodeCommit) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCodeCommit/current/index.json"
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