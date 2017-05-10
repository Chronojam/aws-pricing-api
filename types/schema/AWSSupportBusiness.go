package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSSupportBusiness struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSSupportBusiness_Product
	Terms		map[string]map[string]AWSSupportBusiness_Term
}
type AWSSupportBusiness_Product struct {	Attributes	AWSSupportBusiness_Product_Attributes
	Sku	string
	ProductFamily	string
}
type AWSSupportBusiness_Product_Attributes struct {	ProactiveGuidance	string
	CaseSeverityresponseTimes	string
	LocationType	string
	AccountAssistance	string
	ArchitecturalReview	string
	BestPractices	string
	IncludedServices	string
	LaunchSupport	string
	Training	string
	Servicecode	string
	WhoCanOpenCases	string
	Operation	string
	CustomerServiceAndCommunities	string
	OperationsSupport	string
	TechnicalSupport	string
	Usagetype	string
	ArchitectureSupport	string
	ProgrammaticCaseManagement	string
	ThirdpartySoftwareSupport	string
	Location	string
}

type AWSSupportBusiness_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSSupportBusiness_Term_PriceDimensions
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