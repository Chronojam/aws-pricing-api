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
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSDeveloperSupport_Product{}
	terms := []AWSDeveloperSupport_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSDeveloperSupport_Term_PriceDimensions{}
				tAttributes := []AWSDeveloperSupport_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSDeveloperSupport_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSDeveloperSupport_Term{
					OfferTermCode: term.OfferTermCode,
					Sku: term.Sku,
					EffectiveDate: term.EffectiveDate,
					TermAttributes: tAttributes,
					PriceDimensions: pDimensions,
				}

				terms = append(terms, t)
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
	Products	[]AWSDeveloperSupport_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSDeveloperSupport_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSDeveloperSupport_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AWSDeveloperSupport_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSDeveloperSupport_Product_Attributes struct {
	gorm.Model
		LaunchSupport	string
	ProactiveGuidance	string
	TechnicalSupport	string
	BestPractices	string
	IncludedServices	string
	CaseSeverityresponseTimes	string
	CustomerServiceAndCommunities	string
	ProgrammaticCaseManagement	string
	ThirdpartySoftwareSupport	string
	Training	string
	WhoCanOpenCases	string
	Location	string
	Operation	string
	ArchitecturalReview	string
	ArchitectureSupport	string
	OperationsSupport	string
	LocationType	string
	Usagetype	string
	Servicecode	string
	AccountAssistance	string
}

type AWSDeveloperSupport_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSDeveloperSupport_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSDeveloperSupport_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSDeveloperSupport_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSDeveloperSupport_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDeveloperSupport_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSDeveloperSupport_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSDeveloperSupport) QueryProducts(q func(product AWSDeveloperSupport_Product) bool) []AWSDeveloperSupport_Product{
	ret := []AWSDeveloperSupport_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSDeveloperSupport) QueryTerms(t string, q func(product AWSDeveloperSupport_Term) bool) []AWSDeveloperSupport_Term{
	ret := []AWSDeveloperSupport_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
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