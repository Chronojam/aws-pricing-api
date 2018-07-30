package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonGuardDuty struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonGuardDuty_Product
	Terms           map[string]map[string]map[string]rawAmazonGuardDuty_Term
}

type rawAmazonGuardDuty_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonGuardDuty_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonGuardDuty) UnmarshalJSON(data []byte) error {
	var p rawAmazonGuardDuty
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonGuardDuty_Product{}
	terms := []*AmazonGuardDuty_Term{}

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
				pDimensions := []*AmazonGuardDuty_Term_PriceDimensions{}
				tAttributes := []*AmazonGuardDuty_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonGuardDuty_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonGuardDuty_Term{
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

type AmazonGuardDuty struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonGuardDuty_Product `gorm:"ForeignKey:AmazonGuardDutyID"`
	Terms           []*AmazonGuardDuty_Term    `gorm:"ForeignKey:AmazonGuardDutyID"`
}
type AmazonGuardDuty_Product struct {
	gorm.Model
	AmazonGuardDutyID uint
	Sku               string
	ProductFamily     string
	Attributes        AmazonGuardDuty_Product_Attributes `gorm:"ForeignKey:AmazonGuardDuty_Product_AttributesID"`
}
type AmazonGuardDuty_Product_Attributes struct {
	gorm.Model
	AmazonGuardDuty_Product_AttributesID uint
	LocationType                         string
	Group                                string
	Usagetype                            string
	Operation                            string
	Servicename                          string
	Servicecode                          string
	Location                             string
}

type AmazonGuardDuty_Term struct {
	gorm.Model
	OfferTermCode     string
	AmazonGuardDutyID uint
	Sku               string
	EffectiveDate     string
	PriceDimensions   []*AmazonGuardDuty_Term_PriceDimensions `gorm:"ForeignKey:AmazonGuardDuty_TermID"`
	TermAttributes    []*AmazonGuardDuty_Term_Attributes      `gorm:"ForeignKey:AmazonGuardDuty_TermID"`
}

type AmazonGuardDuty_Term_Attributes struct {
	gorm.Model
	AmazonGuardDuty_TermID uint
	Key                    string
	Value                  string
}

type AmazonGuardDuty_Term_PriceDimensions struct {
	gorm.Model
	AmazonGuardDuty_TermID uint
	RateCode               string
	RateType               string
	Description            string
	BeginRange             string
	EndRange               string
	Unit                   string
	PricePerUnit           *AmazonGuardDuty_Term_PricePerUnit `gorm:"ForeignKey:AmazonGuardDuty_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonGuardDuty_Term_PricePerUnit struct {
	gorm.Model
	AmazonGuardDuty_Term_PriceDimensionsID uint
	USD                                    string
}

func (a *AmazonGuardDuty) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonGuardDuty/current/index.json"
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
