package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSCloudTrail struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSCloudTrail_Product
	Terms           map[string]map[string]map[string]rawAWSCloudTrail_Term
}

type rawAWSCloudTrail_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSCloudTrail_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSCloudTrail) UnmarshalJSON(data []byte) error {
	var p rawAWSCloudTrail
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSCloudTrail_Product{}
	terms := []*AWSCloudTrail_Term{}

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
				pDimensions := []*AWSCloudTrail_Term_PriceDimensions{}
				tAttributes := []*AWSCloudTrail_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSCloudTrail_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSCloudTrail_Term{
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

type AWSCloudTrail struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSCloudTrail_Product `gorm:"ForeignKey:AWSCloudTrailID"`
	Terms           []*AWSCloudTrail_Term    `gorm:"ForeignKey:AWSCloudTrailID"`
}
type AWSCloudTrail_Product struct {
	gorm.Model
	AWSCloudTrailID uint
	Sku             string
	ProductFamily   string
	Attributes      AWSCloudTrail_Product_Attributes `gorm:"ForeignKey:AWSCloudTrail_Product_AttributesID"`
}
type AWSCloudTrail_Product_Attributes struct {
	gorm.Model
	AWSCloudTrail_Product_AttributesID uint
	Servicecode                        string
	Location                           string
	LocationType                       string
	Usagetype                          string
	Operation                          string
}

type AWSCloudTrail_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSCloudTrailID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSCloudTrail_Term_PriceDimensions `gorm:"ForeignKey:AWSCloudTrail_TermID"`
	TermAttributes  []*AWSCloudTrail_Term_Attributes      `gorm:"ForeignKey:AWSCloudTrail_TermID"`
}

type AWSCloudTrail_Term_Attributes struct {
	gorm.Model
	AWSCloudTrail_TermID uint
	Key                  string
	Value                string
}

type AWSCloudTrail_Term_PriceDimensions struct {
	gorm.Model
	AWSCloudTrail_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AWSCloudTrail_Term_PricePerUnit `gorm:"ForeignKey:AWSCloudTrail_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSCloudTrail_Term_PricePerUnit struct {
	gorm.Model
	AWSCloudTrail_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AWSCloudTrail) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCloudTrail/current/index.json"
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
