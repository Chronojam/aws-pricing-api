package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Awskms struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]Awskms_Product
	Terms		map[string]map[string]Awskms_Term
}
type Awskms_Product struct {	Sku	string
	ProductFamily	string
	Attributes	Awskms_Product_Attributes
}
type Awskms_Product_Attributes struct {	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type Awskms_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions Awskms_Term_PriceDimensions
	TermAttributes Awskms_Term_TermAttributes
}

type Awskms_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	Awskms_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type Awskms_Term_PricePerUnit struct {
	USD	string
}

type Awskms_Term_TermAttributes struct {

}
func (a *Awskms) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/awskms/current/index.json"
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