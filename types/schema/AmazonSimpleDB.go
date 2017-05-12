package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonSimpleDB struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonSimpleDB_Product
	Terms		map[string]map[string]map[string]rawAmazonSimpleDB_Term
}


type rawAmazonSimpleDB_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonSimpleDB_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonSimpleDB) UnmarshalJSON(data []byte) error {
	var p rawAmazonSimpleDB
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonSimpleDB_Product{}
	terms := []AmazonSimpleDB_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonSimpleDB_Term_PriceDimensions{}
				tAttributes := []AmazonSimpleDB_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonSimpleDB_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonSimpleDB_Term{
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

type AmazonSimpleDB struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonSimpleDB_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonSimpleDB_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonSimpleDB_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonSimpleDB_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonSimpleDB_Product_Attributes struct {
	gorm.Model
		Servicecode	string
	Location	string
	LocationType	string
	VolumeType	string
	Usagetype	string
	Operation	string
}

type AmazonSimpleDB_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonSimpleDB_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonSimpleDB_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonSimpleDB_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonSimpleDB_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonSimpleDB_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonSimpleDB_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonSimpleDB) QueryProducts(q func(product AmazonSimpleDB_Product) bool) []AmazonSimpleDB_Product{
	ret := []AmazonSimpleDB_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonSimpleDB) QueryTerms(t string, q func(product AmazonSimpleDB_Term) bool) []AmazonSimpleDB_Term{
	ret := []AmazonSimpleDB_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonSimpleDB) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonSimpleDB/current/index.json"
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