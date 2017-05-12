package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonLex struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonLex_Product
	Terms		map[string]map[string]map[string]rawAmazonLex_Term
}


type rawAmazonLex_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonLex_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonLex) UnmarshalJSON(data []byte) error {
	var p rawAmazonLex
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonLex_Product{}
	terms := []AmazonLex_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonLex_Term_PriceDimensions{}
				tAttributes := []AmazonLex_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonLex_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonLex_Term{
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

type AmazonLex struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonLex_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonLex_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonLex_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonLex_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonLex_Product_Attributes struct {
	gorm.Model
		GroupDescription	string
	Usagetype	string
	InputMode	string
	OutputMode	string
	SupportedModes	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	Operation	string
}

type AmazonLex_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonLex_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonLex_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonLex_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonLex_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonLex_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonLex_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonLex) QueryProducts(q func(product AmazonLex_Product) bool) []AmazonLex_Product{
	ret := []AmazonLex_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonLex) QueryTerms(t string, q func(product AmazonLex_Term) bool) []AmazonLex_Term{
	ret := []AmazonLex_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonLex) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonLex/current/index.json"
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