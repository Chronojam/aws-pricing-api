package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSCloudTrail struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSCloudTrail_Product
	Terms		map[string]map[string]AWSCloudTrail_Term
}
type AWSCloudTrail_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSCloudTrail_Product_Attributes
}
type AWSCloudTrail_Product_Attributes struct {	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AWSCloudTrail_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSCloudTrail_Term_PriceDimensions
	TermAttributes AWSCloudTrail_Term_TermAttributes
}

type AWSCloudTrail_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSCloudTrail_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSCloudTrail_Term_PricePerUnit struct {
	USD	string
}

type AWSCloudTrail_Term_TermAttributes struct {

}
func (a *AWSCloudTrail) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCloudTrail/current/index.json"
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