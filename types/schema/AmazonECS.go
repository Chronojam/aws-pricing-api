package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonECS struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonECS_Product
	Terms           map[string]map[string]map[string]rawAmazonECS_Term
}

type rawAmazonECS_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonECS_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonECS) UnmarshalJSON(data []byte) error {
	var p rawAmazonECS
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonECS_Product{}
	terms := []*AmazonECS_Term{}

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
				pDimensions := []*AmazonECS_Term_PriceDimensions{}
				tAttributes := []*AmazonECS_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonECS_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonECS_Term{
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

type AmazonECS struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonECS_Product `gorm:"ForeignKey:AmazonECSID"`
	Terms           []*AmazonECS_Term    `gorm:"ForeignKey:AmazonECSID"`
}
type AmazonECS_Product struct {
	gorm.Model
	AmazonECSID   uint
	Sku           string
	ProductFamily string
	Attributes    AmazonECS_Product_Attributes `gorm:"ForeignKey:AmazonECS_Product_AttributesID"`
}
type AmazonECS_Product_Attributes struct {
	gorm.Model
	AmazonECS_Product_AttributesID uint
	Servicecode                    string
	ToLocation                     string
	ToLocationType                 string
	Usagetype                      string
	TransferType                   string
	FromLocation                   string
	FromLocationType               string
	Operation                      string
	Servicename                    string
}

type AmazonECS_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonECSID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonECS_Term_PriceDimensions `gorm:"ForeignKey:AmazonECS_TermID"`
	TermAttributes  []*AmazonECS_Term_Attributes      `gorm:"ForeignKey:AmazonECS_TermID"`
}

type AmazonECS_Term_Attributes struct {
	gorm.Model
	AmazonECS_TermID uint
	Key              string
	Value            string
}

type AmazonECS_Term_PriceDimensions struct {
	gorm.Model
	AmazonECS_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AmazonECS_Term_PricePerUnit `gorm:"ForeignKey:AmazonECS_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonECS_Term_PricePerUnit struct {
	gorm.Model
	AmazonECS_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AmazonECS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonECS/current/index.json"
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
