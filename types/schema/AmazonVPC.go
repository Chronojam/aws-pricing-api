package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonVPC struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonVPC_Product
	Terms		map[string]map[string]AmazonVPC_Term
}
type AmazonVPC_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonVPC_Product_Attributes
}
type AmazonVPC_Product_Attributes struct {	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
}

type AmazonVPC_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonVPC_Term_PriceDimensions
	TermAttributes AmazonVPC_Term_TermAttributes
}

type AmazonVPC_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonVPC_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonVPC_Term_PricePerUnit struct {
	USD	string
}

type AmazonVPC_Term_TermAttributes struct {

}
func (a *AmazonVPC) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonVPC/current/index.json"
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