package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSElementalMediaConvert struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSElementalMediaConvert_Product
	Terms           map[string]map[string]map[string]rawAWSElementalMediaConvert_Term
}

type rawAWSElementalMediaConvert_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSElementalMediaConvert_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSElementalMediaConvert) UnmarshalJSON(data []byte) error {
	var p rawAWSElementalMediaConvert
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSElementalMediaConvert_Product{}
	terms := []*AWSElementalMediaConvert_Term{}

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
				pDimensions := []*AWSElementalMediaConvert_Term_PriceDimensions{}
				tAttributes := []*AWSElementalMediaConvert_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSElementalMediaConvert_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSElementalMediaConvert_Term{
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

type AWSElementalMediaConvert struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSElementalMediaConvert_Product `gorm:"ForeignKey:AWSElementalMediaConvertID"`
	Terms           []*AWSElementalMediaConvert_Term    `gorm:"ForeignKey:AWSElementalMediaConvertID"`
}
type AWSElementalMediaConvert_Product struct {
	gorm.Model
	AWSElementalMediaConvertID uint
	Sku                        string
	ProductFamily              string
	Attributes                 AWSElementalMediaConvert_Product_Attributes `gorm:"ForeignKey:AWSElementalMediaConvert_Product_AttributesID"`
}
type AWSElementalMediaConvert_Product_Attributes struct {
	gorm.Model
	AWSElementalMediaConvert_Product_AttributesID uint
	VideoQualitySetting                           string
	VideoResolution                               string
	Operation                                     string
	VideoCodec                                    string
	VideoFrameRate                                string
	Usagetype                                     string
	Servicename                                   string
	Tier                                          string
	TranscodingResult                             string
	Servicecode                                   string
	Location                                      string
	LocationType                                  string
}

type AWSElementalMediaConvert_Term struct {
	gorm.Model
	OfferTermCode              string
	AWSElementalMediaConvertID uint
	Sku                        string
	EffectiveDate              string
	PriceDimensions            []*AWSElementalMediaConvert_Term_PriceDimensions `gorm:"ForeignKey:AWSElementalMediaConvert_TermID"`
	TermAttributes             []*AWSElementalMediaConvert_Term_Attributes      `gorm:"ForeignKey:AWSElementalMediaConvert_TermID"`
}

type AWSElementalMediaConvert_Term_Attributes struct {
	gorm.Model
	AWSElementalMediaConvert_TermID uint
	Key                             string
	Value                           string
}

type AWSElementalMediaConvert_Term_PriceDimensions struct {
	gorm.Model
	AWSElementalMediaConvert_TermID uint
	RateCode                        string
	RateType                        string
	Description                     string
	BeginRange                      string
	EndRange                        string
	Unit                            string
	PricePerUnit                    *AWSElementalMediaConvert_Term_PricePerUnit `gorm:"ForeignKey:AWSElementalMediaConvert_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSElementalMediaConvert_Term_PricePerUnit struct {
	gorm.Model
	AWSElementalMediaConvert_Term_PriceDimensionsID uint
	USD                                             string
}

func (a *AWSElementalMediaConvert) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSElementalMediaConvert/current/index.json"
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
