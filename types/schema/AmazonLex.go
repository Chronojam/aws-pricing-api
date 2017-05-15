package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonLex struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonLex_Product
	Terms           map[string]map[string]map[string]rawAmazonLex_Term
}

type rawAmazonLex_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonLex_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonLex) UnmarshalJSON(data []byte) error {
	var p rawAmazonLex
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonLex_Product{}
	terms := []*AmazonLex_Term{}

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
				pDimensions := []*AmazonLex_Term_PriceDimensions{}
				tAttributes := []*AmazonLex_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonLex_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonLex_Term{
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

type AmazonLex struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonLex_Product `gorm:"ForeignKey:AmazonLexID"`
	Terms           []*AmazonLex_Term    `gorm:"ForeignKey:AmazonLexID"`
}
type AmazonLex_Product struct {
	gorm.Model
	AmazonLexID   uint
	Sku           string
	ProductFamily string
	Attributes    AmazonLex_Product_Attributes `gorm:"ForeignKey:AmazonLex_Product_AttributesID"`
}
type AmazonLex_Product_Attributes struct {
	gorm.Model
	AmazonLex_Product_AttributesID uint
	InputMode                      string
	OutputMode                     string
	Location                       string
	Usagetype                      string
	Operation                      string
	GroupDescription               string
	SupportedModes                 string
	Servicecode                    string
	LocationType                   string
	Group                          string
}

type AmazonLex_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonLexID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonLex_Term_PriceDimensions `gorm:"ForeignKey:AmazonLex_TermID"`
	TermAttributes  []*AmazonLex_Term_Attributes      `gorm:"ForeignKey:AmazonLex_TermID"`
}

type AmazonLex_Term_Attributes struct {
	gorm.Model
	AmazonLex_TermID uint
	Key              string
	Value            string
}

type AmazonLex_Term_PriceDimensions struct {
	gorm.Model
	AmazonLex_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AmazonLex_Term_PricePerUnit `gorm:"ForeignKey:AmazonLex_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonLex_Term_PricePerUnit struct {
	gorm.Model
	AmazonLex_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AmazonLex) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonLex/current/index.json"
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
