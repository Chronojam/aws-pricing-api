package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonRedshift struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonRedshift_Product
	Terms           map[string]map[string]map[string]rawAmazonRedshift_Term
}

type rawAmazonRedshift_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonRedshift_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonRedshift) UnmarshalJSON(data []byte) error {
	var p rawAmazonRedshift
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonRedshift_Product{}
	terms := []*AmazonRedshift_Term{}

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
				pDimensions := []*AmazonRedshift_Term_PriceDimensions{}
				tAttributes := []*AmazonRedshift_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonRedshift_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonRedshift_Term{
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

type AmazonRedshift struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonRedshift_Product `gorm:"ForeignKey:AmazonRedshiftID"`
	Terms           []*AmazonRedshift_Term    `gorm:"ForeignKey:AmazonRedshiftID"`
}
type AmazonRedshift_Product struct {
	gorm.Model
	AmazonRedshiftID uint
	Sku              string
	ProductFamily    string
	Attributes       AmazonRedshift_Product_Attributes `gorm:"ForeignKey:AmazonRedshift_Product_AttributesID"`
}
type AmazonRedshift_Product_Attributes struct {
	gorm.Model
	AmazonRedshift_Product_AttributesID uint
	Usagetype                           string
	Operation                           string
	Servicecode                         string
	TransferType                        string
	FromLocation                        string
	FromLocationType                    string
	ToLocation                          string
	ToLocationType                      string
}

type AmazonRedshift_Term struct {
	gorm.Model
	OfferTermCode    string
	AmazonRedshiftID uint
	Sku              string
	EffectiveDate    string
	PriceDimensions  []*AmazonRedshift_Term_PriceDimensions `gorm:"ForeignKey:AmazonRedshift_TermID"`
	TermAttributes   []*AmazonRedshift_Term_Attributes      `gorm:"ForeignKey:AmazonRedshift_TermID"`
}

type AmazonRedshift_Term_Attributes struct {
	gorm.Model
	AmazonRedshift_TermID uint
	Key                   string
	Value                 string
}

type AmazonRedshift_Term_PriceDimensions struct {
	gorm.Model
	AmazonRedshift_TermID uint
	RateCode              string
	RateType              string
	Description           string
	BeginRange            string
	EndRange              string
	Unit                  string
	PricePerUnit          *AmazonRedshift_Term_PricePerUnit `gorm:"ForeignKey:AmazonRedshift_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonRedshift_Term_PricePerUnit struct {
	gorm.Model
	AmazonRedshift_Term_PriceDimensionsID uint
	USD                                   string
}

func (a *AmazonRedshift) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRedshift/current/index.json"
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
