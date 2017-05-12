package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonWorkMail struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWorkMail_Product
	Terms		map[string]map[string]map[string]rawAmazonWorkMail_Term
}


type rawAmazonWorkMail_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonWorkMail_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonWorkMail) UnmarshalJSON(data []byte) error {
	var p rawAmazonWorkMail
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonWorkMail_Product{}
	terms := []AmazonWorkMail_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonWorkMail_Term_PriceDimensions{}
				tAttributes := []AmazonWorkMail_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonWorkMail_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonWorkMail_Term{
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

type AmazonWorkMail struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonWorkMail_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonWorkMail_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonWorkMail_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonWorkMail_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonWorkMail_Product_Attributes struct {
	gorm.Model
		MailboxStorage	string
	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	FreeTier	string
}

type AmazonWorkMail_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonWorkMail_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonWorkMail_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonWorkMail_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonWorkMail_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonWorkMail_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonWorkMail_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonWorkMail) QueryProducts(q func(product AmazonWorkMail_Product) bool) []AmazonWorkMail_Product{
	ret := []AmazonWorkMail_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonWorkMail) QueryTerms(t string, q func(product AmazonWorkMail_Term) bool) []AmazonWorkMail_Term{
	ret := []AmazonWorkMail_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonWorkMail) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkMail/current/index.json"
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