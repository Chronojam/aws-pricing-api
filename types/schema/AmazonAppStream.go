package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonAppStream struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonAppStream_Product
	Terms           map[string]map[string]map[string]rawAmazonAppStream_Term
}

type rawAmazonAppStream_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonAppStream_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonAppStream) UnmarshalJSON(data []byte) error {
	var p rawAmazonAppStream
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonAppStream_Product{}
	terms := []*AmazonAppStream_Term{}

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
				pDimensions := []*AmazonAppStream_Term_PriceDimensions{}
				tAttributes := []*AmazonAppStream_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonAppStream_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonAppStream_Term{
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

type AmazonAppStream struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonAppStream_Product `gorm:"ForeignKey:AmazonAppStreamID"`
	Terms           []*AmazonAppStream_Term    `gorm:"ForeignKey:AmazonAppStreamID"`
}
type AmazonAppStream_Product struct {
	gorm.Model
	AmazonAppStreamID uint
	Sku               string
	ProductFamily     string
	Attributes        AmazonAppStream_Product_Attributes `gorm:"ForeignKey:AmazonAppStream_Product_AttributesID"`
}
type AmazonAppStream_Product_Attributes struct {
	gorm.Model
	AmazonAppStream_Product_AttributesID uint
	Servicecode                          string
	Location                             string
	LocationType                         string
	Usagetype                            string
	Operation                            string
	InstanceFunction                     string
	Servicename                          string
}

type AmazonAppStream_Term struct {
	gorm.Model
	OfferTermCode     string
	AmazonAppStreamID uint
	Sku               string
	EffectiveDate     string
	PriceDimensions   []*AmazonAppStream_Term_PriceDimensions `gorm:"ForeignKey:AmazonAppStream_TermID"`
	TermAttributes    []*AmazonAppStream_Term_Attributes      `gorm:"ForeignKey:AmazonAppStream_TermID"`
}

type AmazonAppStream_Term_Attributes struct {
	gorm.Model
	AmazonAppStream_TermID uint
	Key                    string
	Value                  string
}

type AmazonAppStream_Term_PriceDimensions struct {
	gorm.Model
	AmazonAppStream_TermID uint
	RateCode               string
	RateType               string
	Description            string
	BeginRange             string
	EndRange               string
	Unit                   string
	PricePerUnit           *AmazonAppStream_Term_PricePerUnit `gorm:"ForeignKey:AmazonAppStream_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonAppStream_Term_PricePerUnit struct {
	gorm.Model
	AmazonAppStream_Term_PriceDimensionsID uint
	USD                                    string
}

func (a *AmazonAppStream) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonAppStream/current/index.json"
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
