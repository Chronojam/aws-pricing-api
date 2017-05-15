package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonRekognition struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonRekognition_Product
	Terms           map[string]map[string]map[string]rawAmazonRekognition_Term
}

type rawAmazonRekognition_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonRekognition_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonRekognition) UnmarshalJSON(data []byte) error {
	var p rawAmazonRekognition
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonRekognition_Product{}
	terms := []*AmazonRekognition_Term{}

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
				pDimensions := []*AmazonRekognition_Term_PriceDimensions{}
				tAttributes := []*AmazonRekognition_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonRekognition_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonRekognition_Term{
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

type AmazonRekognition struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonRekognition_Product `gorm:"ForeignKey:AmazonRekognitionID"`
	Terms           []*AmazonRekognition_Term    `gorm:"ForeignKey:AmazonRekognitionID"`
}
type AmazonRekognition_Product struct {
	gorm.Model
	AmazonRekognitionID uint
	Sku                 string
	ProductFamily       string
	Attributes          AmazonRekognition_Product_Attributes `gorm:"ForeignKey:AmazonRekognition_Product_AttributesID"`
}
type AmazonRekognition_Product_Attributes struct {
	gorm.Model
	AmazonRekognition_Product_AttributesID uint
	Location                               string
	LocationType                           string
	Group                                  string
	GroupDescription                       string
	Usagetype                              string
	Operation                              string
	Servicecode                            string
}

type AmazonRekognition_Term struct {
	gorm.Model
	OfferTermCode       string
	AmazonRekognitionID uint
	Sku                 string
	EffectiveDate       string
	PriceDimensions     []*AmazonRekognition_Term_PriceDimensions `gorm:"ForeignKey:AmazonRekognition_TermID"`
	TermAttributes      []*AmazonRekognition_Term_Attributes      `gorm:"ForeignKey:AmazonRekognition_TermID"`
}

type AmazonRekognition_Term_Attributes struct {
	gorm.Model
	AmazonRekognition_TermID uint
	Key                      string
	Value                    string
}

type AmazonRekognition_Term_PriceDimensions struct {
	gorm.Model
	AmazonRekognition_TermID uint
	RateCode                 string
	RateType                 string
	Description              string
	BeginRange               string
	EndRange                 string
	Unit                     string
	PricePerUnit             *AmazonRekognition_Term_PricePerUnit `gorm:"ForeignKey:AmazonRekognition_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonRekognition_Term_PricePerUnit struct {
	gorm.Model
	AmazonRekognition_Term_PriceDimensionsID uint
	USD                                      string
}

func (a *AmazonRekognition) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRekognition/current/index.json"
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
