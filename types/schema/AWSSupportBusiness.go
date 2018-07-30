package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAWSSupportBusiness struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AWSSupportBusiness_Product
	Terms           map[string]map[string]map[string]rawAWSSupportBusiness_Term
}

type rawAWSSupportBusiness_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AWSSupportBusiness_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AWSSupportBusiness) UnmarshalJSON(data []byte) error {
	var p rawAWSSupportBusiness
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AWSSupportBusiness_Product{}
	terms := []*AWSSupportBusiness_Term{}

	// Convert from map to slice
	for i, _ := range p.Products {
		pr := p.Products[i]
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*AWSSupportBusiness_Term_PriceDimensions{}
				tAttributes := []*AWSSupportBusiness_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSSupportBusiness_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AWSSupportBusiness_Term{
					OfferTermCode:   term.OfferTermCode,
					Sku:             term.Sku,
					EffectiveDate:   term.EffectiveDate,
					TermAttributes:  tAttributes,
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

type AWSSupportBusiness struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AWSSupportBusiness_Product `gorm:"ForeignKey:AWSSupportBusinessID"`
	Terms           []*AWSSupportBusiness_Term    `gorm:"ForeignKey:AWSSupportBusinessID"`
}
type AWSSupportBusiness_Product struct {
	gorm.Model
	AWSSupportBusinessID uint
	Sku                  string
	ProductFamily        string
	Attributes           AWSSupportBusiness_Product_Attributes `gorm:"ForeignKey:AWSSupportBusiness_Product_AttributesID"`
}
type AWSSupportBusiness_Product_Attributes struct {
	gorm.Model
	AWSSupportBusiness_Product_AttributesID uint
	Usagetype                               string
	LaunchSupport                           string
	WhoCanOpenCases                         string
	Servicecode                             string
	LocationType                            string
	BestPractices                           string
	ProactiveGuidance                       string
	Training                                string
	Location                                string
	CaseSeverityresponseTimes               string
	OperationsSupport                       string
	ProgrammaticCaseManagement              string
	ArchitectureSupport                     string
	AccountAssistance                       string
	ArchitecturalReview                     string
	CustomerServiceAndCommunities           string
	IncludedServices                        string
	TechnicalSupport                        string
	ThirdpartySoftwareSupport               string
	Operation                               string
}

type AWSSupportBusiness_Term struct {
	gorm.Model
	OfferTermCode        string
	AWSSupportBusinessID uint
	Sku                  string
	EffectiveDate        string
	PriceDimensions      []*AWSSupportBusiness_Term_PriceDimensions `gorm:"ForeignKey:AWSSupportBusiness_TermID"`
	TermAttributes       []*AWSSupportBusiness_Term_Attributes      `gorm:"ForeignKey:AWSSupportBusiness_TermID"`
}

type AWSSupportBusiness_Term_Attributes struct {
	gorm.Model
	AWSSupportBusiness_TermID uint
	Key                       string
	Value                     string
}

type AWSSupportBusiness_Term_PriceDimensions struct {
	gorm.Model
	AWSSupportBusiness_TermID uint
	RateCode                  string
	RateType                  string
	Description               string
	BeginRange                string
	EndRange                  string
	Unit                      string
	PricePerUnit              *AWSSupportBusiness_Term_PricePerUnit `gorm:"ForeignKey:AWSSupportBusiness_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AWSSupportBusiness_Term_PricePerUnit struct {
	gorm.Model
	AWSSupportBusiness_Term_PriceDimensionsID uint
	USD                                       string
}

func (a *AWSSupportBusiness) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AWSSupportBusiness/current/index.json"
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
