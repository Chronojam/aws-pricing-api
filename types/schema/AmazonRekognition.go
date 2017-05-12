package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonRekognition struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonRekognition_Product
	Terms		map[string]map[string]map[string]rawAmazonRekognition_Term
}


type rawAmazonRekognition_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonRekognition_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonRekognition) UnmarshalJSON(data []byte) error {
	var p rawAmazonRekognition
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonRekognition_Product{}
	terms := []AmazonRekognition_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonRekognition_Term_PriceDimensions{}
				tAttributes := []AmazonRekognition_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonRekognition_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonRekognition_Term{
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

type AmazonRekognition struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonRekognition_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonRekognition_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonRekognition_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonRekognition_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonRekognition_Product_Attributes struct {
	gorm.Model
		LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
	Servicecode	string
	Location	string
}

type AmazonRekognition_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonRekognition_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonRekognition_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonRekognition_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonRekognition_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonRekognition_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonRekognition_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonRekognition) QueryProducts(q func(product AmazonRekognition_Product) bool) []AmazonRekognition_Product{
	ret := []AmazonRekognition_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonRekognition) QueryTerms(t string, q func(product AmazonRekognition_Term) bool) []AmazonRekognition_Term{
	ret := []AmazonRekognition_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonRekognition) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRekognition/current/index.json"
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