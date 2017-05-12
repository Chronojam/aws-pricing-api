package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSCodePipeline struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSCodePipeline_Product
	Terms		map[string]map[string]map[string]rawAWSCodePipeline_Term
}


type rawAWSCodePipeline_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSCodePipeline_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSCodePipeline) UnmarshalJSON(data []byte) error {
	var p rawAWSCodePipeline
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSCodePipeline_Product{}
	terms := []*AWSCodePipeline_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AWSCodePipeline_Term_PriceDimensions{}
				tAttributes := []*AWSCodePipeline_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSCodePipeline_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSCodePipeline_Term{
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

type AWSCodePipeline struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AWSCodePipeline_Product `gorm:"ForeignKey:AWSCodePipelineID"`
	Terms		[]*AWSCodePipeline_Term`gorm:"ForeignKey:AWSCodePipelineID"`
}
type AWSCodePipeline_Product struct {
	gorm.Model
		AWSCodePipelineID	uint
	Sku	string
	ProductFamily	string
	Attributes	AWSCodePipeline_Product_Attributes	`gorm:"ForeignKey:AWSCodePipeline_Product_AttributesID"`
}
type AWSCodePipeline_Product_Attributes struct {
	gorm.Model
		AWSCodePipeline_Product_AttributesID	uint
	Usagetype	string
	Operation	string
	Servicecode	string
	Description	string
	Location	string
	LocationType	string
}

type AWSCodePipeline_Term struct {
	gorm.Model
	OfferTermCode string
	AWSCodePipelineID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AWSCodePipeline_Term_PriceDimensions `gorm:"ForeignKey:AWSCodePipeline_TermID"`
	TermAttributes []*AWSCodePipeline_Term_Attributes `gorm:"ForeignKey:AWSCodePipeline_TermID"`
}

type AWSCodePipeline_Term_Attributes struct {
	gorm.Model
	AWSCodePipeline_TermID	uint
	Key	string
	Value	string
}

type AWSCodePipeline_Term_PriceDimensions struct {
	gorm.Model
	AWSCodePipeline_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AWSCodePipeline_Term_PricePerUnit `gorm:"ForeignKey:AWSCodePipeline_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSCodePipeline_Term_PricePerUnit struct {
	gorm.Model
	AWSCodePipeline_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AWSCodePipeline) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSCodePipeline/current/index.json"
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