package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawOpsWorks struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]OpsWorks_Product
	Terms		map[string]map[string]map[string]rawOpsWorks_Term
}


type rawOpsWorks_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]OpsWorks_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *OpsWorks) UnmarshalJSON(data []byte) error {
	var p rawOpsWorks
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*OpsWorks_Product{}
	terms := []*OpsWorks_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*OpsWorks_Term_PriceDimensions{}
				tAttributes := []*OpsWorks_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := OpsWorks_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := OpsWorks_Term{
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

type OpsWorks struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*OpsWorks_Product `gorm:"ForeignKey:OpsWorksID"`
	Terms		[]*OpsWorks_Term`gorm:"ForeignKey:OpsWorksID"`
}
type OpsWorks_Product struct {
	gorm.Model
		OpsWorksID	uint
	ProductFamily	string
	Attributes	OpsWorks_Product_Attributes	`gorm:"ForeignKey:OpsWorks_Product_AttributesID"`
	Sku	string
}
type OpsWorks_Product_Attributes struct {
	gorm.Model
		OpsWorks_Product_AttributesID	uint
	LocationType	string
	Group	string
	Usagetype	string
	Operation	string
	ServerLocation	string
	Servicecode	string
	Location	string
}

type OpsWorks_Term struct {
	gorm.Model
	OfferTermCode string
	OpsWorksID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*OpsWorks_Term_PriceDimensions `gorm:"ForeignKey:OpsWorks_TermID"`
	TermAttributes []*OpsWorks_Term_Attributes `gorm:"ForeignKey:OpsWorks_TermID"`
}

type OpsWorks_Term_Attributes struct {
	gorm.Model
	OpsWorks_TermID	uint
	Key	string
	Value	string
}

type OpsWorks_Term_PriceDimensions struct {
	gorm.Model
	OpsWorks_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*OpsWorks_Term_PricePerUnit `gorm:"ForeignKey:OpsWorks_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type OpsWorks_Term_PricePerUnit struct {
	gorm.Model
	OpsWorks_Term_PriceDimensionsID	uint
	USD	string
}
func (a *OpsWorks) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/OpsWorks/current/index.json"
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