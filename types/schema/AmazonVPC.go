package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonVPC struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonVPC_Product
	Terms		map[string]map[string]map[string]rawAmazonVPC_Term
}


type rawAmazonVPC_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonVPC_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonVPC) UnmarshalJSON(data []byte) error {
	var p rawAmazonVPC
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonVPC_Product{}
	terms := []*AmazonVPC_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AmazonVPC_Term_PriceDimensions{}
				tAttributes := []*AmazonVPC_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonVPC_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonVPC_Term{
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

type AmazonVPC struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AmazonVPC_Product `gorm:"ForeignKey:AmazonVPCID"`
	Terms		[]*AmazonVPC_Term`gorm:"ForeignKey:AmazonVPCID"`
}
type AmazonVPC_Product struct {
	gorm.Model
		AmazonVPCID	uint
	Sku	string
	ProductFamily	string
	Attributes	AmazonVPC_Product_Attributes	`gorm:"ForeignKey:AmazonVPC_Product_AttributesID"`
}
type AmazonVPC_Product_Attributes struct {
	gorm.Model
		AmazonVPC_Product_AttributesID	uint
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
}

type AmazonVPC_Term struct {
	gorm.Model
	OfferTermCode string
	AmazonVPCID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AmazonVPC_Term_PriceDimensions `gorm:"ForeignKey:AmazonVPC_TermID"`
	TermAttributes []*AmazonVPC_Term_Attributes `gorm:"ForeignKey:AmazonVPC_TermID"`
}

type AmazonVPC_Term_Attributes struct {
	gorm.Model
	AmazonVPC_TermID	uint
	Key	string
	Value	string
}

type AmazonVPC_Term_PriceDimensions struct {
	gorm.Model
	AmazonVPC_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AmazonVPC_Term_PricePerUnit `gorm:"ForeignKey:AmazonVPC_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonVPC_Term_PricePerUnit struct {
	gorm.Model
	AmazonVPC_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AmazonVPC) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonVPC/current/index.json"
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