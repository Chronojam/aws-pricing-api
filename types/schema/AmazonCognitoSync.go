package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonCognitoSync struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCognitoSync_Product
	Terms		map[string]map[string]AmazonCognitoSync_Term
}
type AmazonCognitoSync_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonCognitoSync_Product_Attributes
}
type AmazonCognitoSync_Product_Attributes struct {	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
}

type AmazonCognitoSync_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonCognitoSync_Term_PriceDimensions
	TermAttributes AmazonCognitoSync_Term_TermAttributes
}

type AmazonCognitoSync_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCognitoSync_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonCognitoSync_Term_PricePerUnit struct {
	USD	string
}

type AmazonCognitoSync_Term_TermAttributes struct {

}
func (a *AmazonCognitoSync) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCognitoSync/current/index.json"
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