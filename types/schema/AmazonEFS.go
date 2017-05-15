package schema

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type rawAmazonEFS struct {
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        map[string]AmazonEFS_Product
	Terms           map[string]map[string]map[string]rawAmazonEFS_Term
}

type rawAmazonEFS_Term struct {
	OfferTermCode   string
	Sku             string
	EffectiveDate   string
	PriceDimensions map[string]AmazonEFS_Term_PriceDimensions
	TermAttributes  map[string]string
}

func (l *AmazonEFS) UnmarshalJSON(data []byte) error {
	var p rawAmazonEFS
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*AmazonEFS_Product{}
	terms := []*AmazonEFS_Term{}

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
				pDimensions := []*AmazonEFS_Term_PriceDimensions{}
				tAttributes := []*AmazonEFS_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := AmazonEFS_Term_Attributes{
						Key:   key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := AmazonEFS_Term{
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

type AmazonEFS struct {
	gorm.Model
	FormatVersion   string
	Disclaimer      string
	OfferCode       string
	Version         string
	PublicationDate string
	Products        []*AmazonEFS_Product `gorm:"ForeignKey:AmazonEFSID"`
	Terms           []*AmazonEFS_Term    `gorm:"ForeignKey:AmazonEFSID"`
}
type AmazonEFS_Product struct {
	gorm.Model
	AmazonEFSID   uint
	Sku           string
	ProductFamily string
	Attributes    AmazonEFS_Product_Attributes `gorm:"ForeignKey:AmazonEFS_Product_AttributesID"`
}
type AmazonEFS_Product_Attributes struct {
	gorm.Model
	AmazonEFS_Product_AttributesID uint
	LocationType                   string
	StorageClass                   string
	Usagetype                      string
	Operation                      string
	Servicecode                    string
	Location                       string
}

type AmazonEFS_Term struct {
	gorm.Model
	OfferTermCode   string
	AmazonEFSID     uint
	Sku             string
	EffectiveDate   string
	PriceDimensions []*AmazonEFS_Term_PriceDimensions `gorm:"ForeignKey:AmazonEFS_TermID"`
	TermAttributes  []*AmazonEFS_Term_Attributes      `gorm:"ForeignKey:AmazonEFS_TermID"`
}

type AmazonEFS_Term_Attributes struct {
	gorm.Model
	AmazonEFS_TermID uint
	Key              string
	Value            string
}

type AmazonEFS_Term_PriceDimensions struct {
	gorm.Model
	AmazonEFS_TermID uint
	RateCode         string
	RateType         string
	Description      string
	BeginRange       string
	EndRange         string
	Unit             string
	PricePerUnit     *AmazonEFS_Term_PricePerUnit `gorm:"ForeignKey:AmazonEFS_Term_PriceDimensionsID"`
	// AppliesTo	[]string
}

type AmazonEFS_Term_PricePerUnit struct {
	gorm.Model
	AmazonEFS_Term_PriceDimensionsID uint
	USD                              string
}

func (a *AmazonEFS) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonEFS/current/index.json"
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
