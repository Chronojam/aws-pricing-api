package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawElasticMapReduce struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]ElasticMapReduce_Product
	Terms           map[string]map[string]map[string]rawElasticMapReduce_Term
}

type rawElasticMapReduce_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]ElasticMapReduce_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *ElasticMapReduce) UnmarshalJSON(data []byte) error {
	var p rawElasticMapReduce
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*ElasticMapReduce_Product{}
	terms := []*ElasticMapReduce_Term{}

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
				pDimensions := []*ElasticMapReduce_Term_PriceDimensions{}
				tAttributes := []*ElasticMapReduce_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := ElasticMapReduce_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := ElasticMapReduce_Term{
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

type ElasticMapReduce struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*ElasticMapReduce_Product `gorm:"ForeignKey:ElasticMapReduceID"`
	Terms           []*ElasticMapReduce_Term    `gorm:"ForeignKey:ElasticMapReduceID"`
}
type ElasticMapReduce_Product struct {
	gorm.Model
	ElasticMapReduceID uint
	ProductFamily      string
	Attributes         ElasticMapReduce_Product_Attributes `gorm:"ForeignKey:ElasticMapReduce_Product_AttributesID"`
	Sku                string
}
type ElasticMapReduce_Product_Attributes struct {
	gorm.Model
	ElasticMapReduce_Product_AttributesID uint
	Location                              string
	LocationType                          string
	Usagetype                             string
	Operation                             string
	Servicecode                           string
	InstanceFamily                        string
	Servicename                           string
	SoftwareType                          string
	InstanceType                          string
}

type ElasticMapReduce_Term struct {
	gorm.Model
	OfferTermCode      string
	ElasticMapReduceID uint
	Sku                string
	EffectiveDate      string
	PriceDimensions    []*ElasticMapReduce_Term_PriceDimensions `gorm:"ForeignKey:ElasticMapReduce_TermID"`
	TermAttributes     []*ElasticMapReduce_Term_Attributes      `gorm:"ForeignKey:ElasticMapReduce_TermID"`
}

type ElasticMapReduce_Term_Attributes struct {
	gorm.Model
	ElasticMapReduce_TermID uint
	Key                     string
	Value                   string
}

type ElasticMapReduce_Term_PriceDimensions struct {
	gorm.Model
	ElasticMapReduce_TermID uint
	RateCode                string
	RateType                string
	Description             string
	BeginRange              string
	EndRange                string
	Unit                    string
	PricePerUnit            *ElasticMapReduce_Term_PricePerUnit `gorm:"ForeignKey:ElasticMapReduce_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type ElasticMapReduce_Term_PricePerUnit struct {
	gorm.Model
	ElasticMapReduce_Term_PriceDimensionsID uint
	USD                                     string
}

func (a *ElasticMapReduce) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/ElasticMapReduce/current/index.json"
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
