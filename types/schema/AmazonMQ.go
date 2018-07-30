package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonMQ struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonMQ_Product
	Terms           map[string]map[string]map[string]rawAmazonMQ_Term
}

type rawAmazonMQ_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonMQ_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonMQ) UnmarshalJSON(data []byte) error {
	var p rawAmazonMQ
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonMQ_Product{}
	terms := []*AmazonMQ_Term{}

	// Convert from map to slice
	for i, _ := range p.Products {
		pr := p.Products[i]
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AmazonMQ_Term_PriceDimensions{}
				tAttributes := []*AmazonMQ_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonMQ_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonMQ_Term{
					OfferTermCode:   term.OfferTermCode,
					Sku:             term.Sku,
					EffectiveDate:   term.EffectiveDate,
					TermAttributes:  tAttributes,
					PriceDimensions: pDimensions,
				}

				terms = append(terms, &t)
			}
		}
	}

	l.FormatVersion = p.FormatVersion
	l.Disclaimer = p.Disclaimer
	l.OfferCode = p.OfferCode
	l.Version = p.Version
	l.PublicationDate = p.PublicationDate
	l.Products = products
	l.Terms = terms
	return nil
}

type AmazonMQ struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonMQ_Product `gorm:"ForeignKey:AmazonMQID"`
	Terms           []*AmazonMQ_Term    `gorm:"ForeignKey:AmazonMQID"`
}
type AmazonMQ_Product struct {
	gorm.Model
	AmazonMQID    uint
	Sku           string
	ProductFamily string
	Attributes    AmazonMQ_Product_Attributes `gorm:"ForeignKey:AmazonMQ_Product_AttributesID"`
}
type AmazonMQ_Product_Attributes struct {
	gorm.Model
	AmazonMQ_Product_AttributesID uint
	LocationType                  string
	BrokerEngine                  string
	NormalizationSizeFactor       string
	NetworkPerformance            string
	Usagetype                     string
	Servicecode                   string
	Vcpu                          string
	ClockSpeed                    string
	LicenseModel                  string
	EnhancedNetworkingSupport     string
	Location                      string
	Memory                        string
	DeploymentOption              string
	Operation                     string
	Servicename                   string
}

type AmazonMQ_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonMQID      uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonMQ_Term_PriceDimensions `gorm:"ForeignKey:AmazonMQ_TermID"`
	TermAttributes  []*AmazonMQ_Term_Attributes      `gorm:"ForeignKey:AmazonMQ_TermID"`
}

type AmazonMQ_Term_Attributes struct {
	gorm.Model
	AmazonMQ_TermID uint
	Key             string
	Value           string
}

type AmazonMQ_Term_PriceDimensions struct {
	gorm.Model
	AmazonMQ_TermID uint
	RateCode        string
	RateType        string
	Description     string
	BeginRange      string
	EndRange        string
	Unit            string
	PricePerUnit    *AmazonMQ_Term_PricePerUnit `gorm:"ForeignKey:AmazonMQ_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonMQ_Term_PricePerUnit struct {
	gorm.Model
	AmazonMQ_Term_PriceDimensionsID uint
	USD                             string
}

func (a *AmazonMQ) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonMQ/current/index.json"
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
