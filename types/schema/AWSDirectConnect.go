package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSDirectConnect struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSDirectConnect_Product
	Terms           map[string]map[string]map[string]rawAWSDirectConnect_Term
}

type rawAWSDirectConnect_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSDirectConnect_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSDirectConnect) UnmarshalJSON(data []byte) error {
	var p rawAWSDirectConnect
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSDirectConnect_Product{}
	terms := []*AWSDirectConnect_Term{}

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
				pDimensions := []*AWSDirectConnect_Term_PriceDimensions{}
				tAttributes := []*AWSDirectConnect_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSDirectConnect_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSDirectConnect_Term{
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

type AWSDirectConnect struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSDirectConnect_Product `gorm:"ForeignKey:AWSDirectConnectID"`
	Terms           []*AWSDirectConnect_Term    `gorm:"ForeignKey:AWSDirectConnectID"`
}
type AWSDirectConnect_Product struct {
	gorm.Model
	AWSDirectConnectID uint
	Sku                string
	ProductFamily      string
	Attributes         AWSDirectConnect_Product_Attributes `gorm:"ForeignKey:AWSDirectConnect_Product_AttributesID"`
}
type AWSDirectConnect_Product_Attributes struct {
	gorm.Model
	AWSDirectConnect_Product_AttributesID uint
	Location                              string
	LocationType                          string
	Usagetype                             string
	Operation                             string
	DirectConnectLocation                 string
	PortSpeed                             string
	Servicecode                           string
}

type AWSDirectConnect_Term struct {
	gorm.Model
	OfferTermCode      string
	AWSDirectConnectID uint
	Sku                string
	EffectiveDate      string
	PriceDimensions    []*AWSDirectConnect_Term_PriceDimensions `gorm:"ForeignKey:AWSDirectConnect_TermID"`
	TermAttributes     []*AWSDirectConnect_Term_Attributes      `gorm:"ForeignKey:AWSDirectConnect_TermID"`
}

type AWSDirectConnect_Term_Attributes struct {
	gorm.Model
	AWSDirectConnect_TermID uint
	Key                     string
	Value                   string
}

type AWSDirectConnect_Term_PriceDimensions struct {
	gorm.Model
	AWSDirectConnect_TermID uint
	RateCode                string
	RateType                string
	Description             string
	BeginRange              string
	EndRange                string
	Unit                    string
	PricePerUnit            *AWSDirectConnect_Term_PricePerUnit `gorm:"ForeignKey:AWSDirectConnect_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSDirectConnect_Term_PricePerUnit struct {
	gorm.Model
	AWSDirectConnect_Term_PriceDimensionsID uint
	USD                                     string
}

func (a *AWSDirectConnect) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDirectConnect/current/index.json"
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
