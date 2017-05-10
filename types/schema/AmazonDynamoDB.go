package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonDynamoDB struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonDynamoDB_Product
	Terms		map[string]map[string]AmazonDynamoDB_Term
}
type AmazonDynamoDB_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonDynamoDB_Product_Attributes
}
type AmazonDynamoDB_Product_Attributes struct {	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AmazonDynamoDB_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonDynamoDB_Term_PriceDimensions
	TermAttributes AmazonDynamoDB_Term_TermAttributes
}

type AmazonDynamoDB_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonDynamoDB_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonDynamoDB_Term_PricePerUnit struct {
	USD	string
}

type AmazonDynamoDB_Term_TermAttributes struct {

}
func (a *AmazonDynamoDB) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonDynamoDB/current/index.json"
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