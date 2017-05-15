package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawIngestionServiceSnowball struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]IngestionServiceSnowball_Product
	Terms           map[string]map[string]map[string]rawIngestionServiceSnowball_Term
}

type rawIngestionServiceSnowball_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]IngestionServiceSnowball_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *IngestionServiceSnowball) UnmarshalJSON(data []byte) error {
	var p rawIngestionServiceSnowball
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*IngestionServiceSnowball_Product{}
	terms := []*IngestionServiceSnowball_Term{}

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
				pDimensions := []*IngestionServiceSnowball_Term_PriceDimensions{}
				tAttributes := []*IngestionServiceSnowball_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := IngestionServiceSnowball_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := IngestionServiceSnowball_Term{
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

type IngestionServiceSnowball struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*IngestionServiceSnowball_Product `gorm:"ForeignKey:IngestionServiceSnowballID"`
	Terms           []*IngestionServiceSnowball_Term    `gorm:"ForeignKey:IngestionServiceSnowballID"`
}
type IngestionServiceSnowball_Product struct {
	gorm.Model
	IngestionServiceSnowballID uint
	Sku                        string
	ProductFamily              string
	Attributes                 IngestionServiceSnowball_Product_Attributes `gorm:"ForeignKey:IngestionServiceSnowball_Product_AttributesID"`
}
type IngestionServiceSnowball_Product_Attributes struct {
	gorm.Model
	IngestionServiceSnowball_Product_AttributesID uint
	Servicecode                                   string
	Location                                      string
	LocationType                                  string
	Group                                         string
	GroupDescription                              string
	Usagetype                                     string
	Operation                                     string
	SnowballType                                  string
}

type IngestionServiceSnowball_Term struct {
	gorm.Model
	OfferTermCode              string
	IngestionServiceSnowballID uint
	Sku                        string
	EffectiveDate              string
	PriceDimensions            []*IngestionServiceSnowball_Term_PriceDimensions `gorm:"ForeignKey:IngestionServiceSnowball_TermID"`
	TermAttributes             []*IngestionServiceSnowball_Term_Attributes      `gorm:"ForeignKey:IngestionServiceSnowball_TermID"`
}

type IngestionServiceSnowball_Term_Attributes struct {
	gorm.Model
	IngestionServiceSnowball_TermID uint
	Key                             string
	Value                           string
}

type IngestionServiceSnowball_Term_PriceDimensions struct {
	gorm.Model
	IngestionServiceSnowball_TermID uint
	RateCode                        string
	RateType                        string
	Description                     string
	BeginRange                      string
	EndRange                        string
	Unit                            string
	PricePerUnit                    *IngestionServiceSnowball_Term_PricePerUnit `gorm:"ForeignKey:IngestionServiceSnowball_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type IngestionServiceSnowball_Term_PricePerUnit struct {
	gorm.Model
	IngestionServiceSnowball_Term_PriceDimensionsID uint
	USD                                             string
}

func (a *IngestionServiceSnowball) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/IngestionServiceSnowball/current/index.json"
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
