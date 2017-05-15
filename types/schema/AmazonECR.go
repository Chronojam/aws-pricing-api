package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonECR struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonECR_Product
	Terms           map[string]map[string]map[string]rawAmazonECR_Term
}

type rawAmazonECR_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonECR_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonECR) UnmarshalJSON(data []byte) error {
	var p rawAmazonECR
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonECR_Product{}
	terms := []*AmazonECR_Term{}

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
				pDimensions := []*AmazonECR_Term_PriceDimensions{}
				tAttributes := []*AmazonECR_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonECR_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonECR_Term{
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

type AmazonECR struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonECR_Product `gorm:"ForeignKey:AmazonECRID"`
	Terms           []*AmazonECR_Term    `gorm:"ForeignKey:AmazonECRID"`
}
type AmazonECR_Product struct {
	gorm.Model
	AmazonECRID   uint
	Sku           string
	ProductFamily string
	Attributes    AmazonECR_Product_Attributes `gorm:"ForeignKey:AmazonECR_Product_AttributesID"`
}
type AmazonECR_Product_Attributes struct {
	gorm.Model
	AmazonECR_Product_AttributesID uint
	ToLocation                     string
	ToLocationType                 string
	Usagetype                      string
	Operation                      string
	Servicecode                    string
	TransferType                   string
	FromLocation                   string
	FromLocationType               string
}

type AmazonECR_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonECRID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonECR_Term_PriceDimensions `gorm:"ForeignKey:AmazonECR_TermID"`
	TermAttributes  []*AmazonECR_Term_Attributes      `gorm:"ForeignKey:AmazonECR_TermID"`
}

type AmazonECR_Term_Attributes struct {
	gorm.Model
	AmazonECR_TermID uint
	Key              string
	Value            string
}

type AmazonECR_Term_PriceDimensions struct {
	gorm.Model
	AmazonECR_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AmazonECR_Term_PricePerUnit `gorm:"ForeignKey:AmazonECR_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonECR_Term_PricePerUnit struct {
	gorm.Model
	AmazonECR_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AmazonECR) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonECR/current/index.json"
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
