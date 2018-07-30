package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonKinesisVideo struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonKinesisVideo_Product
	Terms           map[string]map[string]map[string]rawAmazonKinesisVideo_Term
}

type rawAmazonKinesisVideo_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonKinesisVideo_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonKinesisVideo) UnmarshalJSON(data []byte) error {
	var p rawAmazonKinesisVideo
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonKinesisVideo_Product{}
	terms := []*AmazonKinesisVideo_Term{}

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
				pDimensions := []*AmazonKinesisVideo_Term_PriceDimensions{}
				tAttributes := []*AmazonKinesisVideo_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonKinesisVideo_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonKinesisVideo_Term{
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

type AmazonKinesisVideo struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonKinesisVideo_Product `gorm:"ForeignKey:AmazonKinesisVideoID"`
	Terms           []*AmazonKinesisVideo_Term    `gorm:"ForeignKey:AmazonKinesisVideoID"`
}
type AmazonKinesisVideo_Product struct {
	gorm.Model
	AmazonKinesisVideoID uint
	ProductFamily        string
	Attributes           AmazonKinesisVideo_Product_Attributes `gorm:"ForeignKey:AmazonKinesisVideo_Product_AttributesID"`
	Sku                  string
}
type AmazonKinesisVideo_Product_Attributes struct {
	gorm.Model
	AmazonKinesisVideo_Product_AttributesID uint
	Servicecode                             string
	FromLocation                            string
	ToLocation                              string
	ToLocationType                          string
	Operation                               string
	TransferType                            string
	FromLocationType                        string
	Usagetype                               string
	Servicename                             string
}

type AmazonKinesisVideo_Term struct {
	gorm.Model
	OfferTermCode        string
	AmazonKinesisVideoID uint
	Sku                  string
	EffectiveDate        string
	PriceDimensions      []*AmazonKinesisVideo_Term_PriceDimensions `gorm:"ForeignKey:AmazonKinesisVideo_TermID"`
	TermAttributes       []*AmazonKinesisVideo_Term_Attributes      `gorm:"ForeignKey:AmazonKinesisVideo_TermID"`
}

type AmazonKinesisVideo_Term_Attributes struct {
	gorm.Model
	AmazonKinesisVideo_TermID uint
	Key                       string
	Value                     string
}

type AmazonKinesisVideo_Term_PriceDimensions struct {
	gorm.Model
	AmazonKinesisVideo_TermID uint
	RateCode                  string
	RateType                  string
	Description               string
	BeginRange                string
	EndRange                  string
	Unit                      string
	PricePerUnit              *AmazonKinesisVideo_Term_PricePerUnit `gorm:"ForeignKey:AmazonKinesisVideo_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonKinesisVideo_Term_PricePerUnit struct {
	gorm.Model
	AmazonKinesisVideo_Term_PriceDimensionsID uint
	USD                                       string
}

func (a *AmazonKinesisVideo) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesisVideo/current/index.json"
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
