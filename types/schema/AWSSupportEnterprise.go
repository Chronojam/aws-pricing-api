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
	err := json.Unmarshal(data, p)
	if err != nil {
		return err
	}

	products := []AWSSupportEnterprise_Product{}
	terms := []AWSSupportEnterprise_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []AWSSupportEnterprise_Term_PriceDimensions{}
				tAttributes := []AWSSupportEnterprise_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, pd)
				}

				for key, value := range term.TermAttributes {
					tr := AWSSupportEnterprise_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, tr)
				}

				t := AWSSupportEnterprise_Term{
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

type AWSSupportEnterprise struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]AWSSupportEnterprise_Product 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Terms		[]AWSSupportEnterprise_Term	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}
type AWSSupportEnterprise_Product struct {
	gorm.Model
		Attributes	AWSSupportEnterprise_Product_Attributes	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	Sku	string
	ProductFamily	string
}
type AWSSupportEnterprise_Product_Attributes struct {
	gorm.Model
		Servicecode	string
	OperationsSupport	string
	CustomerServiceAndCommunities	string
	IncludedServices	string
	ProgrammaticCaseManagement	string
	ThirdpartySoftwareSupport	string
	ArchitectureSupport	string
	LaunchSupport	string
	Training	string
	WhoCanOpenCases	string
	BestPractices	string
	CaseSeverityresponseTimes	string
	Location	string
	LocationType	string
	Usagetype	string
	Operation	string
	AccountAssistance	string
	ArchitecturalReview	string
	ProactiveGuidance	string
	TechnicalSupport	string
}

type AWSSupportEnterprise_Term struct {
	gorm.Model
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions []AWSSupportEnterprise_Term_PriceDimensions 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	TermAttributes []AWSSupportEnterprise_Term_Attributes 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
}

type AWSSupportEnterprise_Term_Attributes struct {
	gorm.Model
	Key	string
	Value	string
}

type AWSSupportEnterprise_Term_PriceDimensions struct {
	gorm.Model
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSSupportEnterprise_Term_PricePerUnit 	`gorm:"ForeignKey:ID,type:varchar(255)[]"`
	AppliesTo	[]interface{}
}

type AWSSupportEnterprise_Term_PricePerUnit struct {
	gorm.Model
	USD	string
}
func (a AWSSupportEnterprise) QueryProducts(q func(product AWSSupportEnterprise_Product) bool) []AWSSupportEnterprise_Product{
	ret := []AWSSupportEnterprise_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AWSSupportEnterprise) QueryTerms(t string, q func(product AWSSupportEnterprise_Term) bool) []AWSSupportEnterprise_Term{
	ret := []AWSSupportEnterprise_Term{}
	for _, v := range a.Terms {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
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