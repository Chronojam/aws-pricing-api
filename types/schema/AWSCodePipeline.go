package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSCodePipeline struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSCodePipeline_Product
	Terms		map[string]map[string]AWSCodePipeline_Term
}
type AWSCodePipeline_Product struct {	ProductFamily	string
	Attributes	AWSCodePipeline_Product_Attributes
	Sku	string
}
type AWSCodePipeline_Product_Attributes struct {	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	Description	string
}

type AWSCodePipeline_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSCodePipeline_Term_PriceDimensions
	TermAttributes AWSCodePipeline_Term_TermAttributes
}

type AWSCodePipeline_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSCodePipeline_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSCodePipeline_Term_PricePerUnit struct {
	USD	string
}

type AWSCodePipeline_Term_TermAttributes struct {

}
func (a *AWSCodePipeline) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCodePipeline/current/index.json"
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