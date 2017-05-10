package schema

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type AmazonEC2 struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonEC2_Product
	Terms           map[string]map[string]AmazonEC2_Term
}
type AmazonEC2_Product struct {
	Sku           string
	ProductFamily string
	Attributes    AmazonEC2_Product_Attributes
}
type AmazonEC2_Product_Attributes struct {
	InstanceType                string
	Usagetype                   string
	LocationType                string
	CurrentGeneration           string
	Storage                     string
	OperatingSystem             string
	LicenseModel                string
	EnhancedNetworkingSupported string
	ProcessorFeatures           string
	Location                    string
	InstanceFamily              string
	PhysicalProcessor           string
	Memory                      string
	ProcessorArchitecture       string
	Servicecode                 string
	Vcpu                        string
	ClockSpeed                  string
	NetworkPerformance          string
	Tenancy                     string
	Operation                   string
	PreInstalledSw              string
}

type AmazonEC2_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions AmazonEC2_Term_PriceDimensions
	TermAttributes  AmazonEC2_Term_TermAttributes
}

type AmazonEC2_Term_PriceDimensions struct {
	RateCode     string
	RateType     string
	Description  string
	BeginRange   string
	EndRange     string
	Unit         string
	PricePerUnit AmazonEC2_Term_PricePerUnit
	AppliesTo    []interface{}
}

type AmazonEC2_Term_PricePerUnit struct {
	USD string
}

type AmazonEC2_Term_TermAttributes struct {
}

func (a *AmazonEC2) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json"
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
