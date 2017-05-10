package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonWorkMail struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWorkMail_Product
	Terms		map[string]map[string]AmazonWorkMail_Term
}
type AmazonWorkMail_Product struct {	ProductFamily	string
	Attributes	AmazonWorkMail_Product_Attributes
	Sku	string
}
type AmazonWorkMail_Product_Attributes struct {	LocationType	string
	Usagetype	string
	Operation	string
	FreeTier	string
	MailboxStorage	string
	Servicecode	string
	Location	string
}

type AmazonWorkMail_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonWorkMail_Term_PriceDimensions
	TermAttributes AmazonWorkMail_Term_TermAttributes
}

type AmazonWorkMail_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonWorkMail_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonWorkMail_Term_PricePerUnit struct {
	USD	string
}

type AmazonWorkMail_Term_TermAttributes struct {

}
func (a *AmazonWorkMail) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkMail/current/index.json"
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