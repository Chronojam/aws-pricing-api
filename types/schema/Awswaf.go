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
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []Awswaf_Product{}
	terms := []Awswaf_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []Awswaf_Term_PriceDimensions{}
				tAttributes := []Awswaf_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := Awswaf_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := Awswaf_Term{
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

type Awswaf struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]Awswaf_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]Awswaf_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type Awswaf_Product struct {
	gorm.Model
		Attributes	Awswaf_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
	ProductFamily	string
}
type Awswaf_Product_Attributes struct {
	gorm.Model
		Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
	Servicecode	string
	Location	string
	LocationType	string
}

type Awswaf_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []Awswaf_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []Awswaf_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type Awswaf_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type Awswaf_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	Awswaf_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type Awswaf_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a Awswaf) QueryProducts(q func(product Awswaf_Product) bool) []Awswaf_Product{
	ret := []Awswaf_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a Awswaf) QueryTerms(t string, q func(product Awswaf_Term) bool) []Awswaf_Term{
	ret := []Awswaf_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
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