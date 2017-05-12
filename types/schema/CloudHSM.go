package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawCloudHSM struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]CloudHSM_Product
	Terms		map[string]map[string]map[string]rawCloudHSM_Term
}


type rawCloudHSM_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]CloudHSM_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *CloudHSM) UnmarshalJSON(data []byte) error {
	var p rawCloudHSM
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []CloudHSM_Product{}
	terms := []CloudHSM_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []CloudHSM_Term_PriceDimensions{}
				tAttributes := []CloudHSM_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := CloudHSM_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := CloudHSM_Term{
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

type CloudHSM struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]CloudHSM_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]CloudHSM_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type CloudHSM_Product struct {
	gorm.Model
		Attributes	CloudHSM_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
	ProductFamily	string
}
type CloudHSM_Product_Attributes struct {
	gorm.Model
		InstanceFamily	string
	Usagetype	string
	Operation	string
	TrialProduct	string
	UpfrontCommitment	string
	Servicecode	string
	Location	string
	LocationType	string
}

type CloudHSM_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []CloudHSM_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []CloudHSM_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type CloudHSM_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type CloudHSM_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	CloudHSM_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type CloudHSM_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a CloudHSM) QueryProducts(q func(product CloudHSM_Product) bool) []CloudHSM_Product{
	ret := []CloudHSM_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a CloudHSM) QueryTerms(t string, q func(product CloudHSM_Term) bool) []CloudHSM_Term{
	ret := []CloudHSM_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *CloudHSM) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/CloudHSM/current/index.json"
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