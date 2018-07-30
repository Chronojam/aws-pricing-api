package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSSecretsManager struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSSecretsManager_Product
	Terms           map[string]map[string]map[string]rawAWSSecretsManager_Term
}

type rawAWSSecretsManager_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSSecretsManager_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSSecretsManager) UnmarshalJSON(data []byte) error {
	var p rawAWSSecretsManager
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSSecretsManager_Product{}
	terms := []*AWSSecretsManager_Term{}

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
				pDimensions := []*AWSSecretsManager_Term_PriceDimensions{}
				tAttributes := []*AWSSecretsManager_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSSecretsManager_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSSecretsManager_Term{
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

type AWSSecretsManager struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSSecretsManager_Product `gorm:"ForeignKey:AWSSecretsManagerID"`
	Terms           []*AWSSecretsManager_Term    `gorm:"ForeignKey:AWSSecretsManagerID"`
}
type AWSSecretsManager_Product struct {
	gorm.Model
	AWSSecretsManagerID uint
	Sku                 string
	ProductFamily       string
	Attributes          AWSSecretsManager_Product_Attributes `gorm:"ForeignKey:AWSSecretsManager_Product_AttributesID"`
}
type AWSSecretsManager_Product_Attributes struct {
	gorm.Model
	AWSSecretsManager_Product_AttributesID uint
	Servicecode                            string
	Location                               string
	LocationType                           string
	Group                                  string
	Usagetype                              string
	Operation                              string
	Servicename                            string
}

type AWSSecretsManager_Term struct {
	gorm.Model
	OfferTermCode       string
	AWSSecretsManagerID uint
	Sku                 string
	EffectiveDate       string
	PriceDimensions     []*AWSSecretsManager_Term_PriceDimensions `gorm:"ForeignKey:AWSSecretsManager_TermID"`
	TermAttributes      []*AWSSecretsManager_Term_Attributes      `gorm:"ForeignKey:AWSSecretsManager_TermID"`
}

type AWSSecretsManager_Term_Attributes struct {
	gorm.Model
	AWSSecretsManager_TermID uint
	Key                      string
	Value                    string
}

type AWSSecretsManager_Term_PriceDimensions struct {
	gorm.Model
	AWSSecretsManager_TermID uint
	RateCode                 string
	RateType                 string
	Description              string
	BeginRange               string
	EndRange                 string
	Unit                     string
	PricePerUnit             *AWSSecretsManager_Term_PricePerUnit `gorm:"ForeignKey:AWSSecretsManager_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSSecretsManager_Term_PricePerUnit struct {
	gorm.Model
	AWSSecretsManager_Term_PriceDimensionsID uint
	USD                                      string
}

func (a *AWSSecretsManager) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSSecretsManager/current/index.json"
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
