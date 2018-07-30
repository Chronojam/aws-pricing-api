package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonConnect struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonConnect_Product
	Terms           map[string]map[string]map[string]rawAmazonConnect_Term
}

type rawAmazonConnect_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonConnect_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonConnect) UnmarshalJSON(data []byte) error {
	var p rawAmazonConnect
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonConnect_Product{}
	terms := []*AmazonConnect_Term{}

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
				pDimensions := []*AmazonConnect_Term_PriceDimensions{}
				tAttributes := []*AmazonConnect_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonConnect_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonConnect_Term{
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

type AmazonConnect struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonConnect_Product `gorm:"ForeignKey:AmazonConnectID"`
	Terms           []*AmazonConnect_Term    `gorm:"ForeignKey:AmazonConnectID"`
}
type AmazonConnect_Product struct {
	gorm.Model
	AmazonConnectID uint
	Sku             string
	ProductFamily   string
	Attributes      AmazonConnect_Product_Attributes `gorm:"ForeignKey:AmazonConnect_Product_AttributesID"`
}
type AmazonConnect_Product_Attributes struct {
	gorm.Model
	AmazonConnect_Product_AttributesID uint
	Operation                          string
	Servicename                        string
	Servicecode                        string
	Location                           string
	LocationType                       string
	Usagetype                          string
}

type AmazonConnect_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonConnectID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonConnect_Term_PriceDimensions `gorm:"ForeignKey:AmazonConnect_TermID"`
	TermAttributes  []*AmazonConnect_Term_Attributes      `gorm:"ForeignKey:AmazonConnect_TermID"`
}

type AmazonConnect_Term_Attributes struct {
	gorm.Model
	AmazonConnect_TermID uint
	Key                  string
	Value                string
}

type AmazonConnect_Term_PriceDimensions struct {
	gorm.Model
	AmazonConnect_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AmazonConnect_Term_PricePerUnit `gorm:"ForeignKey:AmazonConnect_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonConnect_Term_PricePerUnit struct {
	gorm.Model
	AmazonConnect_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AmazonConnect) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonConnect/current/index.json"
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
