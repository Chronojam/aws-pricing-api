package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSDirectConnect struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDirectConnect_Product
	Terms		map[string]map[string]AWSDirectConnect_Term
}
type AWSDirectConnect_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSDirectConnect_Product_Attributes
}
type AWSDirectConnect_Product_Attributes struct {	TransferType	string
	FromLocation	string
	VirtualInterfaceType	string
	Usagetype	string
	Operation	string
	Version	string
	Servicecode	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
}

type AWSDirectConnect_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSDirectConnect_Term_PriceDimensions
	TermAttributes AWSDirectConnect_Term_TermAttributes
}

type AWSDirectConnect_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDirectConnect_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSDirectConnect_Term_PricePerUnit struct {
	USD	string
}

type AWSDirectConnect_Term_TermAttributes struct {

}
func (a *AWSDirectConnect) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDirectConnect/current/index.json"
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