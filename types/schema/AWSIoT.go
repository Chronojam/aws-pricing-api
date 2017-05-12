package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSIoT struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSIoT_Product
	Terms		map[string]map[string]map[string]rawAWSIoT_Term
}


type rawAWSIoT_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSIoT_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSIoT) UnmarshalJSON(data []byte) error {
	var p rawAWSIoT
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSIoT_Product{}
	terms := []*AWSIoT_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AWSIoT_Term_PriceDimensions{}
				tAttributes := []*AWSIoT_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSIoT_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSIoT_Term{
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

type AWSIoT struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AWSIoT_Product `gorm:"ForeignKey:AWSIoTID"`
	Terms		[]*AWSIoT_Term`gorm:"ForeignKey:AWSIoTID"`
}
type AWSIoT_Product struct {
	gorm.Model
		AWSIoTID	uint
	Sku	string
	ProductFamily	string
	Attributes	AWSIoT_Product_Attributes	`gorm:"ForeignKey:AWSIoT_Product_AttributesID"`
}
type AWSIoT_Product_Attributes struct {
	gorm.Model
		AWSIoT_Product_AttributesID	uint
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	Isshadow	string
	Iswebsocket	string
	Protocol	string
	Servicecode	string
}

type AWSIoT_Term struct {
	gorm.Model
	OfferTermCode string
	AWSIoTID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AWSIoT_Term_PriceDimensions `gorm:"ForeignKey:AWSIoT_TermID"`
	TermAttributes []*AWSIoT_Term_Attributes `gorm:"ForeignKey:AWSIoT_TermID"`
}

type AWSIoT_Term_Attributes struct {
	gorm.Model
	AWSIoT_TermID	uint
	Key	string
	Value	string
}

type AWSIoT_Term_PriceDimensions struct {
	gorm.Model
	AWSIoT_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AWSIoT_Term_PricePerUnit `gorm:"ForeignKey:AWSIoT_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSIoT_Term_PricePerUnit struct {
	gorm.Model
	AWSIoT_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AWSIoT) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSIoT/current/index.json"
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