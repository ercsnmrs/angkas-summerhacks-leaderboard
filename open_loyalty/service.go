package open_loyalty

import (
	"context"
	"log/slog"

	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/driver"
)

type OpenLoyaltyService struct {
	Client *OpenLoyaltyClient
	logger *slog.Logger
}

type Label struct {
	Key   string `xml:"key"`
	Value string `xml:"value"`
}
type Address struct {
	Address1 string `xml:"address1,omitempty"`
	Address2 string `xml:"address2,omitempty"`
	City     string `xml:"city,omitempty"`
	Country  string `xml:"country,omitempty"`
	Postal   string `xml:"postal,omitempty"`
	Province string `xml:"province,omitempty"`
	Street   string `xml:"street,omitempty"`
}

type Company struct {
	Name string `xml:"name,omitempty"`
	NIP  string `xml:"nip,omitempty"`
}

type Customer struct {
	Address           Address `xml:"address,omitempty"`
	Agreement1        string  `xml:"agreement1,omitempty"`
	Agreement2        string  `xml:"agreement2,omitempty"`
	Agreement3        string  `xml:"agreement3,omitempty"`
	BirthDate         string  `xml:"birthDate,omitempty"`
	Company           Company `xml:"company,omitempty"`
	Email             string  `xml:"email,omitempty"`
	FirstName         string  `xml:"firstName,omitempty"`
	LastName          string  `xml:"lastName,omitempty"`
	Gender            string  `xml:"gender,omitempty"`
	Labels            []Label `xml:"labels>label,omitempty"`
	LoyaltyCardNumber string  `xml:"loyaltyCardNumber,omitempty"`
	Phone             string  `xml:"phone,omitempty"`
}

type Customers struct {
	Customers []Customer `xml:"customer"`
}

// Creates a new instance of the Open Loyalty Service client with the provided configuration.
func NewProviderSerivce(client OpenLoyaltyClient, logger *slog.Logger) *OpenLoyaltyService {
	svc := &OpenLoyaltyService{
		Client: &client,
		logger: logger,
	}

	return svc
}

func (s *OpenLoyaltyService) ImportDriverRating(ctx context.Context, list []driver.Driver) (err error) {
	// importableList := Customers{}

	// for _, driver := range list {

	// 	customer := Customer{
	// 		LoyaltyCardNumber: driver.ID,
	// 		Labels: []Label{
	// 			{
	// 				Key:   "points",
	// 				Value: strconv.FormatFloat(driver.Rating.Average, 'f', -1, 64),
	// 			},
	// 		},
	// 	}

	// 	importableList.Customers = append(importableList.Customers, customer)
	// }

	// xmlFilePath := time.Now().Format("2006-01-02") + "_temp.xml"

	// // Create XML file
	// err = CreateImportableXMLFile(xmlFilePath, &importableList)
	// if err != nil {
	// 	return fmt.Errorf("failed to create importable XML file: %w", err)
	// }

	// //Call ImportMembers function with XML file path
	// req := importMembersRequest{
	// 	File: xmlFilePath,
	// }

	// resp, err := s.Client.ImportMembers(ctx, req)
	// if err != nil {
	// 	writerErr := DeleteImportableXMLFile(xmlFilePath)
	// 	if writerErr != nil {
	// 		return fmt.Errorf("failed to delete importable XML file: %w", err)
	// 	}

	// 	return err
	// }

	// s.logger.Info(fmt.Sprintf("Success on importing, import ReferenceID: %s", resp.ImportID))

	// // Delete XML file
	// err = DeleteImportableXMLFile(xmlFilePath)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete importable XML file: %w", err)
	// }

	return nil
}
