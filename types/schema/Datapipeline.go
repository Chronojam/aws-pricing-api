package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawDatapipeline struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]Datapipeline_Product
	Terms           map[string]map[string]map[string]rawDatapipeline_Term
}

type rawDatapipeline_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]Datapipeline_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *Datapipeline) UnmarshalJSON(data []byte) error {
	var p rawDatapipeline
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*Datapipeline_Product{}
	terms := []*Datapipeline_Term{}

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
				pDimensions := []*Datapipeline_Term_PriceDimensions{}
				tAttributes := []*Datapipeline_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := Datapipeline_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := Datapipeline_Term{
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

type Datapipeline struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*Datapipeline_Product `gorm:"ForeignKey:DatapipelineID"`
	Terms           []*Datapipeline_Term    `gorm:"ForeignKey:DatapipelineID"`
}
type Datapipeline_Product struct {
	gorm.Model
	DatapipelineID uint
	Sku            string
	ProductFamily  string
	Attributes     Datapipeline_Product_Attributes `gorm:"ForeignKey:Datapipeline_Product_AttributesID"`
}
type Datapipeline_Product_Attributes struct {
	gorm.Model
	Datapipeline_Product_AttributesID uint
	Servicecode                       string
	LocationType                      string
	Operation                         string
	ExecutionLocation                 string
	Location                          string
	Group                             string
	Usagetype                         string
	ExecutionFrequency                string
	FrequencyMode                     string
}

type Datapipeline_Term struct {
	gorm.Model
	OfferTermCode   string
	DatapipelineID  uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*Datapipeline_Term_PriceDimensions `gorm:"ForeignKey:Datapipeline_TermID"`
	TermAttributes  []*Datapipeline_Term_Attributes      `gorm:"ForeignKey:Datapipeline_TermID"`
}

type Datapipeline_Term_Attributes struct {
	gorm.Model
	Datapipeline_TermID uint
	Key                 string
	Value               string
}

type Datapipeline_Term_PriceDimensions struct {
	gorm.Model
	Datapipeline_TermID uint
	RateCode            string
	RateType            string
	Description         string
	BeginRange          string
	EndRange            string
	Unit                string
	PricePerUnit        *Datapipeline_Term_PricePerUnit `gorm:"ForeignKey:Datapipeline_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type Datapipeline_Term_PricePerUnit struct {
	gorm.Model
	Datapipeline_Term_PriceDimensionsID uint
	USD                                 string
}

func (a *Datapipeline) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/datapipeline/current/index.json"
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
