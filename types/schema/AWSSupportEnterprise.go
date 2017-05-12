package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSSupportEnterprise struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSSupportEnterprise_Product
	Terms		map[string]map[string]map[string]rawAWSSupportEnterprise_Term
}


type rawAWSSupportEnterprise_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSSupportEnterprise_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSSupportEnterprise) UnmarshalJSON(data []byte) error {
	var p rawAWSSupportEnterprise
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSSupportEnterprise_Product{}
	terms := []*AWSSupportEnterprise_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AWSSupportEnterprise_Term_PriceDimensions{}
				tAttributes := []*AWSSupportEnterprise_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSSupportEnterprise_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSSupportEnterprise_Term{
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

type AWSSupportEnterprise struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AWSSupportEnterprise_Product `gorm:"ForeignKey:AWSSupportEnterpriseID"`
	Terms		[]*AWSSupportEnterprise_Term`gorm:"ForeignKey:AWSSupportEnterpriseID"`
}
type AWSSupportEnterprise_Product struct {
	gorm.Model
		AWSSupportEnterpriseID	uint
	Sku	string
	ProductFamily	string
	Attributes	AWSSupportEnterprise_Product_Attributes	`gorm:"ForeignKey:AWSSupportEnterprise_Product_AttributesID"`
}
type AWSSupportEnterprise_Product_Attributes struct {
	gorm.Model
		AWSSupportEnterprise_Product_AttributesID	uint
	Servicecode	string
	ArchitecturalReview	string
	BestPractices	string
	ProactiveGuidance	string
	Training	string
	LocationType	string
	Operation	string
	ArchitectureSupport	string
	CustomerServiceAndCommunities	string
	ThirdpartySoftwareSupport	string
	WhoCanOpenCases	string
	Location	string
	Usagetype	string
	CaseSeverityresponseTimes	string
	OperationsSupport	string
	ProgrammaticCaseManagement	string
	AccountAssistance	string
	IncludedServices	string
	LaunchSupport	string
	TechnicalSupport	string
}

type AWSSupportEnterprise_Term struct {
	gorm.Model
	OfferTermCode string
	AWSSupportEnterpriseID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AWSSupportEnterprise_Term_PriceDimensions `gorm:"ForeignKey:AWSSupportEnterprise_TermID"`
	TermAttributes []*AWSSupportEnterprise_Term_Attributes `gorm:"ForeignKey:AWSSupportEnterprise_TermID"`
}

type AWSSupportEnterprise_Term_Attributes struct {
	gorm.Model
	AWSSupportEnterprise_TermID	uint
	Key	string
	Value	string
}

type AWSSupportEnterprise_Term_PriceDimensions struct {
	gorm.Model
	AWSSupportEnterprise_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AWSSupportEnterprise_Term_PricePerUnit `gorm:"ForeignKey:AWSSupportEnterprise_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSSupportEnterprise_Term_PricePerUnit struct {
	gorm.Model
	AWSSupportEnterprise_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AWSSupportEnterprise) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSSupportEnterprise/current/index.json"
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