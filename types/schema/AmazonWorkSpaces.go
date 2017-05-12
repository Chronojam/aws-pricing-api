package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonWorkSpaces struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWorkSpaces_Product
	Terms		map[string]map[string]map[string]rawAmazonWorkSpaces_Term
}


type rawAmazonWorkSpaces_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonWorkSpaces_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonWorkSpaces) UnmarshalJSON(data []byte) error {
	var p rawAmazonWorkSpaces
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonWorkSpaces_Product{}
	terms := []AmazonWorkSpaces_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonWorkSpaces_Term_PriceDimensions{}
				tAttributes := []AmazonWorkSpaces_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonWorkSpaces_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonWorkSpaces_Term{
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

type AmazonWorkSpaces struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonWorkSpaces_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonWorkSpaces_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonWorkSpaces_Product struct {
	gorm.Model
		Attributes	AmazonWorkSpaces_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
	ProductFamily	string
}
type AmazonWorkSpaces_Product_Attributes struct {
	gorm.Model
		Location	string
	Memory	string
	SoftwareIncluded	string
	LocationType	string
	GroupDescription	string
	Bundle	string
	Operation	string
	ResourceType	string
	License	string
	RunningMode	string
	Servicecode	string
	Vcpu	string
	Storage	string
	Group	string
	Usagetype	string
}

type AmazonWorkSpaces_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonWorkSpaces_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonWorkSpaces_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonWorkSpaces_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonWorkSpaces_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonWorkSpaces_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonWorkSpaces_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonWorkSpaces) QueryProducts(q func(product AmazonWorkSpaces_Product) bool) []AmazonWorkSpaces_Product{
	ret := []AmazonWorkSpaces_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonWorkSpaces) QueryTerms(t string, q func(product AmazonWorkSpaces_Term) bool) []AmazonWorkSpaces_Term{
	ret := []AmazonWorkSpaces_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonWorkSpaces) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkSpaces/current/index.json"
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