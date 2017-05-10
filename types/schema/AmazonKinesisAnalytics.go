package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonKinesisAnalytics struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonKinesisAnalytics_Product
	Terms		map[string]map[string]AmazonKinesisAnalytics_Term
}
type AmazonKinesisAnalytics_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonKinesisAnalytics_Product_Attributes
}
type AmazonKinesisAnalytics_Product_Attributes struct {	Description	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AmazonKinesisAnalytics_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonKinesisAnalytics_Term_PriceDimensions
	TermAttributes AmazonKinesisAnalytics_Term_TermAttributes
}

type AmazonKinesisAnalytics_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonKinesisAnalytics_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonKinesisAnalytics_Term_PricePerUnit struct {
	USD	string
}

type AmazonKinesisAnalytics_Term_TermAttributes struct {

}
func (a *AmazonKinesisAnalytics) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesisAnalytics/current/index.json"
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