package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSEvents struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSEvents_Product
	Terms           map[string]map[string]map[string]rawAWSEvents_Term
}

type rawAWSEvents_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSEvents_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSEvents) UnmarshalJSON(data []byte) error {
	var p rawAWSEvents
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSEvents_Product{}
	terms := []*AWSEvents_Term{}

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
				pDimensions := []*AWSEvents_Term_PriceDimensions{}
				tAttributes := []*AWSEvents_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSEvents_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSEvents_Term{
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

type AWSEvents struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSEvents_Product `gorm:"ForeignKey:AWSEventsID"`
	Terms           []*AWSEvents_Term    `gorm:"ForeignKey:AWSEventsID"`
}
type AWSEvents_Product struct {
	gorm.Model
	AWSEventsID   uint
	Sku           string
	ProductFamily string
	Attributes    AWSEvents_Product_Attributes `gorm:"ForeignKey:AWSEvents_Product_AttributesID"`
}
type AWSEvents_Product_Attributes struct {
	gorm.Model
	AWSEvents_Product_AttributesID uint
	LocationType                   string
	Usagetype                      string
	Operation                      string
	EventType                      string
	Servicename                    string
	Servicecode                    string
	Location                       string
}

type AWSEvents_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSEventsID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSEvents_Term_PriceDimensions `gorm:"ForeignKey:AWSEvents_TermID"`
	TermAttributes  []*AWSEvents_Term_Attributes      `gorm:"ForeignKey:AWSEvents_TermID"`
}

type AWSEvents_Term_Attributes struct {
	gorm.Model
	AWSEvents_TermID uint
	Key              string
	Value            string
}

type AWSEvents_Term_PriceDimensions struct {
	gorm.Model
	AWSEvents_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AWSEvents_Term_PricePerUnit `gorm:"ForeignKey:AWSEvents_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSEvents_Term_PricePerUnit struct {
	gorm.Model
	AWSEvents_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AWSEvents) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSEvents/current/index.json"
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
