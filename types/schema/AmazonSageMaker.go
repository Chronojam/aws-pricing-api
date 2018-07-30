package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonSageMaker struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonSageMaker_Product
	Terms           map[string]map[string]map[string]rawAmazonSageMaker_Term
}

type rawAmazonSageMaker_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonSageMaker_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonSageMaker) UnmarshalJSON(data []byte) error {
	var p rawAmazonSageMaker
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonSageMaker_Product{}
	terms := []*AmazonSageMaker_Term{}

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
				pDimensions := []*AmazonSageMaker_Term_PriceDimensions{}
				tAttributes := []*AmazonSageMaker_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonSageMaker_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonSageMaker_Term{
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

type AmazonSageMaker struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonSageMaker_Product `gorm:"ForeignKey:AmazonSageMakerID"`
	Terms           []*AmazonSageMaker_Term    `gorm:"ForeignKey:AmazonSageMakerID"`
}
type AmazonSageMaker_Product struct {
	gorm.Model
	AmazonSageMakerID uint
	Sku               string
	ProductFamily     string
	Attributes        AmazonSageMaker_Product_Attributes `gorm:"ForeignKey:AmazonSageMaker_Product_AttributesID"`
}
type AmazonSageMaker_Product_Attributes struct {
	gorm.Model
	AmazonSageMaker_Product_AttributesID uint
	Usagetype                            string
	VCpu                                 string
	Location                             string
	InstanceType                         string
	Memory                               string
	Operation                            string
	Gpu                                  string
	PhysicalCpu                          string
	PhysicalGpu                          string
	Servicecode                          string
	LocationType                         string
	NetworkPerformance                   string
	Servicename                          string
	GpuMemory                            string
}

type AmazonSageMaker_Term struct {
	gorm.Model
	OfferTermCode     string
	AmazonSageMakerID uint
	Sku               string
	EffectiveDate     string
	PriceDimensions   []*AmazonSageMaker_Term_PriceDimensions `gorm:"ForeignKey:AmazonSageMaker_TermID"`
	TermAttributes    []*AmazonSageMaker_Term_Attributes      `gorm:"ForeignKey:AmazonSageMaker_TermID"`
}

type AmazonSageMaker_Term_Attributes struct {
	gorm.Model
	AmazonSageMaker_TermID uint
	Key                    string
	Value                  string
}

type AmazonSageMaker_Term_PriceDimensions struct {
	gorm.Model
	AmazonSageMaker_TermID uint
	RateCode               string
	RateType               string
	Description            string
	BeginRange             string
	EndRange               string
	Unit                   string
	PricePerUnit           *AmazonSageMaker_Term_PricePerUnit `gorm:"ForeignKey:AmazonSageMaker_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonSageMaker_Term_PricePerUnit struct {
	gorm.Model
	AmazonSageMaker_Term_PriceDimensionsID uint
	USD                                    string
}

func (a *AmazonSageMaker) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonSageMaker/current/index.json"
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
