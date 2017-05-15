package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonCloudFront struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonCloudFront_Product
	Terms           map[string]map[string]map[string]rawAmazonCloudFront_Term
}

type rawAmazonCloudFront_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonCloudFront_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonCloudFront) UnmarshalJSON(data []byte) error {
	var p rawAmazonCloudFront
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonCloudFront_Product{}
	terms := []*AmazonCloudFront_Term{}

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
				pDimensions := []*AmazonCloudFront_Term_PriceDimensions{}
				tAttributes := []*AmazonCloudFront_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonCloudFront_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonCloudFront_Term{
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

type AmazonCloudFront struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonCloudFront_Product `gorm:"ForeignKey:AmazonCloudFrontID"`
	Terms           []*AmazonCloudFront_Term    `gorm:"ForeignKey:AmazonCloudFrontID"`
}
type AmazonCloudFront_Product struct {
	gorm.Model
	AmazonCloudFrontID uint
	Sku                string
	ProductFamily      string
	Attributes         AmazonCloudFront_Product_Attributes `gorm:"ForeignKey:AmazonCloudFront_Product_AttributesID"`
}
type AmazonCloudFront_Product_Attributes struct {
	gorm.Model
	AmazonCloudFront_Product_AttributesID uint
	RequestDescription                    string
	RequestType                           string
	Servicecode                           string
	Location                              string
	LocationType                          string
	Usagetype                             string
	Operation                             string
}

type AmazonCloudFront_Term struct {
	gorm.Model
	OfferTermCode      string
	AmazonCloudFrontID uint
	Sku                string
	EffectiveDate      string
	PriceDimensions    []*AmazonCloudFront_Term_PriceDimensions `gorm:"ForeignKey:AmazonCloudFront_TermID"`
	TermAttributes     []*AmazonCloudFront_Term_Attributes      `gorm:"ForeignKey:AmazonCloudFront_TermID"`
}

type AmazonCloudFront_Term_Attributes struct {
	gorm.Model
	AmazonCloudFront_TermID uint
	Key                     string
	Value                   string
}

type AmazonCloudFront_Term_PriceDimensions struct {
	gorm.Model
	AmazonCloudFront_TermID uint
	RateCode                string
	RateType                string
	Description             string
	BeginRange              string
	EndRange                string
	Unit                    string
	PricePerUnit            *AmazonCloudFront_Term_PricePerUnit `gorm:"ForeignKey:AmazonCloudFront_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonCloudFront_Term_PricePerUnit struct {
	gorm.Model
	AmazonCloudFront_Term_PriceDimensionsID uint
	USD                                     string
}

func (a *AmazonCloudFront) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudFront/current/index.json"
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
