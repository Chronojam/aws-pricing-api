package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonWorkDocs struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWorkDocs_Product
	Terms		map[string]map[string]AmazonWorkDocs_Term
}
type AmazonWorkDocs_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonWorkDocs_Product_Attributes
}
type AmazonWorkDocs_Product_Attributes struct {	Description	string
	Location	string
	Storage	string
	Operation	string
	MaximumStorageVolume	string
	Servicecode	string
	LocationType	string
	Usagetype	string
	FreeTrial	string
	MinimumStorageVolume	string
}

type AmazonWorkDocs_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonWorkDocs_Term_PriceDimensions
	TermAttributes AmazonWorkDocs_Term_TermAttributes
}

type AmazonWorkDocs_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonWorkDocs_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonWorkDocs_Term_PricePerUnit struct {
	USD	string
}

type AmazonWorkDocs_Term_TermAttributes struct {

}
func (a *AmazonWorkDocs) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkDocs/current/index.json"
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