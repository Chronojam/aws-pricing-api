package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonWorkDocs struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonWorkDocs_Product
	Terms		map[string]map[string]map[string]rawAmazonWorkDocs_Term
}


type rawAmazonWorkDocs_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonWorkDocs_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonWorkDocs) UnmarshalJSON(data []byte) error {
	var p rawAmazonWorkDocs
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonWorkDocs_Product{}
	terms := []AmazonWorkDocs_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonWorkDocs_Term_PriceDimensions{}
				tAttributes := []AmazonWorkDocs_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonWorkDocs_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonWorkDocs_Term{
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

type AmazonWorkDocs struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonWorkDocs_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonWorkDocs_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonWorkDocs_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonWorkDocs_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonWorkDocs_Product_Attributes struct {
	gorm.Model
		Storage	string
	Usagetype	string
	Operation	string
	Servicecode	string
	Description	string
	Location	string
	LocationType	string
	FreeTrial	string
	MaximumStorageVolume	string
	MinimumStorageVolume	string
}

type AmazonWorkDocs_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonWorkDocs_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonWorkDocs_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonWorkDocs_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonWorkDocs_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonWorkDocs_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonWorkDocs_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonWorkDocs) QueryProducts(q func(product AmazonWorkDocs_Product) bool) []AmazonWorkDocs_Product{
	ret := []AmazonWorkDocs_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonWorkDocs) QueryTerms(t string, q func(product AmazonWorkDocs_Term) bool) []AmazonWorkDocs_Term{
	ret := []AmazonWorkDocs_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonWorkDocs) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonWorkDocs/current/index.json"
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