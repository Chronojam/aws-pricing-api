package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonSWF struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonSWF_Product
	Terms           map[string]map[string]map[string]rawAmazonSWF_Term
}

type rawAmazonSWF_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonSWF_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonSWF) UnmarshalJSON(data []byte) error {
	var p rawAmazonSWF
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonSWF_Product{}
	terms := []*AmazonSWF_Term{}

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
				pDimensions := []*AmazonSWF_Term_PriceDimensions{}
				tAttributes := []*AmazonSWF_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonSWF_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonSWF_Term{
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

type AmazonSWF struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonSWF_Product `gorm:"ForeignKey:AmazonSWFID"`
	Terms           []*AmazonSWF_Term    `gorm:"ForeignKey:AmazonSWFID"`
}
type AmazonSWF_Product struct {
	gorm.Model
	AmazonSWFID   uint
	Sku           string
	ProductFamily string
	Attributes    AmazonSWF_Product_Attributes `gorm:"ForeignKey:AmazonSWF_Product_AttributesID"`
}
type AmazonSWF_Product_Attributes struct {
	gorm.Model
	AmazonSWF_Product_AttributesID uint
	FromLocation                   string
	FromLocationType               string
	ToLocation                     string
	ToLocationType                 string
	Usagetype                      string
	Operation                      string
	Servicecode                    string
	TransferType                   string
}

type AmazonSWF_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonSWFID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonSWF_Term_PriceDimensions `gorm:"ForeignKey:AmazonSWF_TermID"`
	TermAttributes  []*AmazonSWF_Term_Attributes      `gorm:"ForeignKey:AmazonSWF_TermID"`
}

type AmazonSWF_Term_Attributes struct {
	gorm.Model
	AmazonSWF_TermID uint
	Key              string
	Value            string
}

type AmazonSWF_Term_PriceDimensions struct {
	gorm.Model
	AmazonSWF_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AmazonSWF_Term_PricePerUnit `gorm:"ForeignKey:AmazonSWF_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonSWF_Term_PricePerUnit struct {
	gorm.Model
	AmazonSWF_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AmazonSWF) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonSWF/current/index.json"
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
