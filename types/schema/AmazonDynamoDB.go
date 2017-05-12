package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonDynamoDB struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonDynamoDB_Product
	Terms		map[string]map[string]map[string]rawAmazonDynamoDB_Term
}


type rawAmazonDynamoDB_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonDynamoDB_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonDynamoDB) UnmarshalJSON(data []byte) error {
	var p rawAmazonDynamoDB
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonDynamoDB_Product{}
	terms := []AmazonDynamoDB_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonDynamoDB_Term_PriceDimensions{}
				tAttributes := []AmazonDynamoDB_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonDynamoDB_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonDynamoDB_Term{
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

type AmazonDynamoDB struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonDynamoDB_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonDynamoDB_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonDynamoDB_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonDynamoDB_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonDynamoDB_Product_Attributes struct {
	gorm.Model
		Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
}

type AmazonDynamoDB_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonDynamoDB_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonDynamoDB_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonDynamoDB_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonDynamoDB_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonDynamoDB_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonDynamoDB_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonDynamoDB) QueryProducts(q func(product AmazonDynamoDB_Product) bool) []AmazonDynamoDB_Product{
	ret := []AmazonDynamoDB_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonDynamoDB) QueryTerms(t string, q func(product AmazonDynamoDB_Term) bool) []AmazonDynamoDB_Term{
	ret := []AmazonDynamoDB_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonDynamoDB) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonDynamoDB/current/index.json"
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