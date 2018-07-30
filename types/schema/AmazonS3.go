package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonS3 struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonS3_Product
	Terms           map[string]map[string]map[string]rawAmazonS3_Term
}

type rawAmazonS3_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonS3_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonS3) UnmarshalJSON(data []byte) error {
	var p rawAmazonS3
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonS3_Product{}
	terms := []*AmazonS3_Term{}

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
				pDimensions := []*AmazonS3_Term_PriceDimensions{}
				tAttributes := []*AmazonS3_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonS3_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonS3_Term{
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

type AmazonS3 struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonS3_Product `gorm:"ForeignKey:AmazonS3ID"`
	Terms           []*AmazonS3_Term    `gorm:"ForeignKey:AmazonS3ID"`
}
type AmazonS3_Product struct {
	gorm.Model
	AmazonS3ID    uint
	Sku           string
	ProductFamily string
	Attributes    AmazonS3_Product_Attributes `gorm:"ForeignKey:AmazonS3_Product_AttributesID"`
}
type AmazonS3_Product_Attributes struct {
	gorm.Model
	AmazonS3_Product_AttributesID uint
	TransferType                  string
	FromLocationType              string
	ToLocationType                string
	Usagetype                     string
	Servicecode                   string
	FromLocation                  string
	ToLocation                    string
	Operation                     string
	Servicename                   string
}

type AmazonS3_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonS3ID      uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonS3_Term_PriceDimensions `gorm:"ForeignKey:AmazonS3_TermID"`
	TermAttributes  []*AmazonS3_Term_Attributes      `gorm:"ForeignKey:AmazonS3_TermID"`
}

type AmazonS3_Term_Attributes struct {
	gorm.Model
	AmazonS3_TermID uint
	Key             string
	Value           string
}

type AmazonS3_Term_PriceDimensions struct {
	gorm.Model
	AmazonS3_TermID uint
	RateCode        string
	RateType        string
	Description     string
	BeginRange      string
	EndRange        string
	Unit            string
	PricePerUnit    *AmazonS3_Term_PricePerUnit `gorm:"ForeignKey:AmazonS3_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonS3_Term_PricePerUnit struct {
	gorm.Model
	AmazonS3_Term_PriceDimensionsID uint
	USD                             string
}

func (a *AmazonS3) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonS3/current/index.json"
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
