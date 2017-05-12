package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSStorageGateway struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSStorageGateway_Product
	Terms		map[string]map[string]map[string]rawAWSStorageGateway_Term
}


type rawAWSStorageGateway_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSStorageGateway_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSStorageGateway) UnmarshalJSON(data []byte) error {
	var p rawAWSStorageGateway
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSStorageGateway_Product{}
	terms := []*AWSStorageGateway_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AWSStorageGateway_Term_PriceDimensions{}
				tAttributes := []*AWSStorageGateway_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSStorageGateway_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSStorageGateway_Term{
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

type AWSStorageGateway struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AWSStorageGateway_Product `gorm:"ForeignKey:AWSStorageGatewayID"`
	Terms		[]*AWSStorageGateway_Term`gorm:"ForeignKey:AWSStorageGatewayID"`
}
type AWSStorageGateway_Product struct {
	gorm.Model
		AWSStorageGatewayID	uint
	Sku	string
	ProductFamily	string
	Attributes	AWSStorageGateway_Product_Attributes	`gorm:"ForeignKey:AWSStorageGateway_Product_AttributesID"`
}
type AWSStorageGateway_Product_Attributes struct {
	gorm.Model
		AWSStorageGateway_Product_AttributesID	uint
	Servicecode	string
	TransferType	string
	FromLocation	string
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
}

type AWSStorageGateway_Term struct {
	gorm.Model
	OfferTermCode string
	AWSStorageGatewayID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AWSStorageGateway_Term_PriceDimensions `gorm:"ForeignKey:AWSStorageGateway_TermID"`
	TermAttributes []*AWSStorageGateway_Term_Attributes `gorm:"ForeignKey:AWSStorageGateway_TermID"`
}

type AWSStorageGateway_Term_Attributes struct {
	gorm.Model
	AWSStorageGateway_TermID	uint
	Key	string
	Value	string
}

type AWSStorageGateway_Term_PriceDimensions struct {
	gorm.Model
	AWSStorageGateway_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AWSStorageGateway_Term_PricePerUnit `gorm:"ForeignKey:AWSStorageGateway_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSStorageGateway_Term_PricePerUnit struct {
	gorm.Model
	AWSStorageGateway_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AWSStorageGateway) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSStorageGateway/current/index.json"
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