package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSDirectoryService struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDirectoryService_Product
	Terms		map[string]map[string]map[string]rawAWSDirectoryService_Term
}


type rawAWSDirectoryService_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSDirectoryService_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSDirectoryService) UnmarshalJSON(data []byte) error {
	var p rawAWSDirectoryService
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSDirectoryService_Product{}
	terms := []AWSDirectoryService_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSDirectoryService_Term_PriceDimensions{}
				tAttributes := []AWSDirectoryService_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSDirectoryService_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSDirectoryService_Term{
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

type AWSDirectoryService struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AWSDirectoryService_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSDirectoryService_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSDirectoryService_Product struct {
	gorm.Model
		ProductFamily	string
	Attributes	AWSDirectoryService_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
}
type AWSDirectoryService_Product_Attributes struct {
	gorm.Model
		LocationType	string
	Usagetype	string
	Operation	string
	DirectorySize	string
	DirectoryType	string
	DirectoryTypeDescription	string
	Servicecode	string
	Location	string
}

type AWSDirectoryService_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSDirectoryService_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSDirectoryService_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSDirectoryService_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSDirectoryService_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDirectoryService_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSDirectoryService_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSDirectoryService) QueryProducts(q func(product AWSDirectoryService_Product) bool) []AWSDirectoryService_Product{
	ret := []AWSDirectoryService_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSDirectoryService) QueryTerms(t string, q func(product AWSDirectoryService_Term) bool) []AWSDirectoryService_Term{
	ret := []AWSDirectoryService_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AWSDirectoryService) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDirectoryService/current/index.json"
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