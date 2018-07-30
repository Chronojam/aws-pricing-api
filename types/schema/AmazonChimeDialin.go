package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonChimeDialin struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonChimeDialin_Product
	Terms           map[string]map[string]map[string]rawAmazonChimeDialin_Term
}

type rawAmazonChimeDialin_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonChimeDialin_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonChimeDialin) UnmarshalJSON(data []byte) error {
	var p rawAmazonChimeDialin
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonChimeDialin_Product{}
	terms := []*AmazonChimeDialin_Term{}

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
				pDimensions := []*AmazonChimeDialin_Term_PriceDimensions{}
				tAttributes := []*AmazonChimeDialin_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonChimeDialin_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonChimeDialin_Term{
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

type AmazonChimeDialin struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonChimeDialin_Product `gorm:"ForeignKey:AmazonChimeDialinID"`
	Terms           []*AmazonChimeDialin_Term    `gorm:"ForeignKey:AmazonChimeDialinID"`
}
type AmazonChimeDialin_Product struct {
	gorm.Model
	AmazonChimeDialinID uint
	Sku                 string
	ProductFamily       string
	Attributes          AmazonChimeDialin_Product_Attributes `gorm:"ForeignKey:AmazonChimeDialin_Product_AttributesID"`
}
type AmazonChimeDialin_Product_Attributes struct {
	gorm.Model
	AmazonChimeDialin_Product_AttributesID uint
	Servicename                            string
	Servicecode                            string
	Location                               string
	LocationType                           string
	Usagetype                              string
	Operation                              string
	CallingType                            string
	Country                                string
}

type AmazonChimeDialin_Term struct {
	gorm.Model
	OfferTermCode       string
	AmazonChimeDialinID uint
	Sku                 string
	EffectiveDate       string
	PriceDimensions     []*AmazonChimeDialin_Term_PriceDimensions `gorm:"ForeignKey:AmazonChimeDialin_TermID"`
	TermAttributes      []*AmazonChimeDialin_Term_Attributes      `gorm:"ForeignKey:AmazonChimeDialin_TermID"`
}

type AmazonChimeDialin_Term_Attributes struct {
	gorm.Model
	AmazonChimeDialin_TermID uint
	Key                      string
	Value                    string
}

type AmazonChimeDialin_Term_PriceDimensions struct {
	gorm.Model
	AmazonChimeDialin_TermID uint
	RateCode                 string
	RateType                 string
	Description              string
	BeginRange               string
	EndRange                 string
	Unit                     string
	PricePerUnit             *AmazonChimeDialin_Term_PricePerUnit `gorm:"ForeignKey:AmazonChimeDialin_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonChimeDialin_Term_PricePerUnit struct {
	gorm.Model
	AmazonChimeDialin_Term_PriceDimensionsID uint
	USD                                      string
}

func (a *AmazonChimeDialin) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonChimeDialin/current/index.json"
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
