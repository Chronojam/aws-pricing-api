package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawElasticMapReduce struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]ElasticMapReduce_Product
	Terms		map[string]map[string]map[string]rawElasticMapReduce_Term
}


type rawElasticMapReduce_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]ElasticMapReduce_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *ElasticMapReduce) UnmarshalJSON(data []byte) error {
	var p rawElasticMapReduce
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []ElasticMapReduce_Product{}
	terms := []ElasticMapReduce_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []ElasticMapReduce_Term_PriceDimensions{}
				tAttributes := []ElasticMapReduce_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := ElasticMapReduce_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := ElasticMapReduce_Term{
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

type ElasticMapReduce struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]ElasticMapReduce_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]ElasticMapReduce_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type ElasticMapReduce_Product struct {
	gorm.Model
		ProductFamily	string
	Attributes	ElasticMapReduce_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
}
type ElasticMapReduce_Product_Attributes struct {
	gorm.Model
		Usagetype	string
	Operation	string
	SoftwareType	string
	Servicecode	string
	Location	string
	LocationType	string
	InstanceType	string
	InstanceFamily	string
}

type ElasticMapReduce_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []ElasticMapReduce_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []ElasticMapReduce_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type ElasticMapReduce_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type ElasticMapReduce_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	ElasticMapReduce_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type ElasticMapReduce_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a ElasticMapReduce) QueryProducts(q func(product ElasticMapReduce_Product) bool) []ElasticMapReduce_Product{
	ret := []ElasticMapReduce_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a ElasticMapReduce) QueryTerms(t string, q func(product ElasticMapReduce_Term) bool) []ElasticMapReduce_Term{
	ret := []ElasticMapReduce_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *ElasticMapReduce) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/ElasticMapReduce/current/index.json"
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