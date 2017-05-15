package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonML struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonML_Product
	Terms           map[string]map[string]map[string]rawAmazonML_Term
}

type rawAmazonML_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonML_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonML) UnmarshalJSON(data []byte) error {
	var p rawAmazonML
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonML_Product{}
	terms := []*AmazonML_Term{}

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
				pDimensions := []*AmazonML_Term_PriceDimensions{}
				tAttributes := []*AmazonML_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonML_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonML_Term{
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

type AmazonML struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonML_Product `gorm:"ForeignKey:AmazonMLID"`
	Terms           []*AmazonML_Term    `gorm:"ForeignKey:AmazonMLID"`
}
type AmazonML_Product struct {
	gorm.Model
	AmazonMLID    uint
	Sku           string
	ProductFamily string
	Attributes    AmazonML_Product_Attributes `gorm:"ForeignKey:AmazonML_Product_AttributesID"`
}
type AmazonML_Product_Attributes struct {
	gorm.Model
	AmazonML_Product_AttributesID uint
	MachineLearningProcess        string
	Servicecode                   string
	Location                      string
	LocationType                  string
	Group                         string
	GroupDescription              string
	Usagetype                     string
	Operation                     string
}

type AmazonML_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonMLID      uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonML_Term_PriceDimensions `gorm:"ForeignKey:AmazonML_TermID"`
	TermAttributes  []*AmazonML_Term_Attributes      `gorm:"ForeignKey:AmazonML_TermID"`
}

type AmazonML_Term_Attributes struct {
	gorm.Model
	AmazonML_TermID uint
	Key             string
	Value           string
}

type AmazonML_Term_PriceDimensions struct {
	gorm.Model
	AmazonML_TermID uint
	RateCode        string
	RateType        string
	Description     string
	BeginRange      string
	EndRange        string
	Unit            string
	PricePerUnit    *AmazonML_Term_PricePerUnit `gorm:"ForeignKey:AmazonML_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonML_Term_PricePerUnit struct {
	gorm.Model
	AmazonML_Term_PriceDimensionsID uint
	USD                             string
}

func (a *AmazonML) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonML/current/index.json"
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
