package importers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	uipack "github.com/uipack-io/cli"
)

type FigmaImporter struct {
	Host        string
	AccessToken string
}

func (i *FigmaImporter) Decode(p *uipack.Package, fileKey string) error {
	url := fmt.Sprintf("%s/v1/files/%s/variables/local", i.Host, fileKey)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+i.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch variables: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var figmaResponse figmaResponse
	err = json.Unmarshal(body, &figmaResponse)
	if err != nil {
		return err
	}

	return nil
}

type figmaVariable struct {
	Id           string                 `json:"id"`
	Name         string                 `json:"name"`
	ResolvedType string                 `json:"resolvedType"`
	ValuesByMode map[string]interface{} `json:"valuesByMode"`
}

type figmaMode struct {
	Id   string `json:"modeId"`
	Name string `json:"name"`
}

type figmaVariableCollection struct {
	Id            string      `json:"id"`
	Name          string      `json:"name"`
	DefaultModeId string      `json:"defaultModeId"`
	Modes         []figmaMode `json:"modes"`
	VariableIds   []string    `json:"variableIds"`
}

type figmaMeta struct {
	Variables   map[string]figmaVariable           `json:"variables"`
	Collections map[string]figmaVariableCollection `json:"variableCollections"`
}

type figmaResponse struct {
	Meta figmaMeta `json:"meta"`
}

func (f *figmaResponse) ToBundleMetadata() *uipack.BundleMetadata {
	result := uipack.BundleMetadata{
		Name: "Figma",
		Version: uipack.Version{
			Major: 1,
			Minor: 0,
		},
	}
	vi := 0
	ci := 0
	for _, collection := range f.Meta.Collections {
		mm := uipack.ModeMetadata{
			Identifier: uipack.Uint4(ci),
			Name:       collection.Name,
		}
		for i, mode := range collection.Modes {
			mm.Variants = append(mm.Variants, uipack.ModeVariantMetadata{
				Identifier: uipack.Uint4(i),
				Name:       mode.Name,
			})
			if i == 0 {
				for _, variableId := range collection.VariableIds {
					variable := f.Meta.Variables[variableId]
					vm := uipack.VariableMetadata{
						Identifier: uint64(vi),
						Name:       variable.Name,
						Type:       figmaToVariableType(variable.ResolvedType),
					}
					result.Variables = append(result.Variables, vm)
					vi = vi + 1
				}
			}
			ci++
		}

		result.Modes = append(result.Modes, mm)
	}
	return &result
}
