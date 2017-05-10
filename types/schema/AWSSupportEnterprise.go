package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSSupportEnterprise struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSSupportEnterprise_Product
	Terms		map[string]map[string]AWSSupportEnterprise_Term
}
type AWSSupportEnterprise_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSSupportEnterprise_Product_Attributes
}
type AWSSupportEnterprise_Product_Attributes struct {	ArchitectureSupport	string
	LaunchSupport	string
	OperationsSupport	string
	ProgrammaticCaseManagement	string
	TechnicalSupport	string
	WhoCanOpenCases	string
	Location	string
	LocationType	string
	ProactiveGuidance	string
	AccountAssistance	string
	IncludedServices	string
	Operation	string
	ArchitecturalReview	string
	BestPractices	string
	CaseSeverityresponseTimes	string
	CustomerServiceAndCommunities	string
	ThirdpartySoftwareSupport	string
	Servicecode	string
	Usagetype	string
	Training	string
}

type AWSSupportEnterprise_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSSupportEnterprise_Term_PriceDimensions
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