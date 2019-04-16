package quay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/vbatts/quay-scheduled-build/types"
)

// BuildRequest calls the quay.io swagger API to trigger a build
//
// http status 403 is permission denied.
// http status 200 is that something in the config must not be correct, so you get the GET style output.
// http status 201 is "success" to create the build request
func BuildRequest(bldinfo types.Build) (map[string]interface{}, error) {
	url, err := url.Parse(DefaultURL)
	if err != nil {
		return nil, err
	}
	url.Path = filepath.Join(url.Path, "repository", bldinfo.QuayRepo, "build") + "/"

	buf, err := json.Marshal(buildToRef(bldinfo))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bldinfo.Token))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := map[string]interface{}{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ret)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return ret, fmt.Errorf("failed to create build (%q)", resp.Status)
	}

	return ret, nil
}

// hacky way to downsample to just the BuildRef struct
func buildToRef(bldinfo types.Build) types.BuildRef {
	buf, _ := json.Marshal(bldinfo)
	bldrf := types.BuildRef{}
	_ = json.Unmarshal(buf, &bldrf)
	return bldrf
}
