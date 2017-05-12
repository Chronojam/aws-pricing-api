package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonApiGateway struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonApiGateway_Product
	Terms		map[string]map[string]map[string]rawAmazonApiGateway_Term
}


type rawAmazonApiGateway_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonApiGateway_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonApiGateway) UnmarshalJSON(data []byte) error {
	var p rawAmazonApiGateway
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonApiGateway_Product{}
	terms := []AmazonApiGateway_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonApiGateway_Term_PriceDimensions{}
				tAttributes := []AmazonApiGateway_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonApiGateway_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonApiGateway_Term{
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

type AmazonApiGateway struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonApiGateway_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonApiGateway_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonApiGateway_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonApiGateway_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonApiGateway_Product_Attributes struct {
	gorm.Model
		Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
}

type AmazonApiGateway_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonApiGateway_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonApiGateway_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonApiGateway_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonApiGateway_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonApiGateway_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonApiGateway_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonApiGateway) QueryProducts(q func(product AmazonApiGateway_Product) bool) []AmazonApiGateway_Product{
	ret := []AmazonApiGateway_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonApiGateway) QueryTerms(t string, q func(product AmazonApiGateway_Term) bool) []AmazonApiGateway_Term{
	ret := []AmazonApiGateway_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonApiGateway) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonApiGateway/current/index.json"
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