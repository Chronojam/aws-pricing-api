package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSGlue struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSGlue_Product
	Terms           map[string]map[string]map[string]rawAWSGlue_Term
}

type rawAWSGlue_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSGlue_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSGlue) UnmarshalJSON(data []byte) error {
	var p rawAWSGlue
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSGlue_Product{}
	terms := []*AWSGlue_Term{}

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
				pDimensions := []*AWSGlue_Term_PriceDimensions{}
				tAttributes := []*AWSGlue_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSGlue_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSGlue_Term{
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

type AWSGlue struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSGlue_Product `gorm:"ForeignKey:AWSGlueID"`
	Terms           []*AWSGlue_Term    `gorm:"ForeignKey:AWSGlueID"`
}
type AWSGlue_Product struct {
	gorm.Model
	AWSGlueID     uint
	Sku           string
	ProductFamily string
	Attributes    AWSGlue_Product_Attributes `gorm:"ForeignKey:AWSGlue_Product_AttributesID"`
}
type AWSGlue_Product_Attributes struct {
	gorm.Model
	AWSGlue_Product_AttributesID uint
	Servicecode                  string
	Location                     string
	LocationType                 string
	Group                        string
	Usagetype                    string
	Operation                    string
	Servicename                  string
}

type AWSGlue_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSGlueID       uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSGlue_Term_PriceDimensions `gorm:"ForeignKey:AWSGlue_TermID"`
	TermAttributes  []*AWSGlue_Term_Attributes      `gorm:"ForeignKey:AWSGlue_TermID"`
}

type AWSGlue_Term_Attributes struct {
	gorm.Model
	AWSGlue_TermID uint
	Key            string
	Value          string
}

type AWSGlue_Term_PriceDimensions struct {
	gorm.Model
	AWSGlue_TermID uint
	RateCode       string
	RateType       string
	Description    string
	BeginRange     string
	EndRange       string
	Unit           string
	PricePerUnit   *AWSGlue_Term_PricePerUnit `gorm:"ForeignKey:AWSGlue_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSGlue_Term_PricePerUnit struct {
	gorm.Model
	AWSGlue_Term_PriceDimensionsID uint
	USD                            string
}

func (a *AWSGlue) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSGlue/current/index.json"
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
