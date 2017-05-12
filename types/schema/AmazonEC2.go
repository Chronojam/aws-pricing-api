package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonEC2 struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonEC2_Product
	Terms		map[string]map[string]map[string]rawAmazonEC2_Term
}


type rawAmazonEC2_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonEC2_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonEC2) UnmarshalJSON(data []byte) error {
	var p rawAmazonEC2
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonEC2_Product{}
	terms := []AmazonEC2_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonEC2_Term_PriceDimensions{}
				tAttributes := []AmazonEC2_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonEC2_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonEC2_Term{
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

type AmazonEC2 struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonEC2_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonEC2_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonEC2_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonEC2_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonEC2_Product_Attributes struct {
	gorm.Model
		LocationType	string
	InstanceFamily	string
	Storage	string
	InstanceType	string
	CurrentGeneration	string
	ClockSpeed	string
	NetworkPerformance	string
	Servicecode	string
	Location	string
	Tenancy	string
	LicenseModel	string
	Usagetype	string
	PreInstalledSw	string
	ProcessorFeatures	string
	Vcpu	string
	PhysicalProcessor	string
	Memory	string
	ProcessorArchitecture	string
	OperatingSystem	string
	Operation	string
}

type AmazonEC2_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonEC2_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonEC2_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonEC2_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonEC2_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonEC2_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonEC2_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonEC2) QueryProducts(q func(product AmazonEC2_Product) bool) []AmazonEC2_Product{
	ret := []AmazonEC2_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonEC2) QueryTerms(t string, q func(product AmazonEC2_Term) bool) []AmazonEC2_Term{
	ret := []AmazonEC2_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonEC2) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEC2/current/index.json"
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