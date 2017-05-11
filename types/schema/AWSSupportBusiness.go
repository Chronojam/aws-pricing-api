package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSSupportBusiness struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSSupportBusiness_Product
	Terms		map[string]map[string]map[string]AWSSupportBusiness_Term
}
type AWSSupportBusiness_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSSupportBusiness_Product_Attributes
}
type AWSSupportBusiness_Product_Attributes struct {	ArchitecturalReview	string
	BestPractices	string
	ProgrammaticCaseManagement	string
	ThirdpartySoftwareSupport	string
	Training	string
	WhoCanOpenCases	string
	Servicecode	string
	AccountAssistance	string
	IncludedServices	string
	LaunchSupport	string
	OperationsSupport	string
	TechnicalSupport	string
	Location	string
	Usagetype	string
	Operation	string
	CaseSeverityresponseTimes	string
	LocationType	string
	ArchitectureSupport	string
	CustomerServiceAndCommunities	string
	ProactiveGuidance	string
}

type AWSSupportBusiness_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSSupportBusiness_Term_PriceDimensions
	TermAttributes AWSSupportBusiness_Term_TermAttributes
}

type AWSSupportBusiness_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSSupportBusiness_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSSupportBusiness_Term_PricePerUnit struct {
	USD	string
}

type AWSSupportBusiness_Term_TermAttributes struct {

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
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
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