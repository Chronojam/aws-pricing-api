package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonChime struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonChime_Product
	Terms		map[string]map[string]map[string]rawAmazonChime_Term
}


type rawAmazonChime_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonChime_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonChime) UnmarshalJSON(data []byte) error {
	var p rawAmazonChime
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonChime_Product{}
	terms := []*AmazonChime_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AmazonChime_Term_PriceDimensions{}
				tAttributes := []*AmazonChime_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonChime_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonChime_Term{
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

type AmazonChime struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AmazonChime_Product `gorm:"ForeignKey:AmazonChimeID"`
	Terms		[]*AmazonChime_Term`gorm:"ForeignKey:AmazonChimeID"`
}
type AmazonChime_Product struct {
	gorm.Model
		AmazonChimeID	uint
	Sku	string
	ProductFamily	string
	Attributes	AmazonChime_Product_Attributes	`gorm:"ForeignKey:AmazonChime_Product_AttributesID"`
}
type AmazonChime_Product_Attributes struct {
	gorm.Model
		AmazonChime_Product_AttributesID	uint
	LocationType	string
	Usagetype	string
	Operation	string
	LicenseType	string
	Servicecode	string
	Location	string
}

type AmazonChime_Term struct {
	gorm.Model
	OfferTermCode string
	AmazonChimeID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AmazonChime_Term_PriceDimensions `gorm:"ForeignKey:AmazonChime_TermID"`
	TermAttributes []*AmazonChime_Term_Attributes `gorm:"ForeignKey:AmazonChime_TermID"`
}

type AmazonChime_Term_Attributes struct {
	gorm.Model
	AmazonChime_TermID	uint
	Key	string
	Value	string
}

type AmazonChime_Term_PriceDimensions struct {
	gorm.Model
	AmazonChime_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AmazonChime_Term_PricePerUnit `gorm:"ForeignKey:AmazonChime_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonChime_Term_PricePerUnit struct {
	gorm.Model
	AmazonChime_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AmazonChime) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonChime/current/index.json"
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