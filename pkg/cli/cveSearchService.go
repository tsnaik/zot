package cli

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	zotErrors "github.com/anuvu/zot/errors"
	resty "github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
)

type CveSearchService interface {
	findCveByImageName(imageName, serverURL, username, password string) (CVEListForImageStruct, error)
	findImagesByCveID(cveID, serverURL, username, password string) (ImageListForCVEStruct, error)
}
type searchService struct{}

func NewCveSearchService() CveSearchService {
	return new(searchService)
}

func makeGraphQLRequestBasicAuth(serverURL, query, username,
	password string, resultStructPointer interface{}) error {
	newURL, err := createEndpointURL(serverURL)
	if err != nil {
		return err
	}

	client := resty.New().
		SetDisableWarn(true)

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(query).
		SetBasicAuth(username, password).
		Post(newURL)

	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		if resp.StatusCode() == http.StatusUnauthorized {
			return zotErrors.ErrUnauthorizedAccess
		}

		return errors.New(resp.String())
	}

	err = client.JSONUnmarshal(resp.Body(), resultStructPointer)

	if err != nil {
		return err
	}

	return nil
}

func (service searchService) findCveByImageName(imageName, serverURL,
	username, password string) (CVEListForImageStruct, error) {
	query := fmt.Sprintf(`{ "query": "{ CVEListForImage (image:\"%s\" )`+
		` { Tag CVEList { Id Description Severity } } }" }`, imageName)
	result := &CVEListForImageStruct{}

	if err := makeGraphQLRequestBasicAuth(serverURL, query, username, password, result); err != nil {
		return CVEListForImageStruct{}, err
	}

	return *result, nil
}

type CVEListForImageStruct struct {
	Data struct {
		CVEListForImage []struct {
			Tag     string `json:"Tag"`
			CVEList []struct {
				ID          string `json:"Id"`
				Description string `json:"Description"`
				Severity    string `json:"Severity"`
			} `json:"CVEList"`
		} `json:"CVEListForImage"`
	} `json:"data"`
}

func (c CVEListForImageStruct) String() string {
	stringBuilder := &strings.Builder{}

	for _, image := range c.Data.CVEListForImage {
		stringBuilder.WriteString("Tag:" + image.Tag + "\n")
		stringBuilder.WriteString("CVE List:\n")

		table := getTableWriter(stringBuilder)

		for _, cve := range image.CVEList {
			row := []string{cve.ID, cve.Severity, cve.Description}
			table.Append(row)
		}

		table.Render()
		stringBuilder.WriteString("\n")
	}

	return stringBuilder.String()
}

func (service searchService) findImagesByCveID(cveID, serverURL, username,
	password string) (ImageListForCVEStruct, error) {
	query := fmt.Sprintf(`{ "query": "{ ImageListForCVE (id:\"%s\" ) { Name Tags } }" }`, cveID)
	result := &ImageListForCVEStruct{}

	if err := makeGraphQLRequestBasicAuth(serverURL, query, username, password, result); err != nil {
		return ImageListForCVEStruct{}, err
	}

	return *result, nil
}

func (c ImageListForCVEStruct) String() string {
	stringBuilder := &strings.Builder{}

	for _, image := range c.Data.ImageListForCVE {
		for _, tag := range image.Tags {
			stringBuilder.WriteString(fmt.Sprintf("%s:%s\n", image.Name, tag))
		}
	}

	return stringBuilder.String()
}

type ImageListForCVEStruct struct {
	Data struct {
		ImageListForCVE []struct {
			Name string   `json:"Name"`
			Tags []string `json:"Tags"`
		} `json:"ImageListForCVE"`
	} `json:"data"`
}

func createEndpointURL(rawURL string) (string, error) {
	if !isURL(rawURL) {
		return "", zotErrors.ErrInvalidURL
	}

	newURL, err := url.Parse(rawURL)

	if err != nil {
		return "", zotErrors.ErrInvalidURL
	}

	newURL, _ = newURL.Parse("/query")

	return newURL.String(), nil
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
} // copied from https://stackoverflow.com/a/55551215

func getTableWriter(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("\n")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(false)

	return table
}
