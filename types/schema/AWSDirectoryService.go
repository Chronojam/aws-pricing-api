package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSDirectoryService struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDirectoryService_Product
	Terms		map[string]map[string]AWSDirectoryService_Term
}
type AWSDirectoryService_Product struct {	Attributes	AWSDirectoryService_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AWSDirectoryService_Product_Attributes struct {	Usagetype	string
	Operation	string
	DirectorySize	string
	DirectoryType	string
	DirectoryTypeDescription	string
	Servicecode	string
	Location	string
	LocationType	string
}

type AWSDirectoryService_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSDirectoryService_Term_PriceDimensions
	TermAttributes AWSDirectoryService_Term_TermAttributes
}

type AWSDirectoryService_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDirectoryService_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSDirectoryService_Term_PricePerUnit struct {
	USD	string
}

type AWSDirectoryService_Term_TermAttributes struct {

}
func (a *AWSDirectoryService) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDirectoryService/current/index.json"
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