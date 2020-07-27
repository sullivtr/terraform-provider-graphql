package graphql

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func queryExecute(d *schema.ResourceData, m interface{}, querySource, variableSource string) ([]byte, error) {
	query := d.Get(querySource).(string)
	variables := d.Get(variableSource).(map[string]interface{})
	apiURL := m.(*graphqlProviderConfig).GQLServerUrl
	headers := m.(*graphqlProviderConfig).RequestHeaders

	var queryBodyBuffer bytes.Buffer

	queryObj := gqlQuery{
		Query:     query,
		Variables: variables,
	}

	if err := json.NewEncoder(&queryBodyBuffer).Encode(queryObj); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, &queryBodyBuffer)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value.(string))
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}
