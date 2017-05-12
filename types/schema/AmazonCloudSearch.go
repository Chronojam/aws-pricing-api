package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonCloudSearch struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCloudSearch_Product
	Terms		map[string]map[string]map[string]rawAmazonCloudSearch_Term
}


type rawAmazonCloudSearch_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonCloudSearch_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonCloudSearch) UnmarshalJSON(data []byte) error {
	var p rawAmazonCloudSearch
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonCloudSearch_Product{}
	terms := []AmazonCloudSearch_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonCloudSearch_Term_PriceDimensions{}
				tAttributes := []AmazonCloudSearch_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonCloudSearch_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonCloudSearch_Term{
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

type AmazonCloudSearch struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonCloudSearch_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonCloudSearch_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonCloudSearch_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonCloudSearch_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonCloudSearch_Product_Attributes struct {
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

type AmazonCloudSearch_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonCloudSearch_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonCloudSearch_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonCloudSearch_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonCloudSearch_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonCloudSearch_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonCloudSearch_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonCloudSearch) QueryProducts(q func(product AmazonCloudSearch_Product) bool) []AmazonCloudSearch_Product{
	ret := []AmazonCloudSearch_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonCloudSearch) QueryTerms(t string, q func(product AmazonCloudSearch_Term) bool) []AmazonCloudSearch_Term{
	ret := []AmazonCloudSearch_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonCloudSearch) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudSearch/current/index.json"
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