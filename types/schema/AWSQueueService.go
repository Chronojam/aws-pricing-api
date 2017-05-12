package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSQueueService struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSQueueService_Product
	Terms		map[string]map[string]map[string]rawAWSQueueService_Term
}


type rawAWSQueueService_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSQueueService_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSQueueService) UnmarshalJSON(data []byte) error {
	var p rawAWSQueueService
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSQueueService_Product{}
	terms := []AWSQueueService_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSQueueService_Term_PriceDimensions{}
				tAttributes := []AWSQueueService_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSQueueService_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSQueueService_Term{
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

type AWSQueueService struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AWSQueueService_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSQueueService_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSQueueService_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AWSQueueService_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSQueueService_Product_Attributes struct {
	gorm.Model
		ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
}

type AWSQueueService_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSQueueService_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSQueueService_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSQueueService_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSQueueService_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSQueueService_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSQueueService_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSQueueService) QueryProducts(q func(product AWSQueueService_Product) bool) []AWSQueueService_Product{
	ret := []AWSQueueService_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSQueueService) QueryTerms(t string, q func(product AWSQueueService_Term) bool) []AWSQueueService_Term{
	ret := []AWSQueueService_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a *AWSQueueService) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSQueueService/current/index.json"
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