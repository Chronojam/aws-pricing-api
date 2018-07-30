package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSAppSync struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSAppSync_Product
	Terms           map[string]map[string]map[string]rawAWSAppSync_Term
}

type rawAWSAppSync_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSAppSync_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSAppSync) UnmarshalJSON(data []byte) error {
	var p rawAWSAppSync
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSAppSync_Product{}
	terms := []*AWSAppSync_Term{}

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
				pDimensions := []*AWSAppSync_Term_PriceDimensions{}
				tAttributes := []*AWSAppSync_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSAppSync_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSAppSync_Term{
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

type AWSAppSync struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSAppSync_Product `gorm:"ForeignKey:AWSAppSyncID"`
	Terms           []*AWSAppSync_Term    `gorm:"ForeignKey:AWSAppSyncID"`
}
type AWSAppSync_Product struct {
	gorm.Model
	AWSAppSyncID  uint
	Sku           string
	ProductFamily string
	Attributes    AWSAppSync_Product_Attributes `gorm:"ForeignKey:AWSAppSync_Product_AttributesID"`
}
type AWSAppSync_Product_Attributes struct {
	gorm.Model
	AWSAppSync_Product_AttributesID uint
	Servicecode                     string
	ToLocationType                  string
	Usagetype                       string
	Operation                       string
	Servicename                     string
	TransferType                    string
	FromLocation                    string
	FromLocationType                string
	ToLocation                      string
}

type AWSAppSync_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSAppSyncID    uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSAppSync_Term_PriceDimensions `gorm:"ForeignKey:AWSAppSync_TermID"`
	TermAttributes  []*AWSAppSync_Term_Attributes      `gorm:"ForeignKey:AWSAppSync_TermID"`
}

type AWSAppSync_Term_Attributes struct {
	gorm.Model
	AWSAppSync_TermID uint
	Key               string
	Value             string
}

type AWSAppSync_Term_PriceDimensions struct {
	gorm.Model
	AWSAppSync_TermID uint
	RateCode          string
	RateType          string
	Description       string
	BeginRange        string
	EndRange          string
	Unit              string
	PricePerUnit      *AWSAppSync_Term_PricePerUnit `gorm:"ForeignKey:AWSAppSync_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSAppSync_Term_PricePerUnit struct {
	gorm.Model
	AWSAppSync_Term_PriceDimensionsID uint
	USD                               string
}

func (a *AWSAppSync) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSAppSync/current/index.json"
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
