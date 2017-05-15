package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonCognitoSync struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonCognitoSync_Product
	Terms           map[string]map[string]map[string]rawAmazonCognitoSync_Term
}

type rawAmazonCognitoSync_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonCognitoSync_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonCognitoSync) UnmarshalJSON(data []byte) error {
	var p rawAmazonCognitoSync
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonCognitoSync_Product{}
	terms := []*AmazonCognitoSync_Term{}

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
				pDimensions := []*AmazonCognitoSync_Term_PriceDimensions{}
				tAttributes := []*AmazonCognitoSync_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonCognitoSync_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonCognitoSync_Term{
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

type AmazonCognitoSync struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonCognitoSync_Product `gorm:"ForeignKey:AmazonCognitoSyncID"`
	Terms           []*AmazonCognitoSync_Term    `gorm:"ForeignKey:AmazonCognitoSyncID"`
}
type AmazonCognitoSync_Product struct {
	gorm.Model
	AmazonCognitoSyncID uint
	Sku                 string
	ProductFamily       string
	Attributes          AmazonCognitoSync_Product_Attributes `gorm:"ForeignKey:AmazonCognitoSync_Product_AttributesID"`
}
type AmazonCognitoSync_Product_Attributes struct {
	gorm.Model
	AmazonCognitoSync_Product_AttributesID uint
	Servicecode                            string
	Location                               string
	LocationType                           string
	Usagetype                              string
	Operation                              string
}

type AmazonCognitoSync_Term struct {
	gorm.Model
	OfferTermCode       string
	AmazonCognitoSyncID uint
	Sku                 string
	EffectiveDate       string
	PriceDimensions     []*AmazonCognitoSync_Term_PriceDimensions `gorm:"ForeignKey:AmazonCognitoSync_TermID"`
	TermAttributes      []*AmazonCognitoSync_Term_Attributes      `gorm:"ForeignKey:AmazonCognitoSync_TermID"`
}

type AmazonCognitoSync_Term_Attributes struct {
	gorm.Model
	AmazonCognitoSync_TermID uint
	Key                      string
	Value                    string
}

type AmazonCognitoSync_Term_PriceDimensions struct {
	gorm.Model
	AmazonCognitoSync_TermID uint
	RateCode                 string
	RateType                 string
	Description              string
	BeginRange               string
	EndRange                 string
	Unit                     string
	PricePerUnit             *AmazonCognitoSync_Term_PricePerUnit `gorm:"ForeignKey:AmazonCognitoSync_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonCognitoSync_Term_PricePerUnit struct {
	gorm.Model
	AmazonCognitoSync_Term_PriceDimensionsID uint
	USD                                      string
}

func (a *AmazonCognitoSync) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCognitoSync/current/index.json"
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
