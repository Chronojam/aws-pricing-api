package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonSNS struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonSNS_Product
	Terms           map[string]map[string]map[string]rawAmazonSNS_Term
}

type rawAmazonSNS_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonSNS_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonSNS) UnmarshalJSON(data []byte) error {
	var p rawAmazonSNS
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonSNS_Product{}
	terms := []*AmazonSNS_Term{}

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
				pDimensions := []*AmazonSNS_Term_PriceDimensions{}
				tAttributes := []*AmazonSNS_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonSNS_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonSNS_Term{
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

type AmazonSNS struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonSNS_Product `gorm:"ForeignKey:AmazonSNSID"`
	Terms           []*AmazonSNS_Term    `gorm:"ForeignKey:AmazonSNSID"`
}
type AmazonSNS_Product struct {
	gorm.Model
	AmazonSNSID   uint
	ProductFamily string
	Attributes    AmazonSNS_Product_Attributes `gorm:"ForeignKey:AmazonSNS_Product_AttributesID"`
	Sku           string
}
type AmazonSNS_Product_Attributes struct {
	gorm.Model
	AmazonSNS_Product_AttributesID uint
	FromLocationType               string
	ToLocationType                 string
	Operation                      string
	Servicecode                    string
	FromLocation                   string
	ToLocation                     string
	Usagetype                      string
	Servicename                    string
	TransferType                   string
}

type AmazonSNS_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonSNSID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonSNS_Term_PriceDimensions `gorm:"ForeignKey:AmazonSNS_TermID"`
	TermAttributes  []*AmazonSNS_Term_Attributes      `gorm:"ForeignKey:AmazonSNS_TermID"`
}

type AmazonSNS_Term_Attributes struct {
	gorm.Model
	AmazonSNS_TermID uint
	Key              string
	Value            string
}

type AmazonSNS_Term_PriceDimensions struct {
	gorm.Model
	AmazonSNS_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AmazonSNS_Term_PricePerUnit `gorm:"ForeignKey:AmazonSNS_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonSNS_Term_PricePerUnit struct {
	gorm.Model
	AmazonSNS_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AmazonSNS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonSNS/current/index.json"
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
