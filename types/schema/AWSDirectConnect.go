package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSDirectConnect struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDirectConnect_Product
	Terms		map[string]map[string]map[string]rawAWSDirectConnect_Term
}


type rawAWSDirectConnect_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSDirectConnect_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSDirectConnect) UnmarshalJSON(data []byte) error {
	var p rawAWSDirectConnect
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSDirectConnect_Product{}
	terms := []AWSDirectConnect_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSDirectConnect_Term_PriceDimensions{}
				tAttributes := []AWSDirectConnect_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSDirectConnect_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSDirectConnect_Term{
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

type AWSDirectConnect struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AWSDirectConnect_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSDirectConnect_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSDirectConnect_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AWSDirectConnect_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSDirectConnect_Product_Attributes struct {
	gorm.Model
		ToLocationType	string
	VirtualInterfaceType	string
	Servicecode	string
	FromLocation	string
	ToLocation	string
	Usagetype	string
	Operation	string
	Version	string
	TransferType	string
	FromLocationType	string
}

type AWSDirectConnect_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSDirectConnect_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSDirectConnect_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSDirectConnect_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSDirectConnect_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDirectConnect_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSDirectConnect_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSDirectConnect) QueryProducts(q func(product AWSDirectConnect_Product) bool) []AWSDirectConnect_Product{
	ret := []AWSDirectConnect_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSDirectConnect) QueryTerms(t string, q func(product AWSDirectConnect_Term) bool) []AWSDirectConnect_Term{
	ret := []AWSDirectConnect_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AWSDirectConnect) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDirectConnect/current/index.json"
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