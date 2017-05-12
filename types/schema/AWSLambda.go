package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSLambda struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSLambda_Product
	Terms		map[string]map[string]map[string]rawAWSLambda_Term
}


type rawAWSLambda_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSLambda_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSLambda) UnmarshalJSON(data []byte) error {
	var p rawAWSLambda
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSLambda_Product{}
	terms := []AWSLambda_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSLambda_Term_PriceDimensions{}
				tAttributes := []AWSLambda_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSLambda_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSLambda_Term{
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

type AWSLambda struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AWSLambda_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSLambda_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSLambda_Product struct {
	gorm.Model
		Attributes	AWSLambda_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
	ProductFamily	string
}
type AWSLambda_Product_Attributes struct {
	gorm.Model
		FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
}

type AWSLambda_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSLambda_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSLambda_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSLambda_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSLambda_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSLambda_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSLambda_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSLambda) QueryProducts(q func(product AWSLambda_Product) bool) []AWSLambda_Product{
	ret := []AWSLambda_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSLambda) QueryTerms(t string, q func(product AWSLambda_Term) bool) []AWSLambda_Term{
	ret := []AWSLambda_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AWSLambda) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSLambda/current/index.json"
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