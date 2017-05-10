package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonStates struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonStates_Product
	Terms		map[string]map[string]AmazonStates_Term
}
type AmazonStates_Product struct {	Attributes	AmazonStates_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AmazonStates_Product_Attributes struct {	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
}

type AmazonStates_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonStates_Term_PriceDimensions
	TermAttributes AmazonStates_Term_TermAttributes
}

type AmazonStates_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonStates_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonStates_Term_PricePerUnit struct {
	USD	string
}

type AmazonStates_Term_TermAttributes struct {

}
func (a *AmazonStates) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonStates/current/index.json"
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