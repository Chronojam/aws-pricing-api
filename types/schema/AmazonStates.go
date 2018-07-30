package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonStates struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonStates_Product
	Terms           map[string]map[string]map[string]rawAmazonStates_Term
}

type rawAmazonStates_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonStates_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonStates) UnmarshalJSON(data []byte) error {
	var p rawAmazonStates
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonStates_Product{}
	terms := []*AmazonStates_Term{}

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
				pDimensions := []*AmazonStates_Term_PriceDimensions{}
				tAttributes := []*AmazonStates_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonStates_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonStates_Term{
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

type AmazonStates struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonStates_Product `gorm:"ForeignKey:AmazonStatesID"`
	Terms           []*AmazonStates_Term    `gorm:"ForeignKey:AmazonStatesID"`
}
type AmazonStates_Product struct {
	gorm.Model
	AmazonStatesID uint
	Sku            string
	ProductFamily  string
	Attributes     AmazonStates_Product_Attributes `gorm:"ForeignKey:AmazonStates_Product_AttributesID"`
}
type AmazonStates_Product_Attributes struct {
	gorm.Model
	AmazonStates_Product_AttributesID uint
	Servicecode                       string
	FromLocation                      string
	ToLocation                        string
	Servicename                       string
	TransferType                      string
	FromLocationType                  string
	ToLocationType                    string
	Usagetype                         string
	Operation                         string
}

type AmazonStates_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonStatesID  uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonStates_Term_PriceDimensions `gorm:"ForeignKey:AmazonStates_TermID"`
	TermAttributes  []*AmazonStates_Term_Attributes      `gorm:"ForeignKey:AmazonStates_TermID"`
}

type AmazonStates_Term_Attributes struct {
	gorm.Model
	AmazonStates_TermID uint
	Key                 string
	Value               string
}

type AmazonStates_Term_PriceDimensions struct {
	gorm.Model
	AmazonStates_TermID uint
	RateCode            string
	RateType            string
	Description         string
	BeginRange          string
	EndRange            string
	Unit                string
	PricePerUnit        *AmazonStates_Term_PricePerUnit `gorm:"ForeignKey:AmazonStates_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonStates_Term_PricePerUnit struct {
	gorm.Model
	AmazonStates_Term_PriceDimensionsID uint
	USD                                 string
}

func (a *AmazonStates) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonStates/current/index.json"
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
