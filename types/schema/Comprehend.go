package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawComprehend struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]Comprehend_Product
	Terms           map[string]map[string]map[string]rawComprehend_Term
}

type rawComprehend_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]Comprehend_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *Comprehend) UnmarshalJSON(data []byte) error {
	var p rawComprehend
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*Comprehend_Product{}
	terms := []*Comprehend_Term{}

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
				pDimensions := []*Comprehend_Term_PriceDimensions{}
				tAttributes := []*Comprehend_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := Comprehend_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := Comprehend_Term{
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

type Comprehend struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*Comprehend_Product `gorm:"ForeignKey:ComprehendID"`
	Terms           []*Comprehend_Term    `gorm:"ForeignKey:ComprehendID"`
}
type Comprehend_Product struct {
	gorm.Model
	ComprehendID  uint
	Sku           string
	ProductFamily string
	Attributes    Comprehend_Product_Attributes `gorm:"ForeignKey:Comprehend_Product_AttributesID"`
}
type Comprehend_Product_Attributes struct {
	gorm.Model
	Comprehend_Product_AttributesID uint
	Group                           string
	GroupDescription                string
	Usagetype                       string
	Operation                       string
	Servicename                     string
	Servicecode                     string
	Location                        string
	LocationType                    string
}

type Comprehend_Term struct {
	gorm.Model
	OfferTermCode   string
	ComprehendID    uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*Comprehend_Term_PriceDimensions `gorm:"ForeignKey:Comprehend_TermID"`
	TermAttributes  []*Comprehend_Term_Attributes      `gorm:"ForeignKey:Comprehend_TermID"`
}

type Comprehend_Term_Attributes struct {
	gorm.Model
	Comprehend_TermID uint
	Key               string
	Value             string
}

type Comprehend_Term_PriceDimensions struct {
	gorm.Model
	Comprehend_TermID uint
	RateCode          string
	RateType          string
	Description       string
	BeginRange        string
	EndRange          string
	Unit              string
	PricePerUnit      *Comprehend_Term_PricePerUnit `gorm:"ForeignKey:Comprehend_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type Comprehend_Term_PricePerUnit struct {
	gorm.Model
	Comprehend_Term_PriceDimensionsID uint
	USD                               string
}

func (a *Comprehend) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/comprehend/current/index.json"
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
