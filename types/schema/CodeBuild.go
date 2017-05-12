package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawCodeBuild struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]CodeBuild_Product
	Terms		map[string]map[string]map[string]rawCodeBuild_Term
}


type rawCodeBuild_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]CodeBuild_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *CodeBuild) UnmarshalJSON(data []byte) error {
	var p rawCodeBuild
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []CodeBuild_Product{}
	terms := []CodeBuild_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []CodeBuild_Term_PriceDimensions{}
				tAttributes := []CodeBuild_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := CodeBuild_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := CodeBuild_Term{
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

type CodeBuild struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]CodeBuild_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]CodeBuild_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type CodeBuild_Product struct {
	gorm.Model
		Attributes	CodeBuild_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
	ProductFamily	string
}
type CodeBuild_Product_Attributes struct {
	gorm.Model
		ComputeFamily	string
	Vcpu	string
	OperatingSystem	string
	LocationType	string
	Memory	string
	Usagetype	string
	Operation	string
	ComputeType	string
	Servicecode	string
	Location	string
}

type CodeBuild_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []CodeBuild_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []CodeBuild_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type CodeBuild_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type CodeBuild_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	CodeBuild_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type CodeBuild_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a CodeBuild) QueryProducts(q func(product CodeBuild_Product) bool) []CodeBuild_Product{
	ret := []CodeBuild_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a CodeBuild) QueryTerms(t string, q func(product CodeBuild_Term) bool) []CodeBuild_Term{
	ret := []CodeBuild_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *CodeBuild) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/CodeBuild/current/index.json"
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