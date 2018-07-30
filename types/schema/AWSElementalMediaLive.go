package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSElementalMediaLive struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSElementalMediaLive_Product
	Terms           map[string]map[string]map[string]rawAWSElementalMediaLive_Term
}

type rawAWSElementalMediaLive_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSElementalMediaLive_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSElementalMediaLive) UnmarshalJSON(data []byte) error {
	var p rawAWSElementalMediaLive
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSElementalMediaLive_Product{}
	terms := []*AWSElementalMediaLive_Term{}

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
				pDimensions := []*AWSElementalMediaLive_Term_PriceDimensions{}
				tAttributes := []*AWSElementalMediaLive_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSElementalMediaLive_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSElementalMediaLive_Term{
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

type AWSElementalMediaLive struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSElementalMediaLive_Product `gorm:"ForeignKey:AWSElementalMediaLiveID"`
	Terms           []*AWSElementalMediaLive_Term    `gorm:"ForeignKey:AWSElementalMediaLiveID"`
}
type AWSElementalMediaLive_Product struct {
	gorm.Model
	AWSElementalMediaLiveID uint
	ProductFamily           string
	Attributes              AWSElementalMediaLive_Product_Attributes `gorm:"ForeignKey:AWSElementalMediaLive_Product_AttributesID"`
	Sku                     string
}
type AWSElementalMediaLive_Product_Attributes struct {
	gorm.Model
	AWSElementalMediaLive_Product_AttributesID uint
	TransferType                               string
	Usagetype                                  string
	Operation                                  string
	Servicecode                                string
	FromLocation                               string
	FromLocationType                           string
	ToLocation                                 string
	ToLocationType                             string
	Servicename                                string
}

type AWSElementalMediaLive_Term struct {
	gorm.Model
	OfferTermCode           string
	AWSElementalMediaLiveID uint
	Sku                     string
	EffectiveDate           string
	PriceDimensions         []*AWSElementalMediaLive_Term_PriceDimensions `gorm:"ForeignKey:AWSElementalMediaLive_TermID"`
	TermAttributes          []*AWSElementalMediaLive_Term_Attributes      `gorm:"ForeignKey:AWSElementalMediaLive_TermID"`
}

type AWSElementalMediaLive_Term_Attributes struct {
	gorm.Model
	AWSElementalMediaLive_TermID uint
	Key                          string
	Value                        string
}

type AWSElementalMediaLive_Term_PriceDimensions struct {
	gorm.Model
	AWSElementalMediaLive_TermID uint
	RateCode                     string
	RateType                     string
	Description                  string
	BeginRange                   string
	EndRange                     string
	Unit                         string
	PricePerUnit                 *AWSElementalMediaLive_Term_PricePerUnit `gorm:"ForeignKey:AWSElementalMediaLive_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSElementalMediaLive_Term_PricePerUnit struct {
	gorm.Model
	AWSElementalMediaLive_Term_PriceDimensionsID uint
	USD                                          string
}

func (a *AWSElementalMediaLive) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSElementalMediaLive/current/index.json"
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
