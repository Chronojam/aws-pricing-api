package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSElementalMediaStore struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSElementalMediaStore_Product
	Terms           map[string]map[string]map[string]rawAWSElementalMediaStore_Term
}

type rawAWSElementalMediaStore_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSElementalMediaStore_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSElementalMediaStore) UnmarshalJSON(data []byte) error {
	var p rawAWSElementalMediaStore
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSElementalMediaStore_Product{}
	terms := []*AWSElementalMediaStore_Term{}

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
				pDimensions := []*AWSElementalMediaStore_Term_PriceDimensions{}
				tAttributes := []*AWSElementalMediaStore_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSElementalMediaStore_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSElementalMediaStore_Term{
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

type AWSElementalMediaStore struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSElementalMediaStore_Product `gorm:"ForeignKey:AWSElementalMediaStoreID"`
	Terms           []*AWSElementalMediaStore_Term    `gorm:"ForeignKey:AWSElementalMediaStoreID"`
}
type AWSElementalMediaStore_Product struct {
	gorm.Model
	AWSElementalMediaStoreID uint
	Sku                      string
	ProductFamily            string
	Attributes               AWSElementalMediaStore_Product_Attributes `gorm:"ForeignKey:AWSElementalMediaStore_Product_AttributesID"`
}
type AWSElementalMediaStore_Product_Attributes struct {
	gorm.Model
	AWSElementalMediaStore_Product_AttributesID uint
	Usagetype                                   string
	Operation                                   string
	Servicename                                 string
	Servicecode                                 string
	Description                                 string
	Location                                    string
	LocationType                                string
	Availability                                string
}

type AWSElementalMediaStore_Term struct {
	gorm.Model
	OfferTermCode            string
	AWSElementalMediaStoreID uint
	Sku                      string
	EffectiveDate            string
	PriceDimensions          []*AWSElementalMediaStore_Term_PriceDimensions `gorm:"ForeignKey:AWSElementalMediaStore_TermID"`
	TermAttributes           []*AWSElementalMediaStore_Term_Attributes      `gorm:"ForeignKey:AWSElementalMediaStore_TermID"`
}

type AWSElementalMediaStore_Term_Attributes struct {
	gorm.Model
	AWSElementalMediaStore_TermID uint
	Key                           string
	Value                         string
}

type AWSElementalMediaStore_Term_PriceDimensions struct {
	gorm.Model
	AWSElementalMediaStore_TermID uint
	RateCode                      string
	RateType                      string
	Description                   string
	BeginRange                    string
	EndRange                      string
	Unit                          string
	PricePerUnit                  *AWSElementalMediaStore_Term_PricePerUnit `gorm:"ForeignKey:AWSElementalMediaStore_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSElementalMediaStore_Term_PricePerUnit struct {
	gorm.Model
	AWSElementalMediaStore_Term_PriceDimensionsID uint
	USD                                           string
}

func (a *AWSElementalMediaStore) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSElementalMediaStore/current/index.json"
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
