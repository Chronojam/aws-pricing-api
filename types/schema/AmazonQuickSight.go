package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonQuickSight struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonQuickSight_Product
	Terms		map[string]map[string]AmazonQuickSight_Term
}
type AmazonQuickSight_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonQuickSight_Product_Attributes
}
type AmazonQuickSight_Product_Attributes struct {	Servicecode	string
	GroupDescription	string
	Operation	string
	Edition	string
	Location	string
	LocationType	string
	Group	string
	Usagetype	string
	SubscriptionType	string
}

type AmazonQuickSight_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonQuickSight_Term_PriceDimensions
	TermAttributes AmazonQuickSight_Term_TermAttributes
}

type AmazonQuickSight_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonQuickSight_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonQuickSight_Term_PricePerUnit struct {
	USD	string
}

type AmazonQuickSight_Term_TermAttributes struct {

}
func (a *AmazonQuickSight) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonQuickSight/current/index.json"
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