package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonCloudWatch struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonCloudWatch_Product
	Terms           map[string]map[string]map[string]rawAmazonCloudWatch_Term
}

type rawAmazonCloudWatch_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonCloudWatch_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonCloudWatch) UnmarshalJSON(data []byte) error {
	var p rawAmazonCloudWatch
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonCloudWatch_Product{}
	terms := []*AmazonCloudWatch_Term{}

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
				pDimensions := []*AmazonCloudWatch_Term_PriceDimensions{}
				tAttributes := []*AmazonCloudWatch_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonCloudWatch_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonCloudWatch_Term{
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

type AmazonCloudWatch struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonCloudWatch_Product `gorm:"ForeignKey:AmazonCloudWatchID"`
	Terms           []*AmazonCloudWatch_Term    `gorm:"ForeignKey:AmazonCloudWatchID"`
}
type AmazonCloudWatch_Product struct {
	gorm.Model
	AmazonCloudWatchID uint
	Sku                string
	ProductFamily      string
	Attributes         AmazonCloudWatch_Product_Attributes `gorm:"ForeignKey:AmazonCloudWatch_Product_AttributesID"`
}
type AmazonCloudWatch_Product_Attributes struct {
	gorm.Model
	AmazonCloudWatch_Product_AttributesID uint
	Servicecode                           string
	Location                              string
	LocationType                          string
	Group                                 string
	GroupDescription                      string
	Usagetype                             string
	Operation                             string
}

type AmazonCloudWatch_Term struct {
	gorm.Model
	OfferTermCode      string
	AmazonCloudWatchID uint
	Sku                string
	EffectiveDate      string
	PriceDimensions    []*AmazonCloudWatch_Term_PriceDimensions `gorm:"ForeignKey:AmazonCloudWatch_TermID"`
	TermAttributes     []*AmazonCloudWatch_Term_Attributes      `gorm:"ForeignKey:AmazonCloudWatch_TermID"`
}

type AmazonCloudWatch_Term_Attributes struct {
	gorm.Model
	AmazonCloudWatch_TermID uint
	Key                     string
	Value                   string
}

type AmazonCloudWatch_Term_PriceDimensions struct {
	gorm.Model
	AmazonCloudWatch_TermID uint
	RateCode                string
	RateType                string
	Description             string
	BeginRange              string
	EndRange                string
	Unit                    string
	PricePerUnit            *AmazonCloudWatch_Term_PricePerUnit `gorm:"ForeignKey:AmazonCloudWatch_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonCloudWatch_Term_PricePerUnit struct {
	gorm.Model
	AmazonCloudWatch_Term_PriceDimensionsID uint
	USD                                     string
}

func (a *AmazonCloudWatch) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudWatch/current/index.json"
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
