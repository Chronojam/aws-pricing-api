package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonNeptune struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonNeptune_Product
	Terms           map[string]map[string]map[string]rawAmazonNeptune_Term
}

type rawAmazonNeptune_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonNeptune_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonNeptune) UnmarshalJSON(data []byte) error {
	var p rawAmazonNeptune
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonNeptune_Product{}
	terms := []*AmazonNeptune_Term{}

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
				pDimensions := []*AmazonNeptune_Term_PriceDimensions{}
				tAttributes := []*AmazonNeptune_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonNeptune_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonNeptune_Term{
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

type AmazonNeptune struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonNeptune_Product `gorm:"ForeignKey:AmazonNeptuneID"`
	Terms           []*AmazonNeptune_Term    `gorm:"ForeignKey:AmazonNeptuneID"`
}
type AmazonNeptune_Product struct {
	gorm.Model
	AmazonNeptuneID uint
	Sku             string
	ProductFamily   string
	Attributes      AmazonNeptune_Product_Attributes `gorm:"ForeignKey:AmazonNeptune_Product_AttributesID"`
}
type AmazonNeptune_Product_Attributes struct {
	gorm.Model
	AmazonNeptune_Product_AttributesID uint
	Operation                          string
	Servicecode                        string
	TransferType                       string
	FromLocation                       string
	FromLocationType                   string
	ToLocation                         string
	ToLocationType                     string
	Usagetype                          string
	Servicename                        string
}

type AmazonNeptune_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonNeptuneID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonNeptune_Term_PriceDimensions `gorm:"ForeignKey:AmazonNeptune_TermID"`
	TermAttributes  []*AmazonNeptune_Term_Attributes      `gorm:"ForeignKey:AmazonNeptune_TermID"`
}

type AmazonNeptune_Term_Attributes struct {
	gorm.Model
	AmazonNeptune_TermID uint
	Key                  string
	Value                string
}

type AmazonNeptune_Term_PriceDimensions struct {
	gorm.Model
	AmazonNeptune_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AmazonNeptune_Term_PricePerUnit `gorm:"ForeignKey:AmazonNeptune_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonNeptune_Term_PricePerUnit struct {
	gorm.Model
	AmazonNeptune_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AmazonNeptune) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonNeptune/current/index.json"
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
