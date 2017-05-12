package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSDeveloperSupport struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDeveloperSupport_Product
	Terms		map[string]map[string]map[string]rawAWSDeveloperSupport_Term
}


type rawAWSDeveloperSupport_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSDeveloperSupport_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSDeveloperSupport) UnmarshalJSON(data []byte) error {
	var p rawAWSDeveloperSupport
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSDeveloperSupport_Product{}
	terms := []*AWSDeveloperSupport_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AWSDeveloperSupport_Term_PriceDimensions{}
				tAttributes := []*AWSDeveloperSupport_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSDeveloperSupport_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSDeveloperSupport_Term{
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

type AWSDeveloperSupport struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*AWSDeveloperSupport_Product `gorm:"ForeignKey:AWSDeveloperSupportID"`
	Terms		[]*AWSDeveloperSupport_Term`gorm:"ForeignKey:AWSDeveloperSupportID"`
}
type AWSDeveloperSupport_Product struct {
	gorm.Model
		AWSDeveloperSupportID	uint
	Sku	string
	ProductFamily	string
	Attributes	AWSDeveloperSupport_Product_Attributes	`gorm:"ForeignKey:AWSDeveloperSupport_Product_AttributesID"`
}
type AWSDeveloperSupport_Product_Attributes struct {
	gorm.Model
		AWSDeveloperSupport_Product_AttributesID	uint
	Usagetype	string
	Operation	string
	BestPractices	string
	ProactiveGuidance	string
	ProgrammaticCaseManagement	string
	ArchitecturalReview	string
	IncludedServices	string
	LaunchSupport	string
	Servicecode	string
	Location	string
	AccountAssistance	string
	ArchitectureSupport	string
	CaseSeverityresponseTimes	string
	OperationsSupport	string
	WhoCanOpenCases	string
	LocationType	string
	CustomerServiceAndCommunities	string
	TechnicalSupport	string
	ThirdpartySoftwareSupport	string
	Training	string
}

type AWSDeveloperSupport_Term struct {
	gorm.Model
	OfferTermCode string
	AWSDeveloperSupportID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*AWSDeveloperSupport_Term_PriceDimensions `gorm:"ForeignKey:AWSDeveloperSupport_TermID"`
	TermAttributes []*AWSDeveloperSupport_Term_Attributes `gorm:"ForeignKey:AWSDeveloperSupport_TermID"`
}

type AWSDeveloperSupport_Term_Attributes struct {
	gorm.Model
	AWSDeveloperSupport_TermID	uint
	Key	string
	Value	string
}

type AWSDeveloperSupport_Term_PriceDimensions struct {
	gorm.Model
	AWSDeveloperSupport_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*AWSDeveloperSupport_Term_PricePerUnit `gorm:"ForeignKey:AWSDeveloperSupport_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSDeveloperSupport_Term_PricePerUnit struct {
	gorm.Model
	AWSDeveloperSupport_Term_PriceDimensionsID	uint
	USD	string
}
func (a *AWSDeveloperSupport) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSDeveloperSupport/current/index.json"
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