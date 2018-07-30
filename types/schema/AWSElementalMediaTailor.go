package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSElementalMediaTailor struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSElementalMediaTailor_Product
	Terms           map[string]map[string]map[string]rawAWSElementalMediaTailor_Term
}

type rawAWSElementalMediaTailor_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSElementalMediaTailor_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSElementalMediaTailor) UnmarshalJSON(data []byte) error {
	var p rawAWSElementalMediaTailor
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSElementalMediaTailor_Product{}
	terms := []*AWSElementalMediaTailor_Term{}

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
				pDimensions := []*AWSElementalMediaTailor_Term_PriceDimensions{}
				tAttributes := []*AWSElementalMediaTailor_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSElementalMediaTailor_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSElementalMediaTailor_Term{
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

type AWSElementalMediaTailor struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSElementalMediaTailor_Product `gorm:"ForeignKey:AWSElementalMediaTailorID"`
	Terms           []*AWSElementalMediaTailor_Term    `gorm:"ForeignKey:AWSElementalMediaTailorID"`
}
type AWSElementalMediaTailor_Product struct {
	gorm.Model
	AWSElementalMediaTailorID uint
	Sku                       string
	ProductFamily             string
	Attributes                AWSElementalMediaTailor_Product_Attributes `gorm:"ForeignKey:AWSElementalMediaTailor_Product_AttributesID"`
}
type AWSElementalMediaTailor_Product_Attributes struct {
	gorm.Model
	AWSElementalMediaTailor_Product_AttributesID uint
	LocationType                                 string
	Usagetype                                    string
	Operation                                    string
	OperationType                                string
	Servicename                                  string
	Servicecode                                  string
	Location                                     string
}

type AWSElementalMediaTailor_Term struct {
	gorm.Model
	OfferTermCode             string
	AWSElementalMediaTailorID uint
	Sku                       string
	EffectiveDate             string
	PriceDimensions           []*AWSElementalMediaTailor_Term_PriceDimensions `gorm:"ForeignKey:AWSElementalMediaTailor_TermID"`
	TermAttributes            []*AWSElementalMediaTailor_Term_Attributes      `gorm:"ForeignKey:AWSElementalMediaTailor_TermID"`
}

type AWSElementalMediaTailor_Term_Attributes struct {
	gorm.Model
	AWSElementalMediaTailor_TermID uint
	Key                            string
	Value                          string
}

type AWSElementalMediaTailor_Term_PriceDimensions struct {
	gorm.Model
	AWSElementalMediaTailor_TermID uint
	RateCode                       string
	RateType                       string
	Description                    string
	BeginRange                     string
	EndRange                       string
	Unit                           string
	PricePerUnit                   *AWSElementalMediaTailor_Term_PricePerUnit `gorm:"ForeignKey:AWSElementalMediaTailor_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSElementalMediaTailor_Term_PricePerUnit struct {
	gorm.Model
	AWSElementalMediaTailor_Term_PriceDimensionsID uint
	USD                                            string
}

func (a *AWSElementalMediaTailor) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSElementalMediaTailor/current/index.json"
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
