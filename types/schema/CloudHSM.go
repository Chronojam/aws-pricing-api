package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawCloudHSM struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]CloudHSM_Product
	Terms           map[string]map[string]map[string]rawCloudHSM_Term
}

type rawCloudHSM_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]CloudHSM_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *CloudHSM) UnmarshalJSON(data []byte) error {
	var p rawCloudHSM
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*CloudHSM_Product{}
	terms := []*CloudHSM_Term{}

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
				pDimensions := []*CloudHSM_Term_PriceDimensions{}
				tAttributes := []*CloudHSM_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := CloudHSM_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := CloudHSM_Term{
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

type CloudHSM struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*CloudHSM_Product `gorm:"ForeignKey:CloudHSMID"`
	Terms           []*CloudHSM_Term    `gorm:"ForeignKey:CloudHSMID"`
}
type CloudHSM_Product struct {
	gorm.Model
	CloudHSMID    uint
	Sku           string
	ProductFamily string
	Attributes    CloudHSM_Product_Attributes `gorm:"ForeignKey:CloudHSM_Product_AttributesID"`
}
type CloudHSM_Product_Attributes struct {
	gorm.Model
	CloudHSM_Product_AttributesID uint
	UpfrontCommitment             string
	Servicecode                   string
	Location                      string
	LocationType                  string
	InstanceFamily                string
	Usagetype                     string
	Operation                     string
	TrialProduct                  string
}

type CloudHSM_Term struct {
	gorm.Model
	OfferTermCode   string
	CloudHSMID      uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*CloudHSM_Term_PriceDimensions `gorm:"ForeignKey:CloudHSM_TermID"`
	TermAttributes  []*CloudHSM_Term_Attributes      `gorm:"ForeignKey:CloudHSM_TermID"`
}

type CloudHSM_Term_Attributes struct {
	gorm.Model
	CloudHSM_TermID uint
	Key             string
	Value           string
}

type CloudHSM_Term_PriceDimensions struct {
	gorm.Model
	CloudHSM_TermID uint
	RateCode        string
	RateType        string
	Description     string
	BeginRange      string
	EndRange        string
	Unit            string
	PricePerUnit    *CloudHSM_Term_PricePerUnit `gorm:"ForeignKey:CloudHSM_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type CloudHSM_Term_PricePerUnit struct {
	gorm.Model
	CloudHSM_Term_PriceDimensionsID uint
	USD                             string
}

func (a *CloudHSM) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/CloudHSM/current/index.json"
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
