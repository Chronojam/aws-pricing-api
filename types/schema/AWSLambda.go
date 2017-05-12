package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSLambda struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSLambda_Product
	Terms		map[string]map[string]map[string]rawAWSLambda_Term
}


type rawAWSLambda_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSLambda_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSLambda) UnmarshalJSON(data []byte) error {
	var p rawAWSLambda
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSLambda_Product{}
	terms := []*AWSLambda_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AWSLambda_Term_PriceDimensions{}
				tAttributes := []*AWSLambda_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSLambda_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSLambda_Term{
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

type AWSLambda struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AWSLambda_Product `gorm:"ForeignKey:AWSLambdaID"`
	Terms		[]*AWSLambda_Term`gorm:"ForeignKey:AWSLambdaID"`
}
type AWSLambda_Product struct {
	gorm.Model
		AWSLambdaID	uint
	Sku	string
	ProductFamily	string
	Attributes	AWSLambda_Product_Attributes	`gorm:"ForeignKey:AWSLambda_Product_AttributesID"`
}
type AWSLambda_Product_Attributes struct {
	gorm.Model
		AWSLambda_Product_AttributesID	uint
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
}

type AWSLambda_Term struct {
	gorm.Model
	OfferTermCode string
	AWSLambdaID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AWSLambda_Term_PriceDimensions `gorm:"ForeignKey:AWSLambda_TermID"`
	TermAttributes []*AWSLambda_Term_Attributes `gorm:"ForeignKey:AWSLambda_TermID"`
}

type AWSLambda_Term_Attributes struct {
	gorm.Model
	AWSLambda_TermID	uint
	Key	string
	Value	string
}

type AWSLambda_Term_PriceDimensions struct {
	gorm.Model
	AWSLambda_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AWSLambda_Term_PricePerUnit `gorm:"ForeignKey:AWSLambda_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSLambda_Term_PricePerUnit struct {
	gorm.Model
	AWSLambda_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AWSLambda) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSLambda/current/index.json"
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