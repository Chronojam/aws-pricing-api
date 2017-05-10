package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonML struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonML_Product
	Terms		map[string]map[string]AmazonML_Term
}
type AmazonML_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonML_Product_Attributes
}
type AmazonML_Product_Attributes struct {	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
	MachineLearningProcess	string
}

type AmazonML_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonML_Term_PriceDimensions
	TermAttributes AmazonML_Term_TermAttributes
}

type AmazonML_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonML_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonML_Term_PricePerUnit struct {
	USD	string
}

type AmazonML_Term_TermAttributes struct {

}
func (a *AmazonML) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonML/current/index.json"
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