package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonLightsail struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonLightsail_Product
	Terms           map[string]map[string]map[string]rawAmazonLightsail_Term
}

type rawAmazonLightsail_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonLightsail_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonLightsail) UnmarshalJSON(data []byte) error {
	var p rawAmazonLightsail
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonLightsail_Product{}
	terms := []*AmazonLightsail_Term{}

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
				pDimensions := []*AmazonLightsail_Term_PriceDimensions{}
				tAttributes := []*AmazonLightsail_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonLightsail_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonLightsail_Term{
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

type AmazonLightsail struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonLightsail_Product `gorm:"ForeignKey:AmazonLightsailID"`
	Terms           []*AmazonLightsail_Term    `gorm:"ForeignKey:AmazonLightsailID"`
}
type AmazonLightsail_Product struct {
	gorm.Model
	AmazonLightsailID uint
	ProductFamily     string
	Attributes        AmazonLightsail_Product_Attributes `gorm:"ForeignKey:AmazonLightsail_Product_AttributesID"`
	Sku               string
}
type AmazonLightsail_Product_Attributes struct {
	gorm.Model
	AmazonLightsail_Product_AttributesID uint
	GroupDescription                     string
	FromLocationType                     string
	ToLocationType                       string
	Usagetype                            string
	CountsAgainstQuota                   string
	Servicename                          string
	Operation                            string
	FreeOverage                          string
	TransferType                         string
	FromLocation                         string
	ToLocation                           string
	Servicecode                          string
	Group                                string
	OverageType                          string
}

type AmazonLightsail_Term struct {
	gorm.Model
	OfferTermCode     string
	AmazonLightsailID uint
	Sku               string
	EffectiveDate     string
	PriceDimensions   []*AmazonLightsail_Term_PriceDimensions `gorm:"ForeignKey:AmazonLightsail_TermID"`
	TermAttributes    []*AmazonLightsail_Term_Attributes      `gorm:"ForeignKey:AmazonLightsail_TermID"`
}

type AmazonLightsail_Term_Attributes struct {
	gorm.Model
	AmazonLightsail_TermID uint
	Key                    string
	Value                  string
}

type AmazonLightsail_Term_PriceDimensions struct {
	gorm.Model
	AmazonLightsail_TermID uint
	RateCode               string
	RateType               string
	Description            string
	BeginRange             string
	EndRange               string
	Unit                   string
	PricePerUnit           *AmazonLightsail_Term_PricePerUnit `gorm:"ForeignKey:AmazonLightsail_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonLightsail_Term_PricePerUnit struct {
	gorm.Model
	AmazonLightsail_Term_PriceDimensionsID uint
	USD                                    string
}

func (a *AmazonLightsail) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonLightsail/current/index.json"
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
