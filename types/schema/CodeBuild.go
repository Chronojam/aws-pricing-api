package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type CodeBuild struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]CodeBuild_Product
	Terms		map[string]map[string]CodeBuild_Term
}
type CodeBuild_Product struct {	Attributes	CodeBuild_Product_Attributes
	Sku	string
	ProductFamily	string
}
type CodeBuild_Product_Attributes struct {	Servicecode	string
	Location	string
	Vcpu	string
	Memory	string
	ComputeType	string
	LocationType	string
	OperatingSystem	string
	Usagetype	string
	Operation	string
	ComputeFamily	string
}

type CodeBuild_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions CodeBuild_Term_PriceDimensions
	TermAttributes CodeBuild_Term_TermAttributes
}

type CodeBuild_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	CodeBuild_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type CodeBuild_Term_PricePerUnit struct {
	USD	string
}

type CodeBuild_Term_TermAttributes struct {

}
func (a *CodeBuild) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/CodeBuild/current/index.json"
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