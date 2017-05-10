package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSLambda struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSLambda_Product
	Terms		map[string]map[string]AWSLambda_Term
}
type AWSLambda_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSLambda_Product_Attributes
}
type AWSLambda_Product_Attributes struct {	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
}

type AWSLambda_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSLambda_Term_PriceDimensions
	TermAttributes AWSLambda_Term_TermAttributes
}

type AWSLambda_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSLambda_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSLambda_Term_PricePerUnit struct {
	USD	string
}

type AWSLambda_Term_TermAttributes struct {

}
func (a *AWSLambda) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSLambda/current/index.json"
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