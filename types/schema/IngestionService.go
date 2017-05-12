package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawIngestionService struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]IngestionService_Product
	Terms		map[string]map[string]map[string]rawIngestionService_Term
}


type rawIngestionService_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]IngestionService_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *IngestionService) UnmarshalJSON(data []byte) error {
	var p rawIngestionService
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []IngestionService_Product{}
	terms := []IngestionService_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []IngestionService_Term_PriceDimensions{}
				tAttributes := []IngestionService_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := IngestionService_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := IngestionService_Term{
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

type IngestionService struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]IngestionService_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]IngestionService_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type IngestionService_Product struct {
	gorm.Model
		Attributes	IngestionService_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
	ProductFamily	string
}
type IngestionService_Product_Attributes struct {
	gorm.Model
		Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
	DataAction	string
	Servicecode	string
	Location	string
	LocationType	string
}

type IngestionService_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []IngestionService_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []IngestionService_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type IngestionService_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type IngestionService_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	IngestionService_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type IngestionService_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a IngestionService) QueryProducts(q func(product IngestionService_Product) bool) []IngestionService_Product{
	ret := []IngestionService_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a IngestionService) QueryTerms(t string, q func(product IngestionService_Term) bool) []IngestionService_Term{
	ret := []IngestionService_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *IngestionService) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/IngestionService/current/index.json"
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