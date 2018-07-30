package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAwskms struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]Awskms_Product
	Terms           map[string]map[string]map[string]rawAwskms_Term
}

type rawAwskms_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]Awskms_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *Awskms) UnmarshalJSON(data []byte) error {
	var p rawAwskms
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*Awskms_Product{}
	terms := []*Awskms_Term{}

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
				pDimensions := []*Awskms_Term_PriceDimensions{}
				tAttributes := []*Awskms_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := Awskms_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := Awskms_Term{
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

type Awskms struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*Awskms_Product `gorm:"ForeignKey:AwskmsID"`
	Terms           []*Awskms_Term    `gorm:"ForeignKey:AwskmsID"`
}
type Awskms_Product struct {
	gorm.Model
	AwskmsID      uint
	Attributes    Awskms_Product_Attributes `gorm:"ForeignKey:Awskms_Product_AttributesID"`
	Sku           string
	ProductFamily string
}
type Awskms_Product_Attributes struct {
	gorm.Model
	Awskms_Product_AttributesID uint
	Operation                   string
	Servicename                 string
	Servicecode                 string
	Location                    string
	LocationType                string
	Usagetype                   string
}

type Awskms_Term struct {
	gorm.Model
	OfferTermCode   string
	AwskmsID        uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*Awskms_Term_PriceDimensions `gorm:"ForeignKey:Awskms_TermID"`
	TermAttributes  []*Awskms_Term_Attributes      `gorm:"ForeignKey:Awskms_TermID"`
}

type Awskms_Term_Attributes struct {
	gorm.Model
	Awskms_TermID uint
	Key           string
	Value         string
}

type Awskms_Term_PriceDimensions struct {
	gorm.Model
	Awskms_TermID uint
	RateCode      string
	RateType      string
	Description   string
	BeginRange    string
	EndRange      string
	Unit          string
	PricePerUnit  *Awskms_Term_PricePerUnit `gorm:"ForeignKey:Awskms_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type Awskms_Term_PricePerUnit struct {
	gorm.Model
	Awskms_Term_PriceDimensionsID uint
	USD                           string
}

func (a *Awskms) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/awskms/current/index.json"
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
