package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonPolly struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonPolly_Product
	Terms           map[string]map[string]map[string]rawAmazonPolly_Term
}

type rawAmazonPolly_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonPolly_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonPolly) UnmarshalJSON(data []byte) error {
	var p rawAmazonPolly
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonPolly_Product{}
	terms := []*AmazonPolly_Term{}

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
				pDimensions := []*AmazonPolly_Term_PriceDimensions{}
				tAttributes := []*AmazonPolly_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonPolly_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonPolly_Term{
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

type AmazonPolly struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonPolly_Product `gorm:"ForeignKey:AmazonPollyID"`
	Terms           []*AmazonPolly_Term    `gorm:"ForeignKey:AmazonPollyID"`
}
type AmazonPolly_Product struct {
	gorm.Model
	AmazonPollyID uint
	Sku           string
	ProductFamily string
	Attributes    AmazonPolly_Product_Attributes `gorm:"ForeignKey:AmazonPolly_Product_AttributesID"`
}
type AmazonPolly_Product_Attributes struct {
	gorm.Model
	AmazonPolly_Product_AttributesID uint
	Servicecode                      string
	Location                         string
	LocationType                     string
	Usagetype                        string
	Operation                        string
	Servicename                      string
}

type AmazonPolly_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonPollyID   uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonPolly_Term_PriceDimensions `gorm:"ForeignKey:AmazonPolly_TermID"`
	TermAttributes  []*AmazonPolly_Term_Attributes      `gorm:"ForeignKey:AmazonPolly_TermID"`
}

type AmazonPolly_Term_Attributes struct {
	gorm.Model
	AmazonPolly_TermID uint
	Key                string
	Value              string
}

type AmazonPolly_Term_PriceDimensions struct {
	gorm.Model
	AmazonPolly_TermID uint
	RateCode           string
	RateType           string
	Description        string
	BeginRange         string
	EndRange           string
	Unit               string
	PricePerUnit       *AmazonPolly_Term_PricePerUnit `gorm:"ForeignKey:AmazonPolly_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonPolly_Term_PricePerUnit struct {
	gorm.Model
	AmazonPolly_Term_PriceDimensionsID uint
	USD                                string
}

func (a *AmazonPolly) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonPolly/current/index.json"
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
