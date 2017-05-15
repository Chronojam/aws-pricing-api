package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSServiceCatalog struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSServiceCatalog_Product
	Terms           map[string]map[string]map[string]rawAWSServiceCatalog_Term
}

type rawAWSServiceCatalog_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSServiceCatalog_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSServiceCatalog) UnmarshalJSON(data []byte) error {
	var p rawAWSServiceCatalog
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSServiceCatalog_Product{}
	terms := []*AWSServiceCatalog_Term{}

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
				pDimensions := []*AWSServiceCatalog_Term_PriceDimensions{}
				tAttributes := []*AWSServiceCatalog_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSServiceCatalog_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSServiceCatalog_Term{
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

type AWSServiceCatalog struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSServiceCatalog_Product `gorm:"ForeignKey:AWSServiceCatalogID"`
	Terms           []*AWSServiceCatalog_Term    `gorm:"ForeignKey:AWSServiceCatalogID"`
}
type AWSServiceCatalog_Product struct {
	gorm.Model
	AWSServiceCatalogID uint
	Sku                 string
	ProductFamily       string
	Attributes          AWSServiceCatalog_Product_Attributes `gorm:"ForeignKey:AWSServiceCatalog_Product_AttributesID"`
}
type AWSServiceCatalog_Product_Attributes struct {
	gorm.Model
	AWSServiceCatalog_Product_AttributesID uint
	Servicecode                            string
	Location                               string
	LocationType                           string
	Usagetype                              string
	Operation                              string
	WithActiveUsers                        string
}

type AWSServiceCatalog_Term struct {
	gorm.Model
	OfferTermCode       string
	AWSServiceCatalogID uint
	Sku                 string
	EffectiveDate       string
	PriceDimensions     []*AWSServiceCatalog_Term_PriceDimensions `gorm:"ForeignKey:AWSServiceCatalog_TermID"`
	TermAttributes      []*AWSServiceCatalog_Term_Attributes      `gorm:"ForeignKey:AWSServiceCatalog_TermID"`
}

type AWSServiceCatalog_Term_Attributes struct {
	gorm.Model
	AWSServiceCatalog_TermID uint
	Key                      string
	Value                    string
}

type AWSServiceCatalog_Term_PriceDimensions struct {
	gorm.Model
	AWSServiceCatalog_TermID uint
	RateCode                 string
	RateType                 string
	Description              string
	BeginRange               string
	EndRange                 string
	Unit                     string
	PricePerUnit             *AWSServiceCatalog_Term_PricePerUnit `gorm:"ForeignKey:AWSServiceCatalog_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSServiceCatalog_Term_PricePerUnit struct {
	gorm.Model
	AWSServiceCatalog_Term_PriceDimensionsID uint
	USD                                      string
}

func (a *AWSServiceCatalog) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSServiceCatalog/current/index.json"
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
