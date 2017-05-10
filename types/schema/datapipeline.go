package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Datapipeline struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]Datapipeline_Product
	Terms		map[string]map[string]Datapipeline_Term
}
type Datapipeline_Product struct {	Sku	string
	ProductFamily	string
	Attributes	Datapipeline_Product_Attributes
}
type Datapipeline_Product_Attributes struct {	Group	string
	Usagetype	string
	Operation	string
	ExecutionFrequency	string
	FrequencyMode	string
	Servicecode	string
	LocationType	string
	Location	string
	ExecutionLocation	string
}

type Datapipeline_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions Datapipeline_Term_PriceDimensions
	TermAttributes Datapipeline_Term_TermAttributes
}

type Datapipeline_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	Datapipeline_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type Datapipeline_Term_PricePerUnit struct {
	USD	string
}

type Datapipeline_Term_TermAttributes struct {

}
func (a *Datapipeline) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/datapipeline/current/index.json"
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