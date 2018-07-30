package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSXRay struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSXRay_Product
	Terms           map[string]map[string]map[string]rawAWSXRay_Term
}

type rawAWSXRay_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSXRay_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSXRay) UnmarshalJSON(data []byte) error {
	var p rawAWSXRay
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSXRay_Product{}
	terms := []*AWSXRay_Term{}

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
				pDimensions := []*AWSXRay_Term_PriceDimensions{}
				tAttributes := []*AWSXRay_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSXRay_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSXRay_Term{
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

type AWSXRay struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSXRay_Product `gorm:"ForeignKey:AWSXRayID"`
	Terms           []*AWSXRay_Term    `gorm:"ForeignKey:AWSXRayID"`
}
type AWSXRay_Product struct {
	gorm.Model
	AWSXRayID     uint
	Sku           string
	ProductFamily string
	Attributes    AWSXRay_Product_Attributes `gorm:"ForeignKey:AWSXRay_Product_AttributesID"`
}
type AWSXRay_Product_Attributes struct {
	gorm.Model
	AWSXRay_Product_AttributesID uint
	Operation                    string
	Servicecode                  string
	Location                     string
	LocationType                 string
	Group                        string
	GroupDescription             string
	Usagetype                    string
}

type AWSXRay_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSXRayID       uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSXRay_Term_PriceDimensions `gorm:"ForeignKey:AWSXRay_TermID"`
	TermAttributes  []*AWSXRay_Term_Attributes      `gorm:"ForeignKey:AWSXRay_TermID"`
}

type AWSXRay_Term_Attributes struct {
	gorm.Model
	AWSXRay_TermID uint
	Key            string
	Value          string
}

type AWSXRay_Term_PriceDimensions struct {
	gorm.Model
	AWSXRay_TermID uint
	RateCode       string
	RateType       string
	Description    string
	BeginRange     string
	EndRange       string
	Unit           string
	PricePerUnit   *AWSXRay_Term_PricePerUnit `gorm:"ForeignKey:AWSXRay_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSXRay_Term_PricePerUnit struct {
	gorm.Model
	AWSXRay_Term_PriceDimensionsID uint
	USD                            string
}

func (a *AWSXRay) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSXRay/current/index.json"
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
