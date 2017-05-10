package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AlexaWebInfoService struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AlexaWebInfoService_Product
	Terms		map[string]map[string]AlexaWebInfoService_Term
}
type AlexaWebInfoService_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AlexaWebInfoService_Product_Attributes
}
type AlexaWebInfoService_Product_Attributes struct {	Operation	string
	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
}

type AlexaWebInfoService_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AlexaWebInfoService_Term_PriceDimensions
	TermAttributes AlexaWebInfoService_Term_TermAttributes
}

type AlexaWebInfoService_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AlexaWebInfoService_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AlexaWebInfoService_Term_PricePerUnit struct {
	USD	string
}

type AlexaWebInfoService_Term_TermAttributes struct {

}
func (a *AlexaWebInfoService) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AlexaWebInfoService/current/index.json"
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