package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAwswaf struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]Awswaf_Product
	Terms		map[string]map[string]map[string]rawAwswaf_Term
}


type rawAwswaf_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]Awswaf_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *Awswaf) UnmarshalJSON(data []byte) error {
	var p rawAwswaf
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*Awswaf_Product{}
	terms := []*Awswaf_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*Awswaf_Term_PriceDimensions{}
				tAttributes := []*Awswaf_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := Awswaf_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := Awswaf_Term{
					OfferTermCode: term.OfferTermCode,
					Sku: term.Sku,
					EffectiveDate: term.EffectiveDate,
					TermAttributes: tAttributes,
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

type Awswaf struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*Awswaf_Product `gorm:"ForeignKey:AwswafID"`
	Terms		[]*Awswaf_Term`gorm:"ForeignKey:AwswafID"`
}
type Awswaf_Product struct {
	gorm.Model
		AwswafID	uint
	Sku	string
	ProductFamily	string
	Attributes	Awswaf_Product_Attributes	`gorm:"ForeignKey:Awswaf_Product_AttributesID"`
}
type Awswaf_Product_Attributes struct {
	gorm.Model
		Awswaf_Product_AttributesID	uint
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
}

type Awswaf_Term struct {
	gorm.Model
	OfferTermCode string
	AwswafID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*Awswaf_Term_PriceDimensions `gorm:"ForeignKey:Awswaf_TermID"`
	TermAttributes []*Awswaf_Term_Attributes `gorm:"ForeignKey:Awswaf_TermID"`
}

type Awswaf_Term_Attributes struct {
	gorm.Model
	Awswaf_TermID	uint
	Key	string
	Value	string
}

type Awswaf_Term_PriceDimensions struct {
	gorm.Model
	Awswaf_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*Awswaf_Term_PricePerUnit `gorm:"ForeignKey:Awswaf_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type Awswaf_Term_PricePerUnit struct {
	gorm.Model
	Awswaf_Term_PriceDimensionsID	uint
	USD	string
}
func (a *Awswaf) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/awswaf/current/index.json"
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