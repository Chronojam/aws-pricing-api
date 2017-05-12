package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAmazonRDS struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonRDS_Product
	Terms		map[string]map[string]map[string]rawAmazonRDS_Term
}


type rawAmazonRDS_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonRDS_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AmazonRDS) UnmarshalJSON(data []byte) error {
	var p rawAmazonRDS
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonRDS_Product{}
	terms := []*AmazonRDS_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AmazonRDS_Term_PriceDimensions{}
				tAttributes := []*AmazonRDS_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonRDS_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonRDS_Term{
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

type AmazonRDS struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AmazonRDS_Product `gorm:"ForeignKey:AmazonRDSID"`
	Terms		[]*AmazonRDS_Term`gorm:"ForeignKey:AmazonRDSID"`
}
type AmazonRDS_Product struct {
	gorm.Model
		AmazonRDSID	uint
	Sku	string
	ProductFamily	string
	Attributes	AmazonRDS_Product_Attributes	`gorm:"ForeignKey:AmazonRDS_Product_AttributesID"`
}
type AmazonRDS_Product_Attributes struct {
	gorm.Model
		AmazonRDS_Product_AttributesID	uint
	Servicecode	string
	CurrentGeneration	string
	InstanceFamily	string
	Vcpu	string
	Storage	string
	EngineCode	string
	DatabaseEdition	string
	DeploymentOption	string
	Usagetype	string
	LocationType	string
	Memory	string
	NetworkPerformance	string
	ProcessorArchitecture	string
	Operation	string
	InstanceType	string
	DatabaseEngine	string
	Location	string
	PhysicalProcessor	string
	LicenseModel	string
}

type AmazonRDS_Term struct {
	gorm.Model
	OfferTermCode string
	AmazonRDSID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AmazonRDS_Term_PriceDimensions `gorm:"ForeignKey:AmazonRDS_TermID"`
	TermAttributes []*AmazonRDS_Term_Attributes `gorm:"ForeignKey:AmazonRDS_TermID"`
}

type AmazonRDS_Term_Attributes struct {
	gorm.Model
	AmazonRDS_TermID	uint
	Key	string
	Value	string
}

type AmazonRDS_Term_PriceDimensions struct {
	gorm.Model
	AmazonRDS_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AmazonRDS_Term_PricePerUnit `gorm:"ForeignKey:AmazonRDS_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonRDS_Term_PricePerUnit struct {
	gorm.Model
	AmazonRDS_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AmazonRDS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonRDS/current/index.json"
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