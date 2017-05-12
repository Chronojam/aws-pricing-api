package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSBudgets struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSBudgets_Product
	Terms		map[string]map[string]map[string]rawAWSBudgets_Term
}


type rawAWSBudgets_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSBudgets_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSBudgets) UnmarshalJSON(data []byte) error {
	var p rawAWSBudgets
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSBudgets_Product{}
	terms := []AWSBudgets_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSBudgets_Term_PriceDimensions{}
				tAttributes := []AWSBudgets_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSBudgets_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSBudgets_Term{
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

type AWSBudgets struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AWSBudgets_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSBudgets_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSBudgets_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AWSBudgets_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSBudgets_Product_Attributes struct {
	gorm.Model
		Servicecode	string
	Location	string
	LocationType	string
	GroupDescription	string
	Usagetype	string
	Operation	string
}

type AWSBudgets_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSBudgets_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSBudgets_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSBudgets_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSBudgets_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSBudgets_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSBudgets_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSBudgets) QueryProducts(q func(product AWSBudgets_Product) bool) []AWSBudgets_Product{
	ret := []AWSBudgets_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSBudgets) QueryTerms(t string, q func(product AWSBudgets_Term) bool) []AWSBudgets_Term{
	ret := []AWSBudgets_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AWSBudgets) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSBudgets/current/index.json"
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