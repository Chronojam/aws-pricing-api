package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonES struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonES_Product
	Terms           map[string]map[string]map[string]rawAmazonES_Term
}

type rawAmazonES_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonES_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonES) UnmarshalJSON(data []byte) error {
	var p rawAmazonES
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonES_Product{}
	terms := []*AmazonES_Term{}

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
				pDimensions := []*AmazonES_Term_PriceDimensions{}
				tAttributes := []*AmazonES_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonES_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonES_Term{
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

type AmazonES struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonES_Product `gorm:"ForeignKey:AmazonESID"`
	Terms           []*AmazonES_Term    `gorm:"ForeignKey:AmazonESID"`
}
type AmazonES_Product struct {
	gorm.Model
	AmazonESID    uint
	Sku           string
	ProductFamily string
	Attributes    AmazonES_Product_Attributes `gorm:"ForeignKey:AmazonES_Product_AttributesID"`
}
type AmazonES_Product_Attributes struct {
	gorm.Model
	AmazonES_Product_AttributesID uint
	ToLocationType                string
	Usagetype                     string
	Operation                     string
	Servicecode                   string
	TransferType                  string
	FromLocation                  string
	FromLocationType              string
	ToLocation                    string
}

type AmazonES_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonESID      uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonES_Term_PriceDimensions `gorm:"ForeignKey:AmazonES_TermID"`
	TermAttributes  []*AmazonES_Term_Attributes      `gorm:"ForeignKey:AmazonES_TermID"`
}

type AmazonES_Term_Attributes struct {
	gorm.Model
	AmazonES_TermID uint
	Key             string
	Value           string
}

type AmazonES_Term_PriceDimensions struct {
	gorm.Model
	AmazonES_TermID uint
	RateCode        string
	RateType        string
	Description     string
	BeginRange      string
	EndRange        string
	Unit            string
	PricePerUnit    *AmazonES_Term_PricePerUnit `gorm:"ForeignKey:AmazonES_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonES_Term_PricePerUnit struct {
	gorm.Model
	AmazonES_Term_PriceDimensionsID uint
	USD                             string
}

func (a *AmazonES) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonES/current/index.json"
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
