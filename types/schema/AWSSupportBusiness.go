package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type rawAWSSupportBusiness struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSSupportBusiness_Product
	Terms		map[string]map[string]map[string]rawAWSSupportBusiness_Term
}


type rawAWSSupportBusiness_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSSupportBusiness_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *AWSSupportBusiness) UnmarshalJSON(data []byte) error {
	var p rawAWSSupportBusiness
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSSupportBusiness_Product{}
	terms := []AWSSupportBusiness_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSSupportBusiness_Term_PriceDimensions{}
				tAttributes := []AWSSupportBusiness_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSSupportBusiness_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSSupportBusiness_Term{
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

type AWSSupportBusiness struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AWSSupportBusiness_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSSupportBusiness_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSSupportBusiness_Product struct {
	gorm.Model
		Sku	string
	ProductFamily	string
	Attributes	AWSSupportBusiness_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSSupportBusiness_Product_Attributes struct {
	gorm.Model
		Usagetype	string
	CustomerServiceAndCommunities	string
	ProgrammaticCaseManagement	string
	ThirdpartySoftwareSupport	string
	Location	string
	LocationType	string
	ArchitectureSupport	string
	BestPractices	string
	Servicecode	string
	LaunchSupport	string
	ProactiveGuidance	string
	TechnicalSupport	string
	OperationsSupport	string
	Training	string
	WhoCanOpenCases	string
	Operation	string
	AccountAssistance	string
	ArchitecturalReview	string
	CaseSeverityresponseTimes	string
	IncludedServices	string
}

type AWSSupportBusiness_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSSupportBusiness_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSSupportBusiness_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSSupportBusiness_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSSupportBusiness_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSSupportBusiness_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSSupportBusiness_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSSupportBusiness) QueryProducts(q func(product AWSSupportBusiness_Product) bool) []AWSSupportBusiness_Product{
	ret := []AWSSupportBusiness_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSSupportBusiness) QueryTerms(t string, q func(product AWSSupportBusiness_Term) bool) []AWSSupportBusiness_Term{
	ret := []AWSSupportBusiness_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
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