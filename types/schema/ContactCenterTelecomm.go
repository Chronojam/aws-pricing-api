package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawContactCenterTelecomm struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]ContactCenterTelecomm_Product
	Terms           map[string]map[string]map[string]rawContactCenterTelecomm_Term
}

type rawContactCenterTelecomm_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]ContactCenterTelecomm_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *ContactCenterTelecomm) UnmarshalJSON(data []byte) error {
	var p rawContactCenterTelecomm
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*ContactCenterTelecomm_Product{}
	terms := []*ContactCenterTelecomm_Term{}

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
				pDimensions := []*ContactCenterTelecomm_Term_PriceDimensions{}
				tAttributes := []*ContactCenterTelecomm_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := ContactCenterTelecomm_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := ContactCenterTelecomm_Term{
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

type ContactCenterTelecomm struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*ContactCenterTelecomm_Product `gorm:"ForeignKey:ContactCenterTelecommID"`
	Terms           []*ContactCenterTelecomm_Term    `gorm:"ForeignKey:ContactCenterTelecommID"`
}
type ContactCenterTelecomm_Product struct {
	gorm.Model
	ContactCenterTelecommID uint
	Sku                     string
	ProductFamily           string
	Attributes              ContactCenterTelecomm_Product_Attributes `gorm:"ForeignKey:ContactCenterTelecomm_Product_AttributesID"`
}
type ContactCenterTelecomm_Product_Attributes struct {
	gorm.Model
	ContactCenterTelecomm_Product_AttributesID uint
	Location                                   string
	LocationType                               string
	Usagetype                                  string
	Operation                                  string
	LineType                                   string
	Servicecode                                string
	Country                                    string
	Servicename                                string
	CallingType                                string
}

type ContactCenterTelecomm_Term struct {
	gorm.Model
	OfferTermCode           string
	ContactCenterTelecommID uint
	Sku                     string
	EffectiveDate           string
	PriceDimensions         []*ContactCenterTelecomm_Term_PriceDimensions `gorm:"ForeignKey:ContactCenterTelecomm_TermID"`
	TermAttributes          []*ContactCenterTelecomm_Term_Attributes      `gorm:"ForeignKey:ContactCenterTelecomm_TermID"`
}

type ContactCenterTelecomm_Term_Attributes struct {
	gorm.Model
	ContactCenterTelecomm_TermID uint
	Key                          string
	Value                        string
}

type ContactCenterTelecomm_Term_PriceDimensions struct {
	gorm.Model
	ContactCenterTelecomm_TermID uint
	RateCode                     string
	RateType                     string
	Description                  string
	BeginRange                   string
	EndRange                     string
	Unit                         string
	PricePerUnit                 *ContactCenterTelecomm_Term_PricePerUnit `gorm:"ForeignKey:ContactCenterTelecomm_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type ContactCenterTelecomm_Term_PricePerUnit struct {
	gorm.Model
	ContactCenterTelecomm_Term_PriceDimensionsID uint
	USD                                          string
}

func (a *ContactCenterTelecomm) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/ContactCenterTelecomm/current/index.json"
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
