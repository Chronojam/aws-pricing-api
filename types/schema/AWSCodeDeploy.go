package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSCodeDeploy struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSCodeDeploy_Product
	Terms		map[string]map[string]AWSCodeDeploy_Term
}
type AWSCodeDeploy_Product struct {	ProductFamily	string
	Attributes	AWSCodeDeploy_Product_Attributes
	Sku	string
}
type AWSCodeDeploy_Product_Attributes struct {	DeploymentLocation	string
	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
}

type AWSCodeDeploy_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSCodeDeploy_Term_PriceDimensions
	TermAttributes AWSCodeDeploy_Term_TermAttributes
}

type AWSCodeDeploy_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSCodeDeploy_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSCodeDeploy_Term_PricePerUnit struct {
	USD	string
}

type AWSCodeDeploy_Term_TermAttributes struct {

}
func (a *AWSCodeDeploy) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCodeDeploy/current/index.json"
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