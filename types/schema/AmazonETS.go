package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonETS struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonETS_Product
	Terms		map[string]map[string]AmazonETS_Term
}
type AmazonETS_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonETS_Product_Attributes
}
type AmazonETS_Product_Attributes struct {	TranscodingResult	string
	VideoResolution	string
	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
}

type AmazonETS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonETS_Term_PriceDimensions
	TermAttributes AmazonETS_Term_TermAttributes
}

type AmazonETS_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonETS_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonETS_Term_PricePerUnit struct {
	USD	string
}

type AmazonETS_Term_TermAttributes struct {

}
func (a *AmazonETS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonETS/current/index.json"
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