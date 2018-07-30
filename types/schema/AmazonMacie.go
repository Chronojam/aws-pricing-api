package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonMacie struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonMacie_Product
	Terms           map[string]map[string]map[string]rawAmazonMacie_Term
}

type rawAmazonMacie_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonMacie_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonMacie) UnmarshalJSON(data []byte) error {
	var p rawAmazonMacie
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonMacie_Product{}
	terms := []*AmazonMacie_Term{}

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
				pDimensions := []*AmazonMacie_Term_PriceDimensions{}
				tAttributes := []*AmazonMacie_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonMacie_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonMacie_Term{
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

type AmazonMacie struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonMacie_Product `gorm:"ForeignKey:AmazonMacieID"`
	Terms           []*AmazonMacie_Term    `gorm:"ForeignKey:AmazonMacieID"`
}
type AmazonMacie_Product struct {
	gorm.Model
	AmazonMacieID uint
	Sku           string
	ProductFamily string
	Attributes    AmazonMacie_Product_Attributes `gorm:"ForeignKey:AmazonMacie_Product_AttributesID"`
}
type AmazonMacie_Product_Attributes struct {
	gorm.Model
	AmazonMacie_Product_AttributesID uint
	Location                         string
	Usagetype                        string
	LogsSource                       string
	Servicename                      string
	Servicecode                      string
	LocationType                     string
	Operation                        string
	ActivityType                     string
	LogsType                         string
}

type AmazonMacie_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonMacieID   uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonMacie_Term_PriceDimensions `gorm:"ForeignKey:AmazonMacie_TermID"`
	TermAttributes  []*AmazonMacie_Term_Attributes      `gorm:"ForeignKey:AmazonMacie_TermID"`
}

type AmazonMacie_Term_Attributes struct {
	gorm.Model
	AmazonMacie_TermID uint
	Key                string
	Value              string
}

type AmazonMacie_Term_PriceDimensions struct {
	gorm.Model
	AmazonMacie_TermID uint
	RateCode           string
	RateType           string
	Description        string
	BeginRange         string
	EndRange           string
	Unit               string
	PricePerUnit       *AmazonMacie_Term_PricePerUnit `gorm:"ForeignKey:AmazonMacie_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonMacie_Term_PricePerUnit struct {
	gorm.Model
	AmazonMacie_Term_PriceDimensionsID uint
	USD                                string
}

func (a *AmazonMacie) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonMacie/current/index.json"
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
