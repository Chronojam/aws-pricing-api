package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type CloudHSM struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]CloudHSM_Product
	Terms		map[string]map[string]CloudHSM_Term
}
type CloudHSM_Product struct {	ProductFamily	string
	Attributes	CloudHSM_Product_Attributes
	Sku	string
}
type CloudHSM_Product_Attributes struct {	LocationType	string
	InstanceFamily	string
	Usagetype	string
	Operation	string
	TrialProduct	string
	UpfrontCommitment	string
	Servicecode	string
	Location	string
}

type CloudHSM_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions CloudHSM_Term_PriceDimensions
	TermAttributes CloudHSM_Term_TermAttributes
}

type CloudHSM_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	CloudHSM_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type CloudHSM_Term_PricePerUnit struct {
	USD	string
}

type CloudHSM_Term_TermAttributes struct {

}
func (a *CloudHSM) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/CloudHSM/current/index.json"
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