package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonEKS struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonEKS_Product
	Terms           map[string]map[string]map[string]rawAmazonEKS_Term
}

type rawAmazonEKS_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonEKS_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonEKS) UnmarshalJSON(data []byte) error {
	var p rawAmazonEKS
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonEKS_Product{}
	terms := []*AmazonEKS_Term{}

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
				pDimensions := []*AmazonEKS_Term_PriceDimensions{}
				tAttributes := []*AmazonEKS_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonEKS_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonEKS_Term{
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

type AmazonEKS struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonEKS_Product `gorm:"ForeignKey:AmazonEKSID"`
	Terms           []*AmazonEKS_Term    `gorm:"ForeignKey:AmazonEKSID"`
}
type AmazonEKS_Product struct {
	gorm.Model
	AmazonEKSID   uint
	Sku           string
	ProductFamily string
	Attributes    AmazonEKS_Product_Attributes `gorm:"ForeignKey:AmazonEKS_Product_AttributesID"`
}
type AmazonEKS_Product_Attributes struct {
	gorm.Model
	AmazonEKS_Product_AttributesID uint
	Location                       string
	LocationType                   string
	Usagetype                      string
	Operation                      string
	Servicename                    string
	Tiertype                       string
	Servicecode                    string
}

type AmazonEKS_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonEKSID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonEKS_Term_PriceDimensions `gorm:"ForeignKey:AmazonEKS_TermID"`
	TermAttributes  []*AmazonEKS_Term_Attributes      `gorm:"ForeignKey:AmazonEKS_TermID"`
}

type AmazonEKS_Term_Attributes struct {
	gorm.Model
	AmazonEKS_TermID uint
	Key              string
	Value            string
}

type AmazonEKS_Term_PriceDimensions struct {
	gorm.Model
	AmazonEKS_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AmazonEKS_Term_PricePerUnit `gorm:"ForeignKey:AmazonEKS_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonEKS_Term_PricePerUnit struct {
	gorm.Model
	AmazonEKS_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AmazonEKS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEKS/current/index.json"
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
