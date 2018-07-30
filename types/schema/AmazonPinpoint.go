package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonPinpoint struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonPinpoint_Product
	Terms           map[string]map[string]map[string]rawAmazonPinpoint_Term
}

type rawAmazonPinpoint_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonPinpoint_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonPinpoint) UnmarshalJSON(data []byte) error {
	var p rawAmazonPinpoint
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonPinpoint_Product{}
	terms := []*AmazonPinpoint_Term{}

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
				pDimensions := []*AmazonPinpoint_Term_PriceDimensions{}
				tAttributes := []*AmazonPinpoint_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonPinpoint_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonPinpoint_Term{
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

type AmazonPinpoint struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonPinpoint_Product `gorm:"ForeignKey:AmazonPinpointID"`
	Terms           []*AmazonPinpoint_Term    `gorm:"ForeignKey:AmazonPinpointID"`
}
type AmazonPinpoint_Product struct {
	gorm.Model
	AmazonPinpointID uint
	ProductFamily    string
	Attributes       AmazonPinpoint_Product_Attributes `gorm:"ForeignKey:AmazonPinpoint_Product_AttributesID"`
	Sku              string
}
type AmazonPinpoint_Product_Attributes struct {
	gorm.Model
	AmazonPinpoint_Product_AttributesID uint
	Servicecode                         string
	LocationType                        string
	Group                               string
	Usagetype                           string
	MeteringType                        string
	Servicename                         string
	Location                            string
	GroupDescription                    string
	Operation                           string
}

type AmazonPinpoint_Term struct {
	gorm.Model
	OfferTermCode    string
	AmazonPinpointID uint
	Sku              string
	EffectiveDate    string
	PriceDimensions  []*AmazonPinpoint_Term_PriceDimensions `gorm:"ForeignKey:AmazonPinpoint_TermID"`
	TermAttributes   []*AmazonPinpoint_Term_Attributes      `gorm:"ForeignKey:AmazonPinpoint_TermID"`
}

type AmazonPinpoint_Term_Attributes struct {
	gorm.Model
	AmazonPinpoint_TermID uint
	Key                   string
	Value                 string
}

type AmazonPinpoint_Term_PriceDimensions struct {
	gorm.Model
	AmazonPinpoint_TermID uint
	RateCode              string
	RateType              string
	Description           string
	BeginRange            string
	EndRange              string
	Unit                  string
	PricePerUnit          *AmazonPinpoint_Term_PricePerUnit `gorm:"ForeignKey:AmazonPinpoint_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonPinpoint_Term_PricePerUnit struct {
	gorm.Model
	AmazonPinpoint_Term_PriceDimensionsID uint
	USD                                   string
}

func (a *AmazonPinpoint) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonPinpoint/current/index.json"
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
