package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSDatabaseMigrationSvc struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSDatabaseMigrationSvc_Product
	Terms           map[string]map[string]map[string]rawAWSDatabaseMigrationSvc_Term
}

type rawAWSDatabaseMigrationSvc_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSDatabaseMigrationSvc_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSDatabaseMigrationSvc) UnmarshalJSON(data []byte) error {
	var p rawAWSDatabaseMigrationSvc
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSDatabaseMigrationSvc_Product{}
	terms := []*AWSDatabaseMigrationSvc_Term{}

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
				pDimensions := []*AWSDatabaseMigrationSvc_Term_PriceDimensions{}
				tAttributes := []*AWSDatabaseMigrationSvc_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSDatabaseMigrationSvc_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSDatabaseMigrationSvc_Term{
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

type AWSDatabaseMigrationSvc struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSDatabaseMigrationSvc_Product `gorm:"ForeignKey:AWSDatabaseMigrationSvcID"`
	Terms           []*AWSDatabaseMigrationSvc_Term    `gorm:"ForeignKey:AWSDatabaseMigrationSvcID"`
}
type AWSDatabaseMigrationSvc_Product struct {
	gorm.Model
	AWSDatabaseMigrationSvcID uint
	Sku                       string
	ProductFamily             string
	Attributes                AWSDatabaseMigrationSvc_Product_Attributes `gorm:"ForeignKey:AWSDatabaseMigrationSvc_Product_AttributesID"`
}
type AWSDatabaseMigrationSvc_Product_Attributes struct {
	gorm.Model
	AWSDatabaseMigrationSvc_Product_AttributesID uint
	ToLocation                                   string
	ToLocationType                               string
	Usagetype                                    string
	Operation                                    string
	Servicecode                                  string
	TransferType                                 string
	FromLocation                                 string
	FromLocationType                             string
}

type AWSDatabaseMigrationSvc_Term struct {
	gorm.Model
	OfferTermCode             string
	AWSDatabaseMigrationSvcID uint
	Sku                       string
	EffectiveDate             string
	PriceDimensions           []*AWSDatabaseMigrationSvc_Term_PriceDimensions `gorm:"ForeignKey:AWSDatabaseMigrationSvc_TermID"`
	TermAttributes            []*AWSDatabaseMigrationSvc_Term_Attributes      `gorm:"ForeignKey:AWSDatabaseMigrationSvc_TermID"`
}

type AWSDatabaseMigrationSvc_Term_Attributes struct {
	gorm.Model
	AWSDatabaseMigrationSvc_TermID uint
	Key                            string
	Value                          string
}

type AWSDatabaseMigrationSvc_Term_PriceDimensions struct {
	gorm.Model
	AWSDatabaseMigrationSvc_TermID uint
	RateCode                       string
	RateType                       string
	Description                    string
	BeginRange                     string
	EndRange                       string
	Unit                           string
	PricePerUnit                   *AWSDatabaseMigrationSvc_Term_PricePerUnit `gorm:"ForeignKey:AWSDatabaseMigrationSvc_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSDatabaseMigrationSvc_Term_PricePerUnit struct {
	gorm.Model
	AWSDatabaseMigrationSvc_Term_PriceDimensionsID uint
	USD                                            string
}

func (a *AWSDatabaseMigrationSvc) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDatabaseMigrationSvc/current/index.json"
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
