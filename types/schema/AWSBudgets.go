package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSBudgets struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSBudgets_Product
	Terms           map[string]map[string]map[string]rawAWSBudgets_Term
}

type rawAWSBudgets_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSBudgets_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSBudgets) UnmarshalJSON(data []byte) error {
	var p rawAWSBudgets
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSBudgets_Product{}
	terms := []*AWSBudgets_Term{}

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
				pDimensions := []*AWSBudgets_Term_PriceDimensions{}
				tAttributes := []*AWSBudgets_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSBudgets_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSBudgets_Term{
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

type AWSBudgets struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSBudgets_Product `gorm:"ForeignKey:AWSBudgetsID"`
	Terms           []*AWSBudgets_Term    `gorm:"ForeignKey:AWSBudgetsID"`
}
type AWSBudgets_Product struct {
	gorm.Model
	AWSBudgetsID  uint
	Sku           string
	ProductFamily string
	Attributes    AWSBudgets_Product_Attributes `gorm:"ForeignKey:AWSBudgets_Product_AttributesID"`
}
type AWSBudgets_Product_Attributes struct {
	gorm.Model
	AWSBudgets_Product_AttributesID uint
	Servicecode                     string
	Location                        string
	LocationType                    string
	GroupDescription                string
	Usagetype                       string
	Operation                       string
}

type AWSBudgets_Term struct {
	gorm.Model
	OfferTermCode   string
	AWSBudgetsID    uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AWSBudgets_Term_PriceDimensions `gorm:"ForeignKey:AWSBudgets_TermID"`
	TermAttributes  []*AWSBudgets_Term_Attributes      `gorm:"ForeignKey:AWSBudgets_TermID"`
}

type AWSBudgets_Term_Attributes struct {
	gorm.Model
	AWSBudgets_TermID uint
	Key               string
	Value             string
}

type AWSBudgets_Term_PriceDimensions struct {
	gorm.Model
	AWSBudgets_TermID uint
	RateCode          string
	RateType          string
	Description       string
	BeginRange        string
	EndRange          string
	Unit              string
	PricePerUnit      *AWSBudgets_Term_PricePerUnit `gorm:"ForeignKey:AWSBudgets_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSBudgets_Term_PricePerUnit struct {
	gorm.Model
	AWSBudgets_Term_PriceDimensionsID uint
	USD                               string
}

func (a *AWSBudgets) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSBudgets/current/index.json"
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
