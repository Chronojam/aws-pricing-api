package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonVPC struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonVPC_Product
	Terms		map[string]map[string]map[string]rawAmazonVPC_Term
}


type rawAmazonVPC_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonVPC_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonVPC) UnmarshalJSON(data []byte) error {
	var p rawAmazonVPC
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonVPC_Product{}
	terms := []AmazonVPC_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonVPC_Term_PriceDimensions{}
				tAttributes := []AmazonVPC_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonVPC_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonVPC_Term{
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

type AmazonVPC struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonVPC_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonVPC_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonVPC_Product struct {
	gorm.Model
		ProductFamily	string
	Attributes	AmazonVPC_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
}
type AmazonVPC_Product_Attributes struct {
	gorm.Model
		FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
}

type AmazonVPC_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonVPC_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonVPC_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonVPC_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonVPC_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonVPC_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonVPC_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonVPC) QueryProducts(q func(product AmazonVPC_Product) bool) []AmazonVPC_Product{
	ret := []AmazonVPC_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonVPC) QueryTerms(t string, q func(product AmazonVPC_Term) bool) []AmazonVPC_Term{
	ret := []AmazonVPC_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonVPC) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonVPC/current/index.json"
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