package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonKinesis struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonKinesis_Product
	Terms           map[string]map[string]map[string]rawAmazonKinesis_Term
}

type rawAmazonKinesis_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonKinesis_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonKinesis) UnmarshalJSON(data []byte) error {
	var p rawAmazonKinesis
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonKinesis_Product{}
	terms := []*AmazonKinesis_Term{}

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
				pDimensions := []*AmazonKinesis_Term_PriceDimensions{}
				tAttributes := []*AmazonKinesis_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonKinesis_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonKinesis_Term{
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

type AmazonKinesis struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonKinesis_Product `gorm:"ForeignKey:AmazonKinesisID"`
	Terms           []*AmazonKinesis_Term    `gorm:"ForeignKey:AmazonKinesisID"`
}
type AmazonKinesis_Product struct {
	gorm.Model
	AmazonKinesisID uint
	Sku             string
	ProductFamily   string
	Attributes      AmazonKinesis_Product_Attributes `gorm:"ForeignKey:AmazonKinesis_Product_AttributesID"`
}
type AmazonKinesis_Product_Attributes struct {
	gorm.Model
	AmazonKinesis_Product_AttributesID uint
	Group                              string
	GroupDescription                   string
	Usagetype                          string
	Operation                          string
	Servicename                        string
	Servicecode                        string
	Location                           string
	LocationType                       string
}

type AmazonKinesis_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonKinesisID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonKinesis_Term_PriceDimensions `gorm:"ForeignKey:AmazonKinesis_TermID"`
	TermAttributes  []*AmazonKinesis_Term_Attributes      `gorm:"ForeignKey:AmazonKinesis_TermID"`
}

type AmazonKinesis_Term_Attributes struct {
	gorm.Model
	AmazonKinesis_TermID uint
	Key                  string
	Value                string
}

type AmazonKinesis_Term_PriceDimensions struct {
	gorm.Model
	AmazonKinesis_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AmazonKinesis_Term_PricePerUnit `gorm:"ForeignKey:AmazonKinesis_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonKinesis_Term_PricePerUnit struct {
	gorm.Model
	AmazonKinesis_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AmazonKinesis) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesis/current/index.json"
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
