package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSCodeCommit struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSCodeCommit_Product
	Terms           map[string]map[string]map[string]rawAWSCodeCommit_Term
}

type rawAWSCodeCommit_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSCodeCommit_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSCodeCommit) UnmarshalJSON(data []byte) error {
	var p rawAWSCodeCommit
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSCodeCommit_Product{}
	terms := []*AWSCodeCommit_Term{}

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
				pDimensions := []*AWSCodeCommit_Term_PriceDimensions{}
				tAttributes := []*AWSCodeCommit_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSCodeCommit_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSCodeCommit_Term{
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

type AWSCodeCommit struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSCodeCommit_Product `gorm:"ForeignKey:AWSCodeCommitID"`
	Terms           []*AWSCodeCommit_Term    `gorm:"ForeignKey:AWSCodeCommitID"`
}
type AWSCodeCommit_Product struct {
	gorm.Model
	AWSCodeCommitID uint
	Sku             string
	ProductFamily   string
	Attributes      AWSCodeCommit_Product_Attributes `gorm:"ForeignKey:AWSCodeCommit_Product_AttributesID"`
}
type AWSCodeCommit_Product_Attributes struct {
	gorm.Model
	AWSCodeCommit_Product_AttributesID uint
	Group                              string
	Usagetype                          string
	Operation                          string
	Servicecode                        string
	Location                           string
	LocationType                       string
}

type AWSCodeCommit_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSCodeCommitID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSCodeCommit_Term_PriceDimensions `gorm:"ForeignKey:AWSCodeCommit_TermID"`
	TermAttributes  []*AWSCodeCommit_Term_Attributes      `gorm:"ForeignKey:AWSCodeCommit_TermID"`
}

type AWSCodeCommit_Term_Attributes struct {
	gorm.Model
	AWSCodeCommit_TermID uint
	Key                  string
	Value                string
}

type AWSCodeCommit_Term_PriceDimensions struct {
	gorm.Model
	AWSCodeCommit_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AWSCodeCommit_Term_PricePerUnit `gorm:"ForeignKey:AWSCodeCommit_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSCodeCommit_Term_PricePerUnit struct {
	gorm.Model
	AWSCodeCommit_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AWSCodeCommit) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCodeCommit/current/index.json"
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
