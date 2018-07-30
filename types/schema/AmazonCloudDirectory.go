package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonCloudDirectory struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonCloudDirectory_Product
	Terms           map[string]map[string]map[string]rawAmazonCloudDirectory_Term
}

type rawAmazonCloudDirectory_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonCloudDirectory_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonCloudDirectory) UnmarshalJSON(data []byte) error {
	var p rawAmazonCloudDirectory
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonCloudDirectory_Product{}
	terms := []*AmazonCloudDirectory_Term{}

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
				pDimensions := []*AmazonCloudDirectory_Term_PriceDimensions{}
				tAttributes := []*AmazonCloudDirectory_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonCloudDirectory_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonCloudDirectory_Term{
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

type AmazonCloudDirectory struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonCloudDirectory_Product `gorm:"ForeignKey:AmazonCloudDirectoryID"`
	Terms           []*AmazonCloudDirectory_Term    `gorm:"ForeignKey:AmazonCloudDirectoryID"`
}
type AmazonCloudDirectory_Product struct {
	gorm.Model
	AmazonCloudDirectoryID uint
	Sku                    string
	ProductFamily          string
	Attributes             AmazonCloudDirectory_Product_Attributes `gorm:"ForeignKey:AmazonCloudDirectory_Product_AttributesID"`
}
type AmazonCloudDirectory_Product_Attributes struct {
	gorm.Model
	AmazonCloudDirectory_Product_AttributesID uint
	Usagetype                                 string
	Operation                                 string
	Servicename                               string
	Servicecode                               string
	Location                                  string
	LocationType                              string
	StorageClass                              string
	VolumeType                                string
}

type AmazonCloudDirectory_Term struct {
	gorm.Model
	OfferTermCode          string
	AmazonCloudDirectoryID uint
	Sku                    string
	EffectiveDate          string
	PriceDimensions        []*AmazonCloudDirectory_Term_PriceDimensions `gorm:"ForeignKey:AmazonCloudDirectory_TermID"`
	TermAttributes         []*AmazonCloudDirectory_Term_Attributes      `gorm:"ForeignKey:AmazonCloudDirectory_TermID"`
}

type AmazonCloudDirectory_Term_Attributes struct {
	gorm.Model
	AmazonCloudDirectory_TermID uint
	Key                         string
	Value                       string
}

type AmazonCloudDirectory_Term_PriceDimensions struct {
	gorm.Model
	AmazonCloudDirectory_TermID uint
	RateCode                    string
	RateType                    string
	Description                 string
	BeginRange                  string
	EndRange                    string
	Unit                        string
	PricePerUnit                *AmazonCloudDirectory_Term_PricePerUnit `gorm:"ForeignKey:AmazonCloudDirectory_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonCloudDirectory_Term_PricePerUnit struct {
	gorm.Model
	AmazonCloudDirectory_Term_PriceDimensionsID uint
	USD                                         string
}

func (a *AmazonCloudDirectory) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudDirectory/current/index.json"
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
