package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAwskms struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]Awskms_Product
	Terms		map[string]map[string]map[string]rawAwskms_Term
}


type rawAwskms_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]Awskms_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *Awskms) UnmarshalJSON(data []byte) error {
	var p rawAwskms
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []Awskms_Product{}
	terms := []Awskms_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []Awskms_Term_PriceDimensions{}
				tAttributes := []Awskms_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := Awskms_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := Awskms_Term{
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

type Awskms struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]Awskms_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]Awskms_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type Awskms_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	Awskms_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type Awskms_Product_Attributes struct {
	gorm.Model
		Usagetype	string
	Operation	string
	Servicecode	string
	Location	string
	LocationType	string
}

type Awskms_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []Awskms_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []Awskms_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type Awskms_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type Awskms_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	Awskms_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type Awskms_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a Awskms) QueryProducts(q func(product Awskms_Product) bool) []Awskms_Product{
	ret := []Awskms_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a Awskms) QueryTerms(t string, q func(product Awskms_Term) bool) []Awskms_Term{
	ret := []Awskms_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
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