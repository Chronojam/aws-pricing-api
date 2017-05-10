package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonECR struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonECR_Product
	Terms		map[string]map[string]AmazonECR_Term
}
type AmazonECR_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonECR_Product_Attributes
}
type AmazonECR_Product_Attributes struct {	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
}

type AmazonECR_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonECR_Term_PriceDimensions
	TermAttributes AmazonECR_Term_TermAttributes
}

type AmazonECR_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonECR_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonECR_Term_PricePerUnit struct {
	USD	string
}

type AmazonECR_Term_TermAttributes struct {

}
func (a *AmazonECR) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonECR/current/index.json"
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