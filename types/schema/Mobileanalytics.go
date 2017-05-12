package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawMobileanalytics struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]Mobileanalytics_Product
	Terms		map[string]map[string]map[string]rawMobileanalytics_Term
}


type rawMobileanalytics_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]Mobileanalytics_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *Mobileanalytics) UnmarshalJSON(data []byte) error {
	var p rawMobileanalytics
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []Mobileanalytics_Product{}
	terms := []Mobileanalytics_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []Mobileanalytics_Term_PriceDimensions{}
				tAttributes := []Mobileanalytics_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := Mobileanalytics_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := Mobileanalytics_Term{
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

type Mobileanalytics struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]Mobileanalytics_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]Mobileanalytics_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type Mobileanalytics_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	Mobileanalytics_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type Mobileanalytics_Product_Attributes struct {
	gorm.Model
		Servicecode	string
	Description	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	IncludedEvents	string
}

type Mobileanalytics_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []Mobileanalytics_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []Mobileanalytics_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type Mobileanalytics_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type Mobileanalytics_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	Mobileanalytics_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type Mobileanalytics_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a Mobileanalytics) QueryProducts(q func(product Mobileanalytics_Product) bool) []Mobileanalytics_Product{
	ret := []Mobileanalytics_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a Mobileanalytics) QueryTerms(t string, q func(product Mobileanalytics_Term) bool) []Mobileanalytics_Term{
	ret := []Mobileanalytics_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *Mobileanalytics) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/mobileanalytics/current/index.json"
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