package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSCostExplorer struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSCostExplorer_Product
	Terms           map[string]map[string]map[string]rawAWSCostExplorer_Term
}

type rawAWSCostExplorer_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSCostExplorer_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSCostExplorer) UnmarshalJSON(data []byte) error {
	var p rawAWSCostExplorer
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSCostExplorer_Product{}
	terms := []*AWSCostExplorer_Term{}

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
				pDimensions := []*AWSCostExplorer_Term_PriceDimensions{}
				tAttributes := []*AWSCostExplorer_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSCostExplorer_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSCostExplorer_Term{
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

type AWSCostExplorer struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSCostExplorer_Product `gorm:"ForeignKey:AWSCostExplorerID"`
	Terms           []*AWSCostExplorer_Term    `gorm:"ForeignKey:AWSCostExplorerID"`
}
type AWSCostExplorer_Product struct {
	gorm.Model
	AWSCostExplorerID uint
	Attributes        AWSCostExplorer_Product_Attributes `gorm:"ForeignKey:AWSCostExplorer_Product_AttributesID"`
	Sku               string
	ProductFamily     string
}
type AWSCostExplorer_Product_Attributes struct {
	gorm.Model
	AWSCostExplorer_Product_AttributesID uint
	Servicecode                          string
	Location                             string
	LocationType                         string
	Usagetype                            string
	Operation                            string
	RequestType                          string
	Servicename                          string
}

type AWSCostExplorer_Term struct {
	gorm.Model
	OfferTermCode     string
	AWSCostExplorerID uint
	Sku               string
	EffectiveDate     string
	PriceDimensions   []*AWSCostExplorer_Term_PriceDimensions `gorm:"ForeignKey:AWSCostExplorer_TermID"`
	TermAttributes    []*AWSCostExplorer_Term_Attributes      `gorm:"ForeignKey:AWSCostExplorer_TermID"`
}

type AWSCostExplorer_Term_Attributes struct {
	gorm.Model
	AWSCostExplorer_TermID uint
	Key                    string
	Value                  string
}

type AWSCostExplorer_Term_PriceDimensions struct {
	gorm.Model
	AWSCostExplorer_TermID uint
	RateCode               string
	RateType               string
	Description            string
	BeginRange             string
	EndRange               string
	Unit                   string
	PricePerUnit           *AWSCostExplorer_Term_PricePerUnit `gorm:"ForeignKey:AWSCostExplorer_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSCostExplorer_Term_PricePerUnit struct {
	gorm.Model
	AWSCostExplorer_Term_PriceDimensionsID uint
	USD                                    string
}

func (a *AWSCostExplorer) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCostExplorer/current/index.json"
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
