package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonEFS struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonEFS_Product
	Terms		map[string]map[string]map[string]rawAmazonEFS_Term
}


type rawAmazonEFS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonEFS_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonEFS) UnmarshalJSON(data []byte) error {
	var p rawAmazonEFS
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonEFS_Product{}
	terms := []AmazonEFS_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonEFS_Term_PriceDimensions{}
				tAttributes := []AmazonEFS_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonEFS_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonEFS_Term{
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

type AmazonEFS struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonEFS_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonEFS_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonEFS_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonEFS_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonEFS_Product_Attributes struct {
	gorm.Model
		Servicecode	string
	Location	string
	LocationType	string
	StorageClass	string
	Usagetype	string
	Operation	string
}

type AmazonEFS_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonEFS_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonEFS_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonEFS_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonEFS_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonEFS_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonEFS_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonEFS) QueryProducts(q func(product AmazonEFS_Product) bool) []AmazonEFS_Product{
	ret := []AmazonEFS_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonEFS) QueryTerms(t string, q func(product AmazonEFS_Term) bool) []AmazonEFS_Term{
	ret := []AmazonEFS_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonEFS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEFS/current/index.json"
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