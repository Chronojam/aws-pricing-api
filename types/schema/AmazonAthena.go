package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonAthena struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonAthena_Product
	Terms		map[string]map[string]map[string]rawAmazonAthena_Term
}


type rawAmazonAthena_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonAthena_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonAthena) UnmarshalJSON(data []byte) error {
	var p rawAmazonAthena
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonAthena_Product{}
	terms := []*AmazonAthena_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AmazonAthena_Term_PriceDimensions{}
				tAttributes := []*AmazonAthena_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonAthena_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonAthena_Term{
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

type AmazonAthena struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AmazonAthena_Product `gorm:"ForeignKey:AmazonAthenaID"`
	Terms		[]*AmazonAthena_Term`gorm:"ForeignKey:AmazonAthenaID"`
}
type AmazonAthena_Product struct {
	gorm.Model
		AmazonAthenaID	uint
	Sku	string
	ProductFamily	string
	Attributes	AmazonAthena_Product_Attributes	`gorm:"ForeignKey:AmazonAthena_Product_AttributesID"`
}
type AmazonAthena_Product_Attributes struct {
	gorm.Model
		AmazonAthena_Product_AttributesID	uint
	Servicecode	string
	Description	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	FreeQueryTypes	string
}

type AmazonAthena_Term struct {
	gorm.Model
	OfferTermCode string
	AmazonAthenaID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AmazonAthena_Term_PriceDimensions `gorm:"ForeignKey:AmazonAthena_TermID"`
	TermAttributes []*AmazonAthena_Term_Attributes `gorm:"ForeignKey:AmazonAthena_TermID"`
}

type AmazonAthena_Term_Attributes struct {
	gorm.Model
	AmazonAthena_TermID	uint
	Key	string
	Value	string
}

type AmazonAthena_Term_PriceDimensions struct {
	gorm.Model
	AmazonAthena_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AmazonAthena_Term_PricePerUnit `gorm:"ForeignKey:AmazonAthena_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonAthena_Term_PricePerUnit struct {
	gorm.Model
	AmazonAthena_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AmazonAthena) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonAthena/current/index.json"
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