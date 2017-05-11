package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSDeveloperSupport struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDeveloperSupport_Product
	Terms		map[string]map[string]map[string]AWSDeveloperSupport_Term
}
type AWSDeveloperSupport_Product struct {	ProductFamily	string
	Attributes	AWSDeveloperSupport_Product_Attributes
	Sku	string
}
type AWSDeveloperSupport_Product_Attributes struct {	LocationType	string
	BestPractices	string
	ProactiveGuidance	string
	ProgrammaticCaseManagement	string
	ThirdpartySoftwareSupport	string
	Usagetype	string
	Operation	string
	LaunchSupport	string
	OperationsSupport	string
	CaseSeverityresponseTimes	string
	CustomerServiceAndCommunities	string
	IncludedServices	string
	TechnicalSupport	string
	Location	string
	AccountAssistance	string
	ArchitecturalReview	string
	ArchitectureSupport	string
	WhoCanOpenCases	string
	Servicecode	string
	Training	string
}

type AWSDeveloperSupport_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSDeveloperSupport_Term_PriceDimensions
	TermAttributes AWSDeveloperSupport_Term_TermAttributes
}

type AWSDeveloperSupport_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSDeveloperSupport_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSDeveloperSupport_Term_PricePerUnit struct {
	USD	string
}

type AWSDeveloperSupport_Term_TermAttributes struct {

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
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
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