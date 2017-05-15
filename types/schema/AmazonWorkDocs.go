package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonWorkDocs struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonWorkDocs_Product
	Terms           map[string]map[string]map[string]rawAmazonWorkDocs_Term
}

type rawAmazonWorkDocs_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonWorkDocs_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonWorkDocs) UnmarshalJSON(data []byte) error {
	var p rawAmazonWorkDocs
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonWorkDocs_Product{}
	terms := []*AmazonWorkDocs_Term{}

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
				pDimensions := []*AmazonWorkDocs_Term_PriceDimensions{}
				tAttributes := []*AmazonWorkDocs_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonWorkDocs_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonWorkDocs_Term{
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

type AmazonWorkDocs struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonWorkDocs_Product `gorm:"ForeignKey:AmazonWorkDocsID"`
	Terms           []*AmazonWorkDocs_Term    `gorm:"ForeignKey:AmazonWorkDocsID"`
}
type AmazonWorkDocs_Product struct {
	gorm.Model
	AmazonWorkDocsID uint
	Sku              string
	ProductFamily    string
	Attributes       AmazonWorkDocs_Product_Attributes `gorm:"ForeignKey:AmazonWorkDocs_Product_AttributesID"`
}
type AmazonWorkDocs_Product_Attributes struct {
	gorm.Model
	AmazonWorkDocs_Product_AttributesID uint
	Description                         string
	LocationType                        string
	Storage                             string
	Operation                           string
	FreeTrial                           string
	MaximumStorageVolume                string
	Servicecode                         string
	Location                            string
	Usagetype                           string
	MinimumStorageVolume                string
}

type AmazonWorkDocs_Term struct {
	gorm.Model
	OfferTermCode    string
	AmazonWorkDocsID uint
	Sku              string
	EffectiveDate    string
	PriceDimensions  []*AmazonWorkDocs_Term_PriceDimensions `gorm:"ForeignKey:AmazonWorkDocs_TermID"`
	TermAttributes   []*AmazonWorkDocs_Term_Attributes      `gorm:"ForeignKey:AmazonWorkDocs_TermID"`
}

type AmazonWorkDocs_Term_Attributes struct {
	gorm.Model
	AmazonWorkDocs_TermID uint
	Key                   string
	Value                 string
}

type AmazonWorkDocs_Term_PriceDimensions struct {
	gorm.Model
	AmazonWorkDocs_TermID uint
	RateCode              string
	RateType              string
	Description           string
	BeginRange            string
	EndRange              string
	Unit                  string
	PricePerUnit          *AmazonWorkDocs_Term_PricePerUnit `gorm:"ForeignKey:AmazonWorkDocs_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonWorkDocs_Term_PricePerUnit struct {
	gorm.Model
	AmazonWorkDocs_Term_PriceDimensionsID uint
	USD                                   string
}

func (a *AmazonWorkDocs) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkDocs/current/index.json"
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
