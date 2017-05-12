package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonCognitoSync struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCognitoSync_Product
	Terms		map[string]map[string]map[string]rawAmazonCognitoSync_Term
}


type rawAmazonCognitoSync_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonCognitoSync_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonCognitoSync) UnmarshalJSON(data []byte) error {
	var p rawAmazonCognitoSync
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonCognitoSync_Product{}
	terms := []AmazonCognitoSync_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonCognitoSync_Term_PriceDimensions{}
				tAttributes := []AmazonCognitoSync_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonCognitoSync_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonCognitoSync_Term{
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

type AmazonCognitoSync struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonCognitoSync_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonCognitoSync_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonCognitoSync_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonCognitoSync_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonCognitoSync_Product_Attributes struct {
	gorm.Model
		Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
}

type AmazonCognitoSync_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonCognitoSync_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonCognitoSync_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonCognitoSync_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonCognitoSync_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCognitoSync_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonCognitoSync_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonCognitoSync) QueryProducts(q func(product AmazonCognitoSync_Product) bool) []AmazonCognitoSync_Product{
	ret := []AmazonCognitoSync_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonCognitoSync) QueryTerms(t string, q func(product AmazonCognitoSync_Term) bool) []AmazonCognitoSync_Term{
	ret := []AmazonCognitoSync_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonCognitoSync) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCognitoSync/current/index.json"
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