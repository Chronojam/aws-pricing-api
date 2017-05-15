package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawCodeBuild struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]CodeBuild_Product
	Terms           map[string]map[string]map[string]rawCodeBuild_Term
}

type rawCodeBuild_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]CodeBuild_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *CodeBuild) UnmarshalJSON(data []byte) error {
	var p rawCodeBuild
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*CodeBuild_Product{}
	terms := []*CodeBuild_Term{}

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
				pDimensions := []*CodeBuild_Term_PriceDimensions{}
				tAttributes := []*CodeBuild_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := CodeBuild_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := CodeBuild_Term{
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

type CodeBuild struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*CodeBuild_Product `gorm:"ForeignKey:CodeBuildID"`
	Terms           []*CodeBuild_Term    `gorm:"ForeignKey:CodeBuildID"`
}
type CodeBuild_Product struct {
	gorm.Model
	CodeBuildID   uint
	Attributes    CodeBuild_Product_Attributes `gorm:"ForeignKey:CodeBuild_Product_AttributesID"`
	Sku           string
	ProductFamily string
}
type CodeBuild_Product_Attributes struct {
	gorm.Model
	CodeBuild_Product_AttributesID uint
	Location                       string
	OperatingSystem                string
	Operation                      string
	ComputeFamily                  string
	ComputeType                    string
	Servicecode                    string
	LocationType                   string
	Vcpu                           string
	Memory                         string
	Usagetype                      string
}

type CodeBuild_Term struct {
	gorm.Model
	OfferTermCode   string
	CodeBuildID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*CodeBuild_Term_PriceDimensions `gorm:"ForeignKey:CodeBuild_TermID"`
	TermAttributes  []*CodeBuild_Term_Attributes      `gorm:"ForeignKey:CodeBuild_TermID"`
}

type CodeBuild_Term_Attributes struct {
	gorm.Model
	CodeBuild_TermID uint
	Key              string
	Value            string
}

type CodeBuild_Term_PriceDimensions struct {
	gorm.Model
	CodeBuild_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *CodeBuild_Term_PricePerUnit `gorm:"ForeignKey:CodeBuild_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type CodeBuild_Term_PricePerUnit struct {
	gorm.Model
	CodeBuild_Term_PriceDimensionsID uint
	USD                              string
}

func (a *CodeBuild) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/CodeBuild/current/index.json"
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
