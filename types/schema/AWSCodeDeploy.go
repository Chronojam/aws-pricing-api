package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSCodeDeploy struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSCodeDeploy_Product
	Terms           map[string]map[string]map[string]rawAWSCodeDeploy_Term
}

type rawAWSCodeDeploy_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSCodeDeploy_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSCodeDeploy) UnmarshalJSON(data []byte) error {
	var p rawAWSCodeDeploy
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSCodeDeploy_Product{}
	terms := []*AWSCodeDeploy_Term{}

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
				pDimensions := []*AWSCodeDeploy_Term_PriceDimensions{}
				tAttributes := []*AWSCodeDeploy_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSCodeDeploy_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSCodeDeploy_Term{
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

type AWSCodeDeploy struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSCodeDeploy_Product `gorm:"ForeignKey:AWSCodeDeployID"`
	Terms           []*AWSCodeDeploy_Term    `gorm:"ForeignKey:AWSCodeDeployID"`
}
type AWSCodeDeploy_Product struct {
	gorm.Model
	AWSCodeDeployID uint
	Sku             string
	ProductFamily   string
	Attributes      AWSCodeDeploy_Product_Attributes `gorm:"ForeignKey:AWSCodeDeploy_Product_AttributesID"`
}
type AWSCodeDeploy_Product_Attributes struct {
	gorm.Model
	AWSCodeDeploy_Product_AttributesID uint
	Servicecode                        string
	Location                           string
	LocationType                       string
	Usagetype                          string
	Operation                          string
	DeploymentLocation                 string
}

type AWSCodeDeploy_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSCodeDeployID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSCodeDeploy_Term_PriceDimensions `gorm:"ForeignKey:AWSCodeDeploy_TermID"`
	TermAttributes  []*AWSCodeDeploy_Term_Attributes      `gorm:"ForeignKey:AWSCodeDeploy_TermID"`
}

type AWSCodeDeploy_Term_Attributes struct {
	gorm.Model
	AWSCodeDeploy_TermID uint
	Key                  string
	Value                string
}

type AWSCodeDeploy_Term_PriceDimensions struct {
	gorm.Model
	AWSCodeDeploy_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AWSCodeDeploy_Term_PricePerUnit `gorm:"ForeignKey:AWSCodeDeploy_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSCodeDeploy_Term_PricePerUnit struct {
	gorm.Model
	AWSCodeDeploy_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AWSCodeDeploy) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCodeDeploy/current/index.json"
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
