package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/emman27/jenkinsctl/pkg/builds"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// GetBuild retrieves a particular build of a job
func (c *JenkinsClient) GetBuild(jobName string, buildID int) (*builds.Build, error) {
	resp, err := c.Get(fmt.Sprintf("/job/%s/%d/api/json", jobName, buildID))
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get job %d for %s", buildID, jobName)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not read response body")
	}
	var build = new(builds.Build)
	json.Unmarshal(content, build)
	glog.Infof("Parsed response %v", build)
	return build, nil
}

// CreateBuild starts a build in Jenkins
func (c *JenkinsClient) CreateBuild(jobName string, params map[string]interface{}) (*builds.Build, error) {
	glog.Infof("Creating build %s with parameters %v", jobName, params)
	parameters, err := builds.GenerateParametersBody(params)
	glog.Infof("Parameters: %s", parameters)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to format parameters")
	}
	reader := strings.NewReader(parameters)
	resp, err := c.Post(fmt.Sprintf("/job/%s/build", jobName), reader)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create build")
	}
	glog.Infof("Queued build: %s", resp.Header.Get("Location"))
	return nil, nil
}
