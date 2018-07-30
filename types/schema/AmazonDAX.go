package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonDAX struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonDAX_Product
	Terms           map[string]map[string]map[string]rawAmazonDAX_Term
}

type rawAmazonDAX_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonDAX_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonDAX) UnmarshalJSON(data []byte) error {
	var p rawAmazonDAX
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonDAX_Product{}
	terms := []*AmazonDAX_Term{}

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
				pDimensions := []*AmazonDAX_Term_PriceDimensions{}
				tAttributes := []*AmazonDAX_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonDAX_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonDAX_Term{
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

type AmazonDAX struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonDAX_Product `gorm:"ForeignKey:AmazonDAXID"`
	Terms           []*AmazonDAX_Term    `gorm:"ForeignKey:AmazonDAXID"`
}
type AmazonDAX_Product struct {
	gorm.Model
	AmazonDAXID   uint
	Sku           string
	ProductFamily string
	Attributes    AmazonDAX_Product_Attributes `gorm:"ForeignKey:AmazonDAX_Product_AttributesID"`
}
type AmazonDAX_Product_Attributes struct {
	gorm.Model
	AmazonDAX_Product_AttributesID uint
	Memory                         string
	NetworkPerformance             string
	Usagetype                      string
	Location                       string
	InstanceType                   string
	CurrentGeneration              string
	InstanceFamily                 string
	Vcpu                           string
	Operation                      string
	Servicename                    string
	Servicecode                    string
	LocationType                   string
}

type AmazonDAX_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonDAXID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonDAX_Term_PriceDimensions `gorm:"ForeignKey:AmazonDAX_TermID"`
	TermAttributes  []*AmazonDAX_Term_Attributes      `gorm:"ForeignKey:AmazonDAX_TermID"`
}

type AmazonDAX_Term_Attributes struct {
	gorm.Model
	AmazonDAX_TermID uint
	Key              string
	Value            string
}

type AmazonDAX_Term_PriceDimensions struct {
	gorm.Model
	AmazonDAX_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AmazonDAX_Term_PricePerUnit `gorm:"ForeignKey:AmazonDAX_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonDAX_Term_PricePerUnit struct {
	gorm.Model
	AmazonDAX_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AmazonDAX) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonDAX/current/index.json"
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
