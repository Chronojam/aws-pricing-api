package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type OpsWorks struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]OpsWorks_Product
	Terms		map[string]map[string]OpsWorks_Term
}
type OpsWorks_Product struct {	Sku	string
	ProductFamily	string
	Attributes	OpsWorks_Product_Attributes
}
type OpsWorks_Product_Attributes struct {	ServerLocation	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	Usagetype	string
	Operation	string
}

type OpsWorks_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions OpsWorks_Term_PriceDimensions
	TermAttributes OpsWorks_Term_TermAttributes
}

type OpsWorks_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	OpsWorks_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type OpsWorks_Term_PricePerUnit struct {
	USD	string
}

type OpsWorks_Term_TermAttributes struct {

}
func (a *OpsWorks) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/OpsWorks/current/index.json"
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