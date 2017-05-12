package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonCloudWatch struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCloudWatch_Product
	Terms		map[string]map[string]map[string]rawAmazonCloudWatch_Term
}


type rawAmazonCloudWatch_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonCloudWatch_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonCloudWatch) UnmarshalJSON(data []byte) error {
	var p rawAmazonCloudWatch
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonCloudWatch_Product{}
	terms := []AmazonCloudWatch_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonCloudWatch_Term_PriceDimensions{}
				tAttributes := []AmazonCloudWatch_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonCloudWatch_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonCloudWatch_Term{
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

type AmazonCloudWatch struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonCloudWatch_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonCloudWatch_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonCloudWatch_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonCloudWatch_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonCloudWatch_Product_Attributes struct {
	gorm.Model
		Location	string
	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AmazonCloudWatch_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonCloudWatch_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonCloudWatch_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonCloudWatch_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonCloudWatch_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCloudWatch_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonCloudWatch_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonCloudWatch) QueryProducts(q func(product AmazonCloudWatch_Product) bool) []AmazonCloudWatch_Product{
	ret := []AmazonCloudWatch_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonCloudWatch) QueryTerms(t string, q func(product AmazonCloudWatch_Term) bool) []AmazonCloudWatch_Term{
	ret := []AmazonCloudWatch_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonCloudWatch) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudWatch/current/index.json"
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