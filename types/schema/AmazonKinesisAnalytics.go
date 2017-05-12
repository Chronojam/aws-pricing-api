package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonKinesisAnalytics struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonKinesisAnalytics_Product
	Terms		map[string]map[string]map[string]rawAmazonKinesisAnalytics_Term
}


type rawAmazonKinesisAnalytics_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonKinesisAnalytics_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonKinesisAnalytics) UnmarshalJSON(data []byte) error {
	var p rawAmazonKinesisAnalytics
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonKinesisAnalytics_Product{}
	terms := []AmazonKinesisAnalytics_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonKinesisAnalytics_Term_PriceDimensions{}
				tAttributes := []AmazonKinesisAnalytics_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonKinesisAnalytics_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonKinesisAnalytics_Term{
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

type AmazonKinesisAnalytics struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonKinesisAnalytics_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonKinesisAnalytics_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonKinesisAnalytics_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonKinesisAnalytics_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonKinesisAnalytics_Product_Attributes struct {
	gorm.Model
		Operation	string
	Servicecode	string
	Description	string
	Location	string
	LocationType	string
	Usagetype	string
}

type AmazonKinesisAnalytics_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonKinesisAnalytics_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonKinesisAnalytics_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonKinesisAnalytics_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonKinesisAnalytics_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonKinesisAnalytics_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonKinesisAnalytics_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonKinesisAnalytics) QueryProducts(q func(product AmazonKinesisAnalytics_Product) bool) []AmazonKinesisAnalytics_Product{
	ret := []AmazonKinesisAnalytics_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonKinesisAnalytics) QueryTerms(t string, q func(product AmazonKinesisAnalytics_Term) bool) []AmazonKinesisAnalytics_Term{
	ret := []AmazonKinesisAnalytics_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonKinesisAnalytics) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesisAnalytics/current/index.json"
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