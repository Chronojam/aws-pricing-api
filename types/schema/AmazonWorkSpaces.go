package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonWorkSpaces struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonWorkSpaces_Product
	Terms           map[string]map[string]map[string]rawAmazonWorkSpaces_Term
}

type rawAmazonWorkSpaces_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonWorkSpaces_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonWorkSpaces) UnmarshalJSON(data []byte) error {
	var p rawAmazonWorkSpaces
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonWorkSpaces_Product{}
	terms := []*AmazonWorkSpaces_Term{}

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
				pDimensions := []*AmazonWorkSpaces_Term_PriceDimensions{}
				tAttributes := []*AmazonWorkSpaces_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonWorkSpaces_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonWorkSpaces_Term{
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

type AmazonWorkSpaces struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonWorkSpaces_Product `gorm:"ForeignKey:AmazonWorkSpacesID"`
	Terms           []*AmazonWorkSpaces_Term    `gorm:"ForeignKey:AmazonWorkSpacesID"`
}
type AmazonWorkSpaces_Product struct {
	gorm.Model
	AmazonWorkSpacesID uint
	Sku                string
	ProductFamily      string
	Attributes         AmazonWorkSpaces_Product_Attributes `gorm:"ForeignKey:AmazonWorkSpaces_Product_AttributesID"`
}
type AmazonWorkSpaces_Product_Attributes struct {
	gorm.Model
	AmazonWorkSpaces_Product_AttributesID uint
	Storage                               string
	Group                                 string
	ResourceType                          string
	SoftwareIncluded                      string
	Vcpu                                  string
	Operation                             string
	Bundle                                string
	GroupDescription                      string
	Usagetype                             string
	License                               string
	RunningMode                           string
	Servicecode                           string
	Location                              string
	LocationType                          string
	Memory                                string
}

type AmazonWorkSpaces_Term struct {
	gorm.Model
	OfferTermCode      string
	AmazonWorkSpacesID uint
	Sku                string
	EffectiveDate      string
	PriceDimensions    []*AmazonWorkSpaces_Term_PriceDimensions `gorm:"ForeignKey:AmazonWorkSpaces_TermID"`
	TermAttributes     []*AmazonWorkSpaces_Term_Attributes      `gorm:"ForeignKey:AmazonWorkSpaces_TermID"`
}

type AmazonWorkSpaces_Term_Attributes struct {
	gorm.Model
	AmazonWorkSpaces_TermID uint
	Key                     string
	Value                   string
}

type AmazonWorkSpaces_Term_PriceDimensions struct {
	gorm.Model
	AmazonWorkSpaces_TermID uint
	RateCode                string
	RateType                string
	Description             string
	BeginRange              string
	EndRange                string
	Unit                    string
	PricePerUnit            *AmazonWorkSpaces_Term_PricePerUnit `gorm:"ForeignKey:AmazonWorkSpaces_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonWorkSpaces_Term_PricePerUnit struct {
	gorm.Model
	AmazonWorkSpaces_Term_PriceDimensionsID uint
	USD                                     string
}

func (a *AmazonWorkSpaces) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkSpaces/current/index.json"
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
