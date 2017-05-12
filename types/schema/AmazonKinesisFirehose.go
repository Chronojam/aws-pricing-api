package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonKinesisFirehose struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonKinesisFirehose_Product
	Terms		map[string]map[string]map[string]rawAmazonKinesisFirehose_Term
}


type rawAmazonKinesisFirehose_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonKinesisFirehose_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonKinesisFirehose) UnmarshalJSON(data []byte) error {
	var p rawAmazonKinesisFirehose
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonKinesisFirehose_Product{}
	terms := []AmazonKinesisFirehose_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonKinesisFirehose_Term_PriceDimensions{}
				tAttributes := []AmazonKinesisFirehose_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonKinesisFirehose_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonKinesisFirehose_Term{
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

type AmazonKinesisFirehose struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonKinesisFirehose_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonKinesisFirehose_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonKinesisFirehose_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonKinesisFirehose_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonKinesisFirehose_Product_Attributes struct {
	gorm.Model
		Usagetype	string
	Operation	string
	Servicecode	string
	Description	string
	Location	string
	LocationType	string
	Group	string
}

type AmazonKinesisFirehose_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonKinesisFirehose_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonKinesisFirehose_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonKinesisFirehose_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonKinesisFirehose_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonKinesisFirehose_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonKinesisFirehose_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonKinesisFirehose) QueryProducts(q func(product AmazonKinesisFirehose_Product) bool) []AmazonKinesisFirehose_Product{
	ret := []AmazonKinesisFirehose_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonKinesisFirehose) QueryTerms(t string, q func(product AmazonKinesisFirehose_Term) bool) []AmazonKinesisFirehose_Term{
	ret := []AmazonKinesisFirehose_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonKinesisFirehose) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonKinesisFirehose/current/index.json"
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