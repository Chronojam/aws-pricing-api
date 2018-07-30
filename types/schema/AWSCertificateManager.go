package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSCertificateManager struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSCertificateManager_Product
	Terms           map[string]map[string]map[string]rawAWSCertificateManager_Term
}

type rawAWSCertificateManager_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSCertificateManager_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSCertificateManager) UnmarshalJSON(data []byte) error {
	var p rawAWSCertificateManager
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSCertificateManager_Product{}
	terms := []*AWSCertificateManager_Term{}

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
				pDimensions := []*AWSCertificateManager_Term_PriceDimensions{}
				tAttributes := []*AWSCertificateManager_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSCertificateManager_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSCertificateManager_Term{
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

type AWSCertificateManager struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSCertificateManager_Product `gorm:"ForeignKey:AWSCertificateManagerID"`
	Terms           []*AWSCertificateManager_Term    `gorm:"ForeignKey:AWSCertificateManagerID"`
}
type AWSCertificateManager_Product struct {
	gorm.Model
	AWSCertificateManagerID uint
	Sku                     string
	ProductFamily           string
	Attributes              AWSCertificateManager_Product_Attributes `gorm:"ForeignKey:AWSCertificateManager_Product_AttributesID"`
}
type AWSCertificateManager_Product_Attributes struct {
	gorm.Model
	AWSCertificateManager_Product_AttributesID uint
	Location                                   string
	LocationType                               string
	GroupDescription                           string
	Usagetype                                  string
	Operation                                  string
	Servicecode                                string
	Group                                      string
	Servicename                                string
	Type                                       string
}

type AWSCertificateManager_Term struct {
	gorm.Model
	OfferTermCode           string
	AWSCertificateManagerID uint
	Sku                     string
	EffectiveDate           string
	PriceDimensions         []*AWSCertificateManager_Term_PriceDimensions `gorm:"ForeignKey:AWSCertificateManager_TermID"`
	TermAttributes          []*AWSCertificateManager_Term_Attributes      `gorm:"ForeignKey:AWSCertificateManager_TermID"`
}

type AWSCertificateManager_Term_Attributes struct {
	gorm.Model
	AWSCertificateManager_TermID uint
	Key                          string
	Value                        string
}

type AWSCertificateManager_Term_PriceDimensions struct {
	gorm.Model
	AWSCertificateManager_TermID uint
	RateCode                     string
	RateType                     string
	Description                  string
	BeginRange                   string
	EndRange                     string
	Unit                         string
	PricePerUnit                 *AWSCertificateManager_Term_PricePerUnit `gorm:"ForeignKey:AWSCertificateManager_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSCertificateManager_Term_PricePerUnit struct {
	gorm.Model
	AWSCertificateManager_Term_PriceDimensionsID uint
	USD                                          string
}

func (a *AWSCertificateManager) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCertificateManager/current/index.json"
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
