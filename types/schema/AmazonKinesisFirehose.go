package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonKinesisFirehose struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonKinesisFirehose_Product
	Terms		map[string]map[string]AmazonKinesisFirehose_Term
}
type AmazonKinesisFirehose_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonKinesisFirehose_Product_Attributes
}
type AmazonKinesisFirehose_Product_Attributes struct {	Servicecode	string
	Description	string
	Location	string
	LocationType	string
	Group	string
	Usagetype	string
	Operation	string
}

type AmazonKinesisFirehose_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonKinesisFirehose_Term_PriceDimensions
	TermAttributes AmazonKinesisFirehose_Term_TermAttributes
}

type AmazonKinesisFirehose_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonKinesisFirehose_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonKinesisFirehose_Term_PricePerUnit struct {
	USD	string
}

type AmazonKinesisFirehose_Term_TermAttributes struct {

}
func (a *AmazonKinesisFirehose) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesisFirehose/current/index.json"
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