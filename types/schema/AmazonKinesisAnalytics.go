package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonKinesisAnalytics struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonKinesisAnalytics_Product
	Terms           map[string]map[string]map[string]rawAmazonKinesisAnalytics_Term
}

type rawAmazonKinesisAnalytics_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonKinesisAnalytics_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonKinesisAnalytics) UnmarshalJSON(data []byte) error {
	var p rawAmazonKinesisAnalytics
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonKinesisAnalytics_Product{}
	terms := []*AmazonKinesisAnalytics_Term{}

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
				pDimensions := []*AmazonKinesisAnalytics_Term_PriceDimensions{}
				tAttributes := []*AmazonKinesisAnalytics_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonKinesisAnalytics_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonKinesisAnalytics_Term{
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

type AmazonKinesisAnalytics struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonKinesisAnalytics_Product `gorm:"ForeignKey:AmazonKinesisAnalyticsID"`
	Terms           []*AmazonKinesisAnalytics_Term    `gorm:"ForeignKey:AmazonKinesisAnalyticsID"`
}
type AmazonKinesisAnalytics_Product struct {
	gorm.Model
	AmazonKinesisAnalyticsID uint
	Attributes               AmazonKinesisAnalytics_Product_Attributes `gorm:"ForeignKey:AmazonKinesisAnalytics_Product_AttributesID"`
	Sku                      string
	ProductFamily            string
}
type AmazonKinesisAnalytics_Product_Attributes struct {
	gorm.Model
	AmazonKinesisAnalytics_Product_AttributesID uint
	Usagetype                                   string
	Operation                                   string
	Servicecode                                 string
	Description                                 string
	Location                                    string
	LocationType                                string
}

type AmazonKinesisAnalytics_Term struct {
	gorm.Model
	OfferTermCode            string
	AmazonKinesisAnalyticsID uint
	Sku                      string
	EffectiveDate            string
	PriceDimensions          []*AmazonKinesisAnalytics_Term_PriceDimensions `gorm:"ForeignKey:AmazonKinesisAnalytics_TermID"`
	TermAttributes           []*AmazonKinesisAnalytics_Term_Attributes      `gorm:"ForeignKey:AmazonKinesisAnalytics_TermID"`
}

type AmazonKinesisAnalytics_Term_Attributes struct {
	gorm.Model
	AmazonKinesisAnalytics_TermID uint
	Key                           string
	Value                         string
}

type AmazonKinesisAnalytics_Term_PriceDimensions struct {
	gorm.Model
	AmazonKinesisAnalytics_TermID uint
	RateCode                      string
	RateType                      string
	Description                   string
	BeginRange                    string
	EndRange                      string
	Unit                          string
	PricePerUnit                  *AmazonKinesisAnalytics_Term_PricePerUnit `gorm:"ForeignKey:AmazonKinesisAnalytics_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonKinesisAnalytics_Term_PricePerUnit struct {
	gorm.Model
	AmazonKinesisAnalytics_Term_PriceDimensionsID uint
	USD                                           string
}

func (a *AmazonKinesisAnalytics) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesisAnalytics/current/index.json"
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
