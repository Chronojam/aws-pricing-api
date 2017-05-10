package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSDeviceFarm struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDeviceFarm_Product
	Terms		map[string]map[string]AWSDeviceFarm_Term
}
type AWSDeviceFarm_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSDeviceFarm_Product_Attributes
}
type AWSDeviceFarm_Product_Attributes struct {	Usagetype	string
	DeviceOs	string
	ExecutionMode	string
	Servicecode	string
	Description	string
	Operation	string
	MeterMode	string
	Location	string
	LocationType	string
}

type AWSDeviceFarm_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSDeviceFarm_Term_PriceDimensions
	TermAttributes AWSDeviceFarm_Term_TermAttributes
}

type AWSDeviceFarm_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDeviceFarm_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSDeviceFarm_Term_PricePerUnit struct {
	USD	string
}

type AWSDeviceFarm_Term_TermAttributes struct {

}
func (a *AWSDeviceFarm) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDeviceFarm/current/index.json"
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