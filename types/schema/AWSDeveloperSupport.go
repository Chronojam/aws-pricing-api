package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AWSDeveloperSupport struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AWSDeveloperSupport_Product
	Terms		map[string]map[string]AWSDeveloperSupport_Term
}
type AWSDeveloperSupport_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AWSDeveloperSupport_Product_Attributes
}
type AWSDeveloperSupport_Product_Attributes struct {	ThirdpartySoftwareSupport	string
	CustomerServiceAndCommunities	string
	IncludedServices	string
	OperationsSupport	string
	ProactiveGuidance	string
	Operation	string
	ArchitecturalReview	string
	ArchitectureSupport	string
	CaseSeverityresponseTimes	string
	ProgrammaticCaseManagement	string
	Usagetype	string
	AccountAssistance	string
	TechnicalSupport	string
	Training	string
	LaunchSupport	string
	WhoCanOpenCases	string
	Servicecode	string
	Location	string
	LocationType	string
	BestPractices	string
}

type AWSDeveloperSupport_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions AWSDeveloperSupport_Term_PriceDimensions
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