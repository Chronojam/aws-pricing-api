package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonETS struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonETS_Product
	Terms		map[string]map[string]map[string]rawAmazonETS_Term
}


type rawAmazonETS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonETS_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonETS) UnmarshalJSON(data []byte) error {
	var p rawAmazonETS
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonETS_Product{}
	terms := []*AmazonETS_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AmazonETS_Term_PriceDimensions{}
				tAttributes := []*AmazonETS_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonETS_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonETS_Term{
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

type AmazonETS struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AmazonETS_Product `gorm:"ForeignKey:AmazonETSID"`
	Terms		[]*AmazonETS_Term`gorm:"ForeignKey:AmazonETSID"`
}
type AmazonETS_Product struct {
	gorm.Model
		AmazonETSID	uint
	Attributes	AmazonETS_Product_Attributes	`gorm:"ForeignKey:AmazonETS_Product_AttributesID"`
	Sku	string
	ProductFamily	string
}
type AmazonETS_Product_Attributes struct {
	gorm.Model
		AmazonETS_Product_AttributesID	uint
	Operation	string
	TranscodingResult	string
	VideoResolution	string
	Servicecode	string
	Location	string
	LocationType	string
	Usagetype	string
}

type AmazonETS_Term struct {
	gorm.Model
	OfferTermCode string
	AmazonETSID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AmazonETS_Term_PriceDimensions `gorm:"ForeignKey:AmazonETS_TermID"`
	TermAttributes []*AmazonETS_Term_Attributes `gorm:"ForeignKey:AmazonETS_TermID"`
}

type AmazonETS_Term_Attributes struct {
	gorm.Model
	AmazonETS_TermID	uint
	Key	string
	Value	string
}

type AmazonETS_Term_PriceDimensions struct {
	gorm.Model
	AmazonETS_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AmazonETS_Term_PricePerUnit `gorm:"ForeignKey:AmazonETS_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonETS_Term_PricePerUnit struct {
	gorm.Model
	AmazonETS_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AmazonETS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonETS/current/index.json"
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