package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonWAM struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonWAM_Product
	Terms           map[string]map[string]map[string]rawAmazonWAM_Term
}

type rawAmazonWAM_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonWAM_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonWAM) UnmarshalJSON(data []byte) error {
	var p rawAmazonWAM
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonWAM_Product{}
	terms := []*AmazonWAM_Term{}

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
				pDimensions := []*AmazonWAM_Term_PriceDimensions{}
				tAttributes := []*AmazonWAM_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonWAM_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonWAM_Term{
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

type AmazonWAM struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonWAM_Product `gorm:"ForeignKey:AmazonWAMID"`
	Terms           []*AmazonWAM_Term    `gorm:"ForeignKey:AmazonWAMID"`
}
type AmazonWAM_Product struct {
	gorm.Model
	AmazonWAMID   uint
	ProductFamily string
	Attributes    AmazonWAM_Product_Attributes `gorm:"ForeignKey:AmazonWAM_Product_AttributesID"`
	Sku           string
}
type AmazonWAM_Product_Attributes struct {
	gorm.Model
	AmazonWAM_Product_AttributesID uint
	Operation                      string
	PlanType                       string
	Servicecode                    string
	Location                       string
	LocationType                   string
	Group                          string
	Usagetype                      string
}

type AmazonWAM_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonWAMID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonWAM_Term_PriceDimensions `gorm:"ForeignKey:AmazonWAM_TermID"`
	TermAttributes  []*AmazonWAM_Term_Attributes      `gorm:"ForeignKey:AmazonWAM_TermID"`
}

type AmazonWAM_Term_Attributes struct {
	gorm.Model
	AmazonWAM_TermID uint
	Key              string
	Value            string
}

type AmazonWAM_Term_PriceDimensions struct {
	gorm.Model
	AmazonWAM_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AmazonWAM_Term_PricePerUnit `gorm:"ForeignKey:AmazonWAM_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonWAM_Term_PricePerUnit struct {
	gorm.Model
	AmazonWAM_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AmazonWAM) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWAM/current/index.json"
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
