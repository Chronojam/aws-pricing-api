package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonGlacier struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonGlacier_Product
	Terms           map[string]map[string]map[string]rawAmazonGlacier_Term
}

type rawAmazonGlacier_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonGlacier_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonGlacier) UnmarshalJSON(data []byte) error {
	var p rawAmazonGlacier
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonGlacier_Product{}
	terms := []*AmazonGlacier_Term{}

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
				pDimensions := []*AmazonGlacier_Term_PriceDimensions{}
				tAttributes := []*AmazonGlacier_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonGlacier_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonGlacier_Term{
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

type AmazonGlacier struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonGlacier_Product `gorm:"ForeignKey:AmazonGlacierID"`
	Terms           []*AmazonGlacier_Term    `gorm:"ForeignKey:AmazonGlacierID"`
}
type AmazonGlacier_Product struct {
	gorm.Model
	AmazonGlacierID uint
	Sku             string
	ProductFamily   string
	Attributes      AmazonGlacier_Product_Attributes `gorm:"ForeignKey:AmazonGlacier_Product_AttributesID"`
}
type AmazonGlacier_Product_Attributes struct {
	gorm.Model
	AmazonGlacier_Product_AttributesID uint
	Servicecode                        string
	Location                           string
	LocationType                       string
	FeeCode                            string
	FeeDescription                     string
	Usagetype                          string
	Operation                          string
	Servicename                        string
}

type AmazonGlacier_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonGlacierID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonGlacier_Term_PriceDimensions `gorm:"ForeignKey:AmazonGlacier_TermID"`
	TermAttributes  []*AmazonGlacier_Term_Attributes      `gorm:"ForeignKey:AmazonGlacier_TermID"`
}

type AmazonGlacier_Term_Attributes struct {
	gorm.Model
	AmazonGlacier_TermID uint
	Key                  string
	Value                string
}

type AmazonGlacier_Term_PriceDimensions struct {
	gorm.Model
	AmazonGlacier_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AmazonGlacier_Term_PricePerUnit `gorm:"ForeignKey:AmazonGlacier_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonGlacier_Term_PricePerUnit struct {
	gorm.Model
	AmazonGlacier_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AmazonGlacier) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonGlacier/current/index.json"
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
