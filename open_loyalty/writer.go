package open_loyalty

import (
	"encoding/xml"
	"fmt"
	"os"
)

func CreateImportableXMLFile(filename string, list *Customers) error {
	xmlData, err := xml.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	// Write the XML declaration
	fmt.Fprintf(file, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")

	// Write the XML data
	_, err = file.Write(xmlData)
	if err != nil {
		return err
	}

	fmt.Printf("XML file '%s' created successfully.\n", filename)
	return nil
}

func DeleteImportableXMLFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}

	fmt.Printf("XML file '%s' deleted successfully.\n", filename)
	return nil
}
