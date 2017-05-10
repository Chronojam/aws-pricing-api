package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonGlacier struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonGlacier_Product
	Terms		map[string]map[string]AmazonGlacier_Term
}
type AmazonGlacier_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonGlacier_Product_Attributes
}
type AmazonGlacier_Product_Attributes struct {	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
}

type AmazonGlacier_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonGlacier_Term_PriceDimensions
	TermAttributes AmazonGlacier_Term_TermAttributes
}

type AmazonGlacier_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonGlacier_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonGlacier_Term_PricePerUnit struct {
	USD	string
}

type AmazonGlacier_Term_TermAttributes struct {

}
func (a *AmazonGlacier) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonGlacier/current/index.json"
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