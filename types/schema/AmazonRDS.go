package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AmazonRDS struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonRDS_Product
	Terms		map[string]map[string]AmazonRDS_Term
}
type AmazonRDS_Product struct {	ProductFamily	string
	Attributes	AmazonRDS_Product_Attributes
	Sku	string
}
type AmazonRDS_Product_Attributes struct {	Location	string
	LocationType	string
	InstanceType	string
	CurrentGeneration	string
	DatabaseEngine	string
	DeploymentOption	string
	Usagetype	string
	Operation	string
	Servicecode	string
	Vcpu	string
	Storage	string
	NetworkPerformance	string
	EngineCode	string
	InstanceFamily	string
	ClockSpeed	string
	DatabaseEdition	string
	LicenseModel	string
	PhysicalProcessor	string
	Memory	string
	ProcessorArchitecture	string
	ProcessorFeatures	string
}

type AmazonRDS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AmazonRDS_Term_PriceDimensions
	TermAttributes AmazonRDS_Term_TermAttributes
}

type AmazonRDS_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonRDS_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonRDS_Term_PricePerUnit struct {
	USD	string
}

type AmazonRDS_Term_TermAttributes struct {

}
func (a *AmazonRDS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRDS/current/index.json"
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