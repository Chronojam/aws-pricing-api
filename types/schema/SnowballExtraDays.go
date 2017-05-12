package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawSnowballExtraDays struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]SnowballExtraDays_Product
	Terms		map[string]map[string]map[string]rawSnowballExtraDays_Term
}


type rawSnowballExtraDays_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]SnowballExtraDays_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *SnowballExtraDays) UnmarshalJSON(data []byte) error {
	var p rawSnowballExtraDays
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []SnowballExtraDays_Product{}
	terms := []SnowballExtraDays_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []SnowballExtraDays_Term_PriceDimensions{}
				tAttributes := []SnowballExtraDays_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := SnowballExtraDays_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := SnowballExtraDays_Term{
					OfferTermCode: term.OfferTermCode,
					Sku: term.Sku,
					EffectiveDate: term.EffectiveDate,
					TermAttributes: tAttributes,
					PriceDimensions: pDimensions,
				}

				terms = append(terms, t)
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

type SnowballExtraDays struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]SnowballExtraDays_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]SnowballExtraDays_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type SnowballExtraDays_Product struct {
	gorm.Model
		Attributes	SnowballExtraDays_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
	ProductFamily	string
}
type SnowballExtraDays_Product_Attributes struct {
	gorm.Model
		FeeCode	string
	FeeDescription	string
	Usagetype	string
	Operation	string
	SnowballType	string
	Servicecode	string
	Location	string
	LocationType	string
}

type SnowballExtraDays_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []SnowballExtraDays_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []SnowballExtraDays_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type SnowballExtraDays_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type SnowballExtraDays_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	SnowballExtraDays_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type SnowballExtraDays_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a SnowballExtraDays) QueryProducts(q func(product SnowballExtraDays_Product) bool) []SnowballExtraDays_Product{
	ret := []SnowballExtraDays_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a SnowballExtraDays) QueryTerms(t string, q func(product SnowballExtraDays_Term) bool) []SnowballExtraDays_Term{
	ret := []SnowballExtraDays_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *SnowballExtraDays) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/SnowballExtraDays/current/index.json"
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