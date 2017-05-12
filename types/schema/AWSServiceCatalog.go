package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSServiceCatalog struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSServiceCatalog_Product
	Terms		map[string]map[string]map[string]rawAWSServiceCatalog_Term
}


type rawAWSServiceCatalog_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSServiceCatalog_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSServiceCatalog) UnmarshalJSON(data []byte) error {
	var p rawAWSServiceCatalog
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSServiceCatalog_Product{}
	terms := []AWSServiceCatalog_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSServiceCatalog_Term_PriceDimensions{}
				tAttributes := []AWSServiceCatalog_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSServiceCatalog_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSServiceCatalog_Term{
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

type AWSServiceCatalog struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AWSServiceCatalog_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSServiceCatalog_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSServiceCatalog_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AWSServiceCatalog_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSServiceCatalog_Product_Attributes struct {
	gorm.Model
		Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	WithActiveUsers	string
}

type AWSServiceCatalog_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSServiceCatalog_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSServiceCatalog_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSServiceCatalog_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSServiceCatalog_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSServiceCatalog_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSServiceCatalog_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSServiceCatalog) QueryProducts(q func(product AWSServiceCatalog_Product) bool) []AWSServiceCatalog_Product{
	ret := []AWSServiceCatalog_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSServiceCatalog) QueryTerms(t string, q func(product AWSServiceCatalog_Term) bool) []AWSServiceCatalog_Term{
	ret := []AWSServiceCatalog_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AWSServiceCatalog) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSServiceCatalog/current/index.json"
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