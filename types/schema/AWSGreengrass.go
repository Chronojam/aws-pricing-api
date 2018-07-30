package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSGreengrass struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSGreengrass_Product
	Terms           map[string]map[string]map[string]rawAWSGreengrass_Term
}

type rawAWSGreengrass_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSGreengrass_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSGreengrass) UnmarshalJSON(data []byte) error {
	var p rawAWSGreengrass
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSGreengrass_Product{}
	terms := []*AWSGreengrass_Term{}

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
				pDimensions := []*AWSGreengrass_Term_PriceDimensions{}
				tAttributes := []*AWSGreengrass_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSGreengrass_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSGreengrass_Term{
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

type AWSGreengrass struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSGreengrass_Product `gorm:"ForeignKey:AWSGreengrassID"`
	Terms           []*AWSGreengrass_Term    `gorm:"ForeignKey:AWSGreengrassID"`
}
type AWSGreengrass_Product struct {
	gorm.Model
	AWSGreengrassID uint
	Sku             string
	ProductFamily   string
	Attributes      AWSGreengrass_Product_Attributes `gorm:"ForeignKey:AWSGreengrass_Product_AttributesID"`
}
type AWSGreengrass_Product_Attributes struct {
	gorm.Model
	AWSGreengrass_Product_AttributesID uint
	Usagetype                          string
	Operation                          string
	Servicename                        string
	TenancySupport                     string
	Servicecode                        string
	Location                           string
	LocationType                       string
}

type AWSGreengrass_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSGreengrassID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSGreengrass_Term_PriceDimensions `gorm:"ForeignKey:AWSGreengrass_TermID"`
	TermAttributes  []*AWSGreengrass_Term_Attributes      `gorm:"ForeignKey:AWSGreengrass_TermID"`
}

type AWSGreengrass_Term_Attributes struct {
	gorm.Model
	AWSGreengrass_TermID uint
	Key                  string
	Value                string
}

type AWSGreengrass_Term_PriceDimensions struct {
	gorm.Model
	AWSGreengrass_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AWSGreengrass_Term_PricePerUnit `gorm:"ForeignKey:AWSGreengrass_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSGreengrass_Term_PricePerUnit struct {
	gorm.Model
	AWSGreengrass_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AWSGreengrass) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSGreengrass/current/index.json"
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
