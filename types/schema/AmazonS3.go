package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonS3 struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonS3_Product
	Terms		map[string]map[string]AmazonS3_Term
}
type AmazonS3_Product struct {	ProductFamily	string
	Attributes	AmazonS3_Product_Attributes
	Sku	string
}
type AmazonS3_Product_Attributes struct {	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
}

type AmazonS3_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonS3_Term_PriceDimensions
	TermAttributes AmazonS3_Term_TermAttributes
}

type AmazonS3_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonS3_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonS3_Term_PricePerUnit struct {
	USD	string
}

type AmazonS3_Term_TermAttributes struct {

}
func (a *AmazonS3) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonS3/current/index.json"
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