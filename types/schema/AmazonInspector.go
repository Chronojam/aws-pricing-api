package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonInspector struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonInspector_Product
	Terms           map[string]map[string]map[string]rawAmazonInspector_Term
}

type rawAmazonInspector_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonInspector_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonInspector) UnmarshalJSON(data []byte) error {
	var p rawAmazonInspector
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonInspector_Product{}
	terms := []*AmazonInspector_Term{}

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
				pDimensions := []*AmazonInspector_Term_PriceDimensions{}
				tAttributes := []*AmazonInspector_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonInspector_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonInspector_Term{
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

type AmazonInspector struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonInspector_Product `gorm:"ForeignKey:AmazonInspectorID"`
	Terms           []*AmazonInspector_Term    `gorm:"ForeignKey:AmazonInspectorID"`
}
type AmazonInspector_Product struct {
	gorm.Model
	AmazonInspectorID uint
	Sku               string
	ProductFamily     string
	Attributes        AmazonInspector_Product_Attributes `gorm:"ForeignKey:AmazonInspector_Product_AttributesID"`
}
type AmazonInspector_Product_Attributes struct {
	gorm.Model
	AmazonInspector_Product_AttributesID uint
	Operation                            string
	FreeUsageIncluded                    string
	Servicename                          string
	Servicecode                          string
	Description                          string
	Location                             string
	LocationType                         string
	Usagetype                            string
}

type AmazonInspector_Term struct {
	gorm.Model
	OfferTermCode     string
	AmazonInspectorID uint
	Sku               string
	EffectiveDate     string
	PriceDimensions   []*AmazonInspector_Term_PriceDimensions `gorm:"ForeignKey:AmazonInspector_TermID"`
	TermAttributes    []*AmazonInspector_Term_Attributes      `gorm:"ForeignKey:AmazonInspector_TermID"`
}

type AmazonInspector_Term_Attributes struct {
	gorm.Model
	AmazonInspector_TermID uint
	Key                    string
	Value                  string
}

type AmazonInspector_Term_PriceDimensions struct {
	gorm.Model
	AmazonInspector_TermID uint
	RateCode               string
	RateType               string
	Description            string
	BeginRange             string
	EndRange               string
	Unit                   string
	PricePerUnit           *AmazonInspector_Term_PricePerUnit `gorm:"ForeignKey:AmazonInspector_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonInspector_Term_PricePerUnit struct {
	gorm.Model
	AmazonInspector_Term_PriceDimensionsID uint
	USD                                    string
}

func (a *AmazonInspector) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonInspector/current/index.json"
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
