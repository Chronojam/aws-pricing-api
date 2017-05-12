package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAlexaTopSites struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AlexaTopSites_Product
	Terms		map[string]map[string]map[string]rawAlexaTopSites_Term
}


type rawAlexaTopSites_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AlexaTopSites_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AlexaTopSites) UnmarshalJSON(data []byte) error {
	var p rawAlexaTopSites
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AlexaTopSites_Product{}
	terms := []AlexaTopSites_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AlexaTopSites_Term_PriceDimensions{}
				tAttributes := []AlexaTopSites_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AlexaTopSites_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AlexaTopSites_Term{
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

type AlexaTopSites struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AlexaTopSites_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AlexaTopSites_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AlexaTopSites_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AlexaTopSites_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AlexaTopSites_Product_Attributes struct {
	gorm.Model
		Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
}

type AlexaTopSites_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AlexaTopSites_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AlexaTopSites_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AlexaTopSites_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AlexaTopSites_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AlexaTopSites_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AlexaTopSites_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AlexaTopSites) QueryProducts(q func(product AlexaTopSites_Product) bool) []AlexaTopSites_Product{
	ret := []AlexaTopSites_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AlexaTopSites) QueryTerms(t string, q func(product AlexaTopSites_Term) bool) []AlexaTopSites_Term{
	ret := []AlexaTopSites_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AlexaTopSites) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AlexaTopSites/current/index.json"
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