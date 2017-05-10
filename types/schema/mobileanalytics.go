package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Mobileanalytics struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]Mobileanalytics_Product
	Terms		map[string]map[string]Mobileanalytics_Term
}
type Mobileanalytics_Product struct {	Sku	string
	ProductFamily	string
	Attributes	Mobileanalytics_Product_Attributes
}
type Mobileanalytics_Product_Attributes struct {	Servicecode	string
	Description	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	IncludedEvents	string
}

type Mobileanalytics_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions Mobileanalytics_Term_PriceDimensions
	TermAttributes Mobileanalytics_Term_TermAttributes
}

type Mobileanalytics_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	Mobileanalytics_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type Mobileanalytics_Term_PricePerUnit struct {
	USD	string
}

type Mobileanalytics_Term_TermAttributes struct {

}
func (a *Mobileanalytics) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/mobileanalytics/current/index.json"
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