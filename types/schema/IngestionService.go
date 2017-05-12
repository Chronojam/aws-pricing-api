package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawIngestionService struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]IngestionService_Product
	Terms		map[string]map[string]map[string]rawIngestionService_Term
}


type rawIngestionService_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]IngestionService_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *IngestionService) UnmarshalJSON(data []byte) error {
	var p rawIngestionService
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*IngestionService_Product{}
	terms := []*IngestionService_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*IngestionService_Term_PriceDimensions{}
				tAttributes := []*IngestionService_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := IngestionService_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := IngestionService_Term{
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

type IngestionService struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*IngestionService_Product `gorm:"ForeignKey:IngestionServiceID"`
	Terms		[]*IngestionService_Term`gorm:"ForeignKey:IngestionServiceID"`
}
type IngestionService_Product struct {
	gorm.Model
		IngestionServiceID	uint
	Sku	string
	ProductFamily	string
	Attributes	IngestionService_Product_Attributes	`gorm:"ForeignKey:IngestionService_Product_AttributesID"`
}
type IngestionService_Product_Attributes struct {
	gorm.Model
		IngestionService_Product_AttributesID	uint
	DataAction	string
	Servicecode	string
	Location	string
	LocationType	string
	Group	string
	GroupDescription	string
	Usagetype	string
	Operation	string
}

type IngestionService_Term struct {
	gorm.Model
	OfferTermCode string
	IngestionServiceID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*IngestionService_Term_PriceDimensions `gorm:"ForeignKey:IngestionService_TermID"`
	TermAttributes []*IngestionService_Term_Attributes `gorm:"ForeignKey:IngestionService_TermID"`
}

type IngestionService_Term_Attributes struct {
	gorm.Model
	IngestionService_TermID	uint
	Key	string
	Value	string
}

type IngestionService_Term_PriceDimensions struct {
	gorm.Model
	IngestionService_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*IngestionService_Term_PricePerUnit `gorm:"ForeignKey:IngestionService_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type IngestionService_Term_PricePerUnit struct {
	gorm.Model
	IngestionService_Term_PriceDimensionsID	uint
	USD	string
}
func (a *IngestionService) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/IngestionService/current/index.json"
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