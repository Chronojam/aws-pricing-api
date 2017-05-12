package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonCloudSearch struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonCloudSearch_Product
	Terms		map[string]map[string]map[string]rawAmazonCloudSearch_Term
}


type rawAmazonCloudSearch_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonCloudSearch_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonCloudSearch) UnmarshalJSON(data []byte) error {
	var p rawAmazonCloudSearch
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonCloudSearch_Product{}
	terms := []*AmazonCloudSearch_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AmazonCloudSearch_Term_PriceDimensions{}
				tAttributes := []*AmazonCloudSearch_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonCloudSearch_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonCloudSearch_Term{
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

type AmazonCloudSearch struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AmazonCloudSearch_Product `gorm:"ForeignKey:AmazonCloudSearchID"`
	Terms		[]*AmazonCloudSearch_Term`gorm:"ForeignKey:AmazonCloudSearchID"`
}
type AmazonCloudSearch_Product struct {
	gorm.Model
		AmazonCloudSearchID	uint
	Attributes	AmazonCloudSearch_Product_Attributes	`gorm:"ForeignKey:AmazonCloudSearch_Product_AttributesID"`
	Sku	string
	ProductFamily	string
}
type AmazonCloudSearch_Product_Attributes struct {
	gorm.Model
		AmazonCloudSearch_Product_AttributesID	uint
	FromLocationType	string
	ToLocation	string
	ToLocationType	string
	Usagetype	string
	Operation	string
	Servicecode	string
	TransferType	string
	FromLocation	string
}

type AmazonCloudSearch_Term struct {
	gorm.Model
	OfferTermCode string
	AmazonCloudSearchID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AmazonCloudSearch_Term_PriceDimensions `gorm:"ForeignKey:AmazonCloudSearch_TermID"`
	TermAttributes []*AmazonCloudSearch_Term_Attributes `gorm:"ForeignKey:AmazonCloudSearch_TermID"`
}

type AmazonCloudSearch_Term_Attributes struct {
	gorm.Model
	AmazonCloudSearch_TermID	uint
	Key	string
	Value	string
}

type AmazonCloudSearch_Term_PriceDimensions struct {
	gorm.Model
	AmazonCloudSearch_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AmazonCloudSearch_Term_PricePerUnit `gorm:"ForeignKey:AmazonCloudSearch_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonCloudSearch_Term_PricePerUnit struct {
	gorm.Model
	AmazonCloudSearch_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AmazonCloudSearch) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonCloudSearch/current/index.json"
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