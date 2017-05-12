package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonRoute53 struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonRoute53_Product
	Terms		map[string]map[string]map[string]rawAmazonRoute53_Term
}


type rawAmazonRoute53_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonRoute53_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonRoute53) UnmarshalJSON(data []byte) error {
	var p rawAmazonRoute53
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonRoute53_Product{}
	terms := []*AmazonRoute53_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AmazonRoute53_Term_PriceDimensions{}
				tAttributes := []*AmazonRoute53_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonRoute53_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonRoute53_Term{
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

type AmazonRoute53 struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AmazonRoute53_Product `gorm:"ForeignKey:AmazonRoute53ID"`
	Terms		[]*AmazonRoute53_Term`gorm:"ForeignKey:AmazonRoute53ID"`
}
type AmazonRoute53_Product struct {
	gorm.Model
		AmazonRoute53ID	uint
	Sku	string
	ProductFamily	string
	Attributes	AmazonRoute53_Product_Attributes	`gorm:"ForeignKey:AmazonRoute53_Product_AttributesID"`
}
type AmazonRoute53_Product_Attributes struct {
	gorm.Model
		AmazonRoute53_Product_AttributesID	uint
	Servicecode	string
	RoutingType	string
	RoutingTarget	string
	Usagetype	string
	Operation	string
}

type AmazonRoute53_Term struct {
	gorm.Model
	OfferTermCode string
	AmazonRoute53ID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AmazonRoute53_Term_PriceDimensions `gorm:"ForeignKey:AmazonRoute53_TermID"`
	TermAttributes []*AmazonRoute53_Term_Attributes `gorm:"ForeignKey:AmazonRoute53_TermID"`
}

type AmazonRoute53_Term_Attributes struct {
	gorm.Model
	AmazonRoute53_TermID	uint
	Key	string
	Value	string
}

type AmazonRoute53_Term_PriceDimensions struct {
	gorm.Model
	AmazonRoute53_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AmazonRoute53_Term_PricePerUnit `gorm:"ForeignKey:AmazonRoute53_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonRoute53_Term_PricePerUnit struct {
	gorm.Model
	AmazonRoute53_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AmazonRoute53) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRoute53/current/index.json"
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