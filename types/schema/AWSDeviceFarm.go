package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSDeviceFarm struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSDeviceFarm_Product
	Terms           map[string]map[string]map[string]rawAWSDeviceFarm_Term
}

type rawAWSDeviceFarm_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSDeviceFarm_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSDeviceFarm) UnmarshalJSON(data []byte) error {
	var p rawAWSDeviceFarm
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSDeviceFarm_Product{}
	terms := []*AWSDeviceFarm_Term{}

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
				pDimensions := []*AWSDeviceFarm_Term_PriceDimensions{}
				tAttributes := []*AWSDeviceFarm_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSDeviceFarm_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSDeviceFarm_Term{
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

type AWSDeviceFarm struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSDeviceFarm_Product `gorm:"ForeignKey:AWSDeviceFarmID"`
	Terms           []*AWSDeviceFarm_Term    `gorm:"ForeignKey:AWSDeviceFarmID"`
}
type AWSDeviceFarm_Product struct {
	gorm.Model
	AWSDeviceFarmID uint
	Sku             string
	ProductFamily   string
	Attributes      AWSDeviceFarm_Product_Attributes `gorm:"ForeignKey:AWSDeviceFarm_Product_AttributesID"`
}
type AWSDeviceFarm_Product_Attributes struct {
	gorm.Model
	AWSDeviceFarm_Product_AttributesID uint
	Operation                          string
	Servicecode                        string
	Location                           string
	LocationType                       string
	Usagetype                          string
	DeviceOs                           string
	ExecutionMode                      string
	MeterMode                          string
	Description                        string
}

type AWSDeviceFarm_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSDeviceFarmID uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSDeviceFarm_Term_PriceDimensions `gorm:"ForeignKey:AWSDeviceFarm_TermID"`
	TermAttributes  []*AWSDeviceFarm_Term_Attributes      `gorm:"ForeignKey:AWSDeviceFarm_TermID"`
}

type AWSDeviceFarm_Term_Attributes struct {
	gorm.Model
	AWSDeviceFarm_TermID uint
	Key                  string
	Value                string
}

type AWSDeviceFarm_Term_PriceDimensions struct {
	gorm.Model
	AWSDeviceFarm_TermID uint
	RateCode             string
	RateType             string
	Description          string
	BeginRange           string
	EndRange             string
	Unit                 string
	PricePerUnit         *AWSDeviceFarm_Term_PricePerUnit `gorm:"ForeignKey:AWSDeviceFarm_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSDeviceFarm_Term_PricePerUnit struct {
	gorm.Model
	AWSDeviceFarm_Term_PriceDimensionsID uint
	USD                                  string
}

func (a *AWSDeviceFarm) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDeviceFarm/current/index.json"
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
