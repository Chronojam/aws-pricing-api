package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSShield struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSShield_Product
	Terms           map[string]map[string]map[string]rawAWSShield_Term
}

type rawAWSShield_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSShield_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSShield) UnmarshalJSON(data []byte) error {
	var p rawAWSShield
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSShield_Product{}
	terms := []*AWSShield_Term{}

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
				pDimensions := []*AWSShield_Term_PriceDimensions{}
				tAttributes := []*AWSShield_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSShield_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSShield_Term{
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

type AWSShield struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSShield_Product `gorm:"ForeignKey:AWSShieldID"`
	Terms           []*AWSShield_Term    `gorm:"ForeignKey:AWSShieldID"`
}
type AWSShield_Product struct {
	gorm.Model
	AWSShieldID   uint
	Sku           string
	ProductFamily string
	Attributes    AWSShield_Product_Attributes `gorm:"ForeignKey:AWSShield_Product_AttributesID"`
}
type AWSShield_Product_Attributes struct {
	gorm.Model
	AWSShield_Product_AttributesID uint
	Servicecode                    string
	FromLocation                   string
	ToLocation                     string
	ToLocationType                 string
	Usagetype                      string
	ResourceType                   string
	Servicename                    string
	FromLocationType               string
	Operation                      string
}

type AWSShield_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSShieldID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSShield_Term_PriceDimensions `gorm:"ForeignKey:AWSShield_TermID"`
	TermAttributes  []*AWSShield_Term_Attributes      `gorm:"ForeignKey:AWSShield_TermID"`
}

type AWSShield_Term_Attributes struct {
	gorm.Model
	AWSShield_TermID uint
	Key              string
	Value            string
}

type AWSShield_Term_PriceDimensions struct {
	gorm.Model
	AWSShield_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AWSShield_Term_PricePerUnit `gorm:"ForeignKey:AWSShield_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSShield_Term_PricePerUnit struct {
	gorm.Model
	AWSShield_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AWSShield) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSShield/current/index.json"
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
