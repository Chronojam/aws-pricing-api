package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonAthena struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonAthena_Product
	Terms		map[string]map[string]map[string]rawAmazonAthena_Term
}


type rawAmazonAthena_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonAthena_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonAthena) UnmarshalJSON(data []byte) error {
	var p rawAmazonAthena
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonAthena_Product{}
	terms := []AmazonAthena_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonAthena_Term_PriceDimensions{}
				tAttributes := []AmazonAthena_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonAthena_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonAthena_Term{
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

type AmazonAthena struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonAthena_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonAthena_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonAthena_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonAthena_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonAthena_Product_Attributes struct {
	gorm.Model
		Servicecode	string
	Description	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	FreeQueryTypes	string
}

type AmazonAthena_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonAthena_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonAthena_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonAthena_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonAthena_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonAthena_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonAthena_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonAthena) QueryProducts(q func(product AmazonAthena_Product) bool) []AmazonAthena_Product{
	ret := []AmazonAthena_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonAthena) QueryTerms(t string, q func(product AmazonAthena_Term) bool) []AmazonAthena_Term{
	ret := []AmazonAthena_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonAthena) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonAthena/current/index.json"
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