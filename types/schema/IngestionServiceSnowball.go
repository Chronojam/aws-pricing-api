package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type IngestionServiceSnowball struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]IngestionServiceSnowball_Product
	Terms		map[string]map[string]IngestionServiceSnowball_Term
}
type IngestionServiceSnowball_Product struct {	Attributes	IngestionServiceSnowball_Product_Attributes
	Sku	string
	ProductFamily	string
}
type IngestionServiceSnowball_Product_Attributes struct {	GroupDescription	string
	ToLocationType	string
	Usagetype	string
	Servicecode	string
	Group	string
	FromLocationType	string
	ToLocation	string
	Operation	string
	TransferType	string
	FromLocation	string
}

type IngestionServiceSnowball_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions IngestionServiceSnowball_Term_PriceDimensions
	TermAttributes IngestionServiceSnowball_Term_TermAttributes
}

type IngestionServiceSnowball_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	IngestionServiceSnowball_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type IngestionServiceSnowball_Term_PricePerUnit struct {
	USD	string
}

type IngestionServiceSnowball_Term_TermAttributes struct {

}
func (a *IngestionServiceSnowball) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/IngestionServiceSnowball/current/index.json"
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