package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonKinesis struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonKinesis_Product
	Terms		map[string]map[string]AmazonKinesis_Term
}
type AmazonKinesis_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonKinesis_Product_Attributes
}
type AmazonKinesis_Product_Attributes struct {	StandardStorageRetentionIncluded	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
	MaximumExtendedStorage	string
}

type AmazonKinesis_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonKinesis_Term_PriceDimensions
	TermAttributes AmazonKinesis_Term_TermAttributes
}

type AmazonKinesis_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonKinesis_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonKinesis_Term_PricePerUnit struct {
	USD	string
}

type AmazonKinesis_Term_TermAttributes struct {

}
func (a *AmazonKinesis) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesis/current/index.json"
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