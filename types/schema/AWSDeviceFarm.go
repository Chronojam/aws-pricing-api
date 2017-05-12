package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSDeviceFarm struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDeviceFarm_Product
	Terms		map[string]map[string]map[string]rawAWSDeviceFarm_Term
}


type rawAWSDeviceFarm_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSDeviceFarm_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSDeviceFarm) UnmarshalJSON(data []byte) error {
	var p rawAWSDeviceFarm
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSDeviceFarm_Product{}
	terms := []AWSDeviceFarm_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSDeviceFarm_Term_PriceDimensions{}
				tAttributes := []AWSDeviceFarm_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSDeviceFarm_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSDeviceFarm_Term{
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

type AWSDeviceFarm struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AWSDeviceFarm_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSDeviceFarm_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSDeviceFarm_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AWSDeviceFarm_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSDeviceFarm_Product_Attributes struct {
	gorm.Model
		Location	string
	LocationType	string
	Usagetype	string
	DeviceOs	string
	Servicecode	string
	Description	string
	Operation	string
	ExecutionMode	string
	MeterMode	string
}

type AWSDeviceFarm_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSDeviceFarm_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSDeviceFarm_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSDeviceFarm_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSDeviceFarm_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDeviceFarm_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSDeviceFarm_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSDeviceFarm) QueryProducts(q func(product AWSDeviceFarm_Product) bool) []AWSDeviceFarm_Product{
	ret := []AWSDeviceFarm_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSDeviceFarm) QueryTerms(t string, q func(product AWSDeviceFarm_Term) bool) []AWSDeviceFarm_Term{
	ret := []AWSDeviceFarm_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AWSDeviceFarm) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDeviceFarm/current/index.json"
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