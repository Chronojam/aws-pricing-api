package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonGameLift struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonGameLift_Product
	Terms           map[string]map[string]map[string]rawAmazonGameLift_Term
}

type rawAmazonGameLift_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonGameLift_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonGameLift) UnmarshalJSON(data []byte) error {
	var p rawAmazonGameLift
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonGameLift_Product{}
	terms := []*AmazonGameLift_Term{}

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
				pDimensions := []*AmazonGameLift_Term_PriceDimensions{}
				tAttributes := []*AmazonGameLift_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonGameLift_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonGameLift_Term{
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

type AmazonGameLift struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonGameLift_Product `gorm:"ForeignKey:AmazonGameLiftID"`
	Terms           []*AmazonGameLift_Term    `gorm:"ForeignKey:AmazonGameLiftID"`
}
type AmazonGameLift_Product struct {
	gorm.Model
	AmazonGameLiftID uint
	Sku              string
	ProductFamily    string
	Attributes       AmazonGameLift_Product_Attributes `gorm:"ForeignKey:AmazonGameLift_Product_AttributesID"`
}
type AmazonGameLift_Product_Attributes struct {
	gorm.Model
	AmazonGameLift_Product_AttributesID uint
	ToLocation                          string
	ToLocationType                      string
	Usagetype                           string
	Operation                           string
	Servicecode                         string
	TransferType                        string
	FromLocation                        string
	FromLocationType                    string
}

type AmazonGameLift_Term struct {
	gorm.Model
	OfferTermCode    string
	AmazonGameLiftID uint
	Sku              string
	EffectiveDate    string
	PriceDimensions  []*AmazonGameLift_Term_PriceDimensions `gorm:"ForeignKey:AmazonGameLift_TermID"`
	TermAttributes   []*AmazonGameLift_Term_Attributes      `gorm:"ForeignKey:AmazonGameLift_TermID"`
}

type AmazonGameLift_Term_Attributes struct {
	gorm.Model
	AmazonGameLift_TermID uint
	Key                   string
	Value                 string
}

type AmazonGameLift_Term_PriceDimensions struct {
	gorm.Model
	AmazonGameLift_TermID uint
	RateCode              string
	RateType              string
	Description           string
	BeginRange            string
	EndRange              string
	Unit                  string
	PricePerUnit          *AmazonGameLift_Term_PricePerUnit `gorm:"ForeignKey:AmazonGameLift_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonGameLift_Term_PricePerUnit struct {
	gorm.Model
	AmazonGameLift_Term_PriceDimensionsID uint
	USD                                   string
}

func (a *AmazonGameLift) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonGameLift/current/index.json"
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
