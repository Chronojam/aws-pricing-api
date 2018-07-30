package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawTranslate struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]Translate_Product
	Terms           map[string]map[string]map[string]rawTranslate_Term
}

type rawTranslate_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]Translate_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *Translate) UnmarshalJSON(data []byte) error {
	var p rawTranslate
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*Translate_Product{}
	terms := []*Translate_Term{}

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
				pDimensions := []*Translate_Term_PriceDimensions{}
				tAttributes := []*Translate_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := Translate_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := Translate_Term{
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

type Translate struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*Translate_Product `gorm:"ForeignKey:TranslateID"`
	Terms           []*Translate_Term    `gorm:"ForeignKey:TranslateID"`
}
type Translate_Product struct {
	gorm.Model
	TranslateID   uint
	Sku           string
	ProductFamily string
	Attributes    Translate_Product_Attributes `gorm:"ForeignKey:Translate_Product_AttributesID"`
}
type Translate_Product_Attributes struct {
	gorm.Model
	Translate_Product_AttributesID uint
	OutputMode                     string
	Servicename                    string
	Servicecode                    string
	Location                       string
	LocationType                   string
	Usagetype                      string
	Operation                      string
	InputMode                      string
}

type Translate_Term struct {
	gorm.Model
	OfferTermCode   string
	TranslateID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*Translate_Term_PriceDimensions `gorm:"ForeignKey:Translate_TermID"`
	TermAttributes  []*Translate_Term_Attributes      `gorm:"ForeignKey:Translate_TermID"`
}

type Translate_Term_Attributes struct {
	gorm.Model
	Translate_TermID uint
	Key              string
	Value            string
}

type Translate_Term_PriceDimensions struct {
	gorm.Model
	Translate_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *Translate_Term_PricePerUnit `gorm:"ForeignKey:Translate_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type Translate_Term_PricePerUnit struct {
	gorm.Model
	Translate_Term_PriceDimensionsID uint
	USD                              string
}

func (a *Translate) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/translate/current/index.json"
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
