package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonRDS struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonRDS_Product
	Terms		map[string]map[string]map[string]rawAmazonRDS_Term
}


type rawAmazonRDS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonRDS_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonRDS) UnmarshalJSON(data []byte) error {
	var p rawAmazonRDS
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AmazonRDS_Product{}
	terms := []AmazonRDS_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AmazonRDS_Term_PriceDimensions{}
				tAttributes := []AmazonRDS_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonRDS_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AmazonRDS_Term{
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

type AmazonRDS struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AmazonRDS_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AmazonRDS_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonRDS_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AmazonRDS_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AmazonRDS_Product_Attributes struct {
	gorm.Model
		Vcpu	string
	PhysicalProcessor	string
	NetworkPerformance	string
	EngineCode	string
	LicenseModel	string
	InstanceType	string
	CurrentGeneration	string
	InstanceFamily	string
	ProcessorFeatures	string
	Storage	string
	DatabaseEngine	string
	DatabaseEdition	string
	Usagetype	string
	EnhancedNetworkingSupported	string
	Servicecode	string
	Location	string
	Memory	string
	LocationType	string
	ClockSpeed	string
	ProcessorArchitecture	string
	DeploymentOption	string
	Operation	string
}

type AmazonRDS_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AmazonRDS_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AmazonRDS_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AmazonRDS_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AmazonRDS_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonRDS_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AmazonRDS_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AmazonRDS) QueryProducts(q func(product AmazonRDS_Product) bool) []AmazonRDS_Product{
	ret := []AmazonRDS_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonRDS) QueryTerms(t string, q func(product AmazonRDS_Term) bool) []AmazonRDS_Term{
	ret := []AmazonRDS_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AmazonRDS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRDS/current/index.json"
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