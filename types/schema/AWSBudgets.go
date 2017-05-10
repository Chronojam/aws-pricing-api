package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSBudgets struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSBudgets_Product
	Terms		map[string]map[string]AWSBudgets_Term
}
type AWSBudgets_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSBudgets_Product_Attributes
}
type AWSBudgets_Product_Attributes struct {	Location	string
	LocationType	string
	GroupDescription	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AWSBudgets_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSBudgets_Term_PriceDimensions
	TermAttributes AWSBudgets_Term_TermAttributes
}

type AWSBudgets_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSBudgets_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSBudgets_Term_PricePerUnit struct {
	USD	string
}

type AWSBudgets_Term_TermAttributes struct {

}
func (a *AWSBudgets) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSBudgets/current/index.json"
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