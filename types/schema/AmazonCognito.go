package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonCognito struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonCognito_Product
	Terms           map[string]map[string]map[string]rawAmazonCognito_Term
}

type rawAmazonCognito_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonCognito_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonCognito) UnmarshalJSON(data []byte) error {
	var p rawAmazonCognito
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonCognito_Product{}
	terms := []*AmazonCognito_Term{}

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
				pDimensions := []*AmazonCognito_Term_PriceDimensions{}
				tAttributes := []*AmazonCognito_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonCognito_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonCognito_Term{
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

type AmazonCognito struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonCognito_Product `gorm:"ForeignKey:AmazonCognitoID"`
	Terms           []*AmazonCognito_Term    `gorm:"ForeignKey:AmazonCognitoID"`
}
type AmazonCognito_Product struct {
	gorm.Model
	AmazonCognitoID uint
	Sku             string
	ProductFamily   string
	Attributes      AmazonCognito_Product_Attributes `gorm:"ForeignKey:AmazonCognito_Product_AttributesID"`
}
type AmazonCognito_Product_Attributes struct {
	gorm.Model
	AmazonCognito_Product_AttributesID uint
	Usagetype                          string
	Operation                          string
	Servicecode                        string
	Location                           string
	LocationType                       string
}

type AmazonCognito_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonCognitoID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonCognito_Term_PriceDimensions `gorm:"ForeignKey:AmazonCognito_TermID"`
	TermAttributes  []*AmazonCognito_Term_Attributes      `gorm:"ForeignKey:AmazonCognito_TermID"`
}

type AmazonCognito_Term_Attributes struct {
	gorm.Model
	AmazonCognito_TermID uint
	Key                  string
	Value                string
}

type AmazonCognito_Term_PriceDimensions struct {
	gorm.Model
	AmazonCognito_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AmazonCognito_Term_PricePerUnit `gorm:"ForeignKey:AmazonCognito_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonCognito_Term_PricePerUnit struct {
	gorm.Model
	AmazonCognito_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AmazonCognito) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCognito/current/index.json"
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
