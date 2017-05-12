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
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSQueueService_Product{}
	terms := []*AWSQueueService_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AWSQueueService_Term_PriceDimensions{}
				tAttributes := []*AWSQueueService_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSQueueService_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSQueueService_Term{
					OfferTermCode: term.OfferTermCode,
					Sku: term.Sku,
					EffectiveDate: term.EffectiveDate,
					TermAttributes: tAttributes,
					PriceDimensions: pDimensions,
				}

				terms = append(terms, &t)
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
	Products	[]*AWSQueueService_Product `gorm:"ForeignKey:AWSQueueServiceID"`
	Terms		[]*AWSQueueService_Term`gorm:"ForeignKey:AWSQueueServiceID"`
}
type AWSQueueService_Product struct {
	gorm.Model
		AWSQueueServiceID	uint
	Sku	string
	ProductFamily	string
	Attributes	AWSQueueService_Product_Attributes	`gorm:"ForeignKey:AWSQueueService_Product_AttributesID"`
}
type AWSQueueService_Product_Attributes struct {
	gorm.Model
		AWSQueueService_Product_AttributesID	uint
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
}

type AWSQueueService_Term struct {
	gorm.Model
	OfferTermCode string
	AWSQueueServiceID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AWSQueueService_Term_PriceDimensions `gorm:"ForeignKey:AWSQueueService_TermID"`
	TermAttributes []*AWSQueueService_Term_Attributes `gorm:"ForeignKey:AWSQueueService_TermID"`
}

type AWSQueueService_Term_Attributes struct {
	gorm.Model
	AWSQueueService_TermID	uint
	Key	string
	Value	string
}

type AWSQueueService_Term_PriceDimensions struct {
	gorm.Model
	AWSQueueService_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AWSQueueService_Term_PricePerUnit `gorm:"ForeignKey:AWSQueueService_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSQueueService_Term_PricePerUnit struct {
	gorm.Model
	AWSQueueService_Term_PriceDimensionsID	uint
	USD	string
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