package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawTranscribe struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]Transcribe_Product
	Terms           map[string]map[string]map[string]rawTranscribe_Term
}

type rawTranscribe_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]Transcribe_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *Transcribe) UnmarshalJSON(data []byte) error {
	var p rawTranscribe
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*Transcribe_Product{}
	terms := []*Transcribe_Term{}

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
				pDimensions := []*Transcribe_Term_PriceDimensions{}
				tAttributes := []*Transcribe_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := Transcribe_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := Transcribe_Term{
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

type Transcribe struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*Transcribe_Product `gorm:"ForeignKey:TranscribeID"`
	Terms           []*Transcribe_Term    `gorm:"ForeignKey:TranscribeID"`
}
type Transcribe_Product struct {
	gorm.Model
	TranscribeID  uint
	Sku           string
	ProductFamily string
	Attributes    Transcribe_Product_Attributes `gorm:"ForeignKey:Transcribe_Product_AttributesID"`
}
type Transcribe_Product_Attributes struct {
	gorm.Model
	Transcribe_Product_AttributesID uint
	Servicecode                     string
	Location                        string
	LocationType                    string
	Usagetype                       string
	Operation                       string
	Servicename                     string
	SupportedModes                  string
}

type Transcribe_Term struct {
	gorm.Model
	OfferTermCode   string
	TranscribeID    uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*Transcribe_Term_PriceDimensions `gorm:"ForeignKey:Transcribe_TermID"`
	TermAttributes  []*Transcribe_Term_Attributes      `gorm:"ForeignKey:Transcribe_TermID"`
}

type Transcribe_Term_Attributes struct {
	gorm.Model
	Transcribe_TermID uint
	Key               string
	Value             string
}

type Transcribe_Term_PriceDimensions struct {
	gorm.Model
	Transcribe_TermID uint
	RateCode          string
	RateType          string
	Description       string
	BeginRange        string
	EndRange          string
	Unit              string
	PricePerUnit      *Transcribe_Term_PricePerUnit `gorm:"ForeignKey:Transcribe_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type Transcribe_Term_PricePerUnit struct {
	gorm.Model
	Transcribe_Term_PriceDimensionsID uint
	USD                               string
}

func (a *Transcribe) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/transcribe/current/index.json"
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
