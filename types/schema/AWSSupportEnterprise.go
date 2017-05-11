package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AWSSupportEnterprise struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSSupportEnterprise_Product
	Terms		map[string]map[string]map[string]AWSSupportEnterprise_Term
}
type AWSSupportEnterprise_Product struct {	Attributes	AWSSupportEnterprise_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AWSSupportEnterprise_Product_Attributes struct {	BestPractices	string
	OperationsSupport	string
	TechnicalSupport	string
	WhoCanOpenCases	string
	Usagetype	string
	Operation	string
	ArchitecturalReview	string
	IncludedServices	string
	ThirdpartySoftwareSupport	string
	Servicecode	string
	Location	string
	ProactiveGuidance	string
	ProgrammaticCaseManagement	string
	LocationType	string
	CustomerServiceAndCommunities	string
	CaseSeverityresponseTimes	string
	LaunchSupport	string
	Training	string
	AccountAssistance	string
	ArchitectureSupport	string
}

type AWSSupportEnterprise_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AWSSupportEnterprise_Term_PriceDimensions
	TermAttributes AWSSupportEnterprise_Term_TermAttributes
}

type AWSSupportEnterprise_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AWSSupportEnterprise_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AWSSupportEnterprise_Term_PricePerUnit struct {
	USD	string
}

type AWSSupportEnterprise_Term_TermAttributes struct {

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
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
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