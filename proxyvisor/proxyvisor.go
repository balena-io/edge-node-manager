package proxyvisor

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/cavaliercoder/grab"

	"github.com/josephroberts/edge-node-manager/application"
	"github.com/josephroberts/edge-node-manager/config"
	"github.com/josephroberts/edge-node-manager/device"
	"github.com/parnurzeal/gorequest"
)

var (
	address string
	version string
	key     string
	rawKey  string
)

// DependantApplicationsList returns all dependant applications assigned to the edge-node-manager
func DependantApplicationsList() (map[string]*application.Application, []error) {
	url, err := buildPath(address, []string{version, "applications"})
	if err != nil {
		return nil, []error{err}
	}

	req := gorequest.New()
	req.Get(url)
	req.Query(key)

	log.WithFields(log.Fields{
		"URL":    req.Url,
		"Method": req.Method,
		"Query":  req.QueryData,
	}).Debug("Requesting dependant applications list")

	resp, body, errs := req.EndBytes()
	if errs = handleResp(resp, errs); errs != nil {
		return nil, errs
	}

	var buffer []application.Application
	if err := json.Unmarshal(body, &buffer); err != nil {
		return nil, []error{err}
	}

	applications := make(map[string]*application.Application)
	for key, app := range buffer {
		// The value of the map must be set using the key rather than using the app object directly,
		// otherwise the memory address pointed to gets overwritten with each iteration of the loop
		applications[app.Name] = &buffer[key]
	}

	return applications, nil
}

// DependantApplicationUpdate downloads the binary.tar for a specific application and commit,
// saving it to {ENM_ASSETS_DIRECTORY}/{appUUID}/{commit}/binary.tar
func DependantApplicationUpdate(appUUID int, commit string) error {
	url, err := buildPath(address, []string{version, "assets", strconv.Itoa(appUUID), commit})
	if err != nil {
		return err
	}

	req, err := grab.NewRequest(url)
	if err != nil {
		return err
	}

	q := req.HTTPRequest.URL.Query()
	q.Set("apikey", rawKey)
	req.HTTPRequest.URL.RawQuery = q.Encode()

	filePath := config.GetAssetsDir()
	filePath = path.Join(filePath, strconv.Itoa(appUUID))
	filePath = path.Join(filePath, commit)
	if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
		return err
	}
	filePath = path.Join(filePath, "binary.tar")
	req.Filename = filePath

	log.WithFields(log.Fields{
		"URL":         req.HTTPRequest.URL,
		"Method":      req.HTTPRequest.Method,
		"Query":       req.HTTPRequest.URL.RawQuery,
		"Destination": req.Filename,
	}).Debug("Requesting dependant application update")

	client := grab.NewClient()
	_, err = client.Do(req)

	return err
}

// DependantDeviceLog transmits a log message and timestamp for a specific device
func DependantDeviceLog(resinUUID, message string) []error {
	url, err := buildPath(address, []string{version, "devices", resinUUID, "logs"})
	if err != nil {
		return []error{err}
	}

	type dependantDeviceLog struct {
		Message   string `json:"message"`
		TimeStamp int64  `json:"timestamp"`
	}

	content := &dependantDeviceLog{
		Message:   message,
		TimeStamp: time.Now().UTC().Unix(),
	}

	bytes, err := json.Marshal(content)
	if err != nil {
		return []error{err}
	}

	req := gorequest.New()
	req.Put(url)
	req.Query(key)
	req.Send((string)(bytes))

	log.WithFields(log.Fields{
		"URL":    req.Url,
		"Method": req.Method,
		"Query":  req.QueryData,
		"Body":   (string)(bytes),
	}).Debug("Transmitting dependant device log")

	resp, _, errs := req.End()
	return handleResp(resp, errs)
}

// DependantDeviceInfoUpdate transmits status and is_online for a specific device
func DependantDeviceInfoUpdate(resinUUID, status string, online bool) []error {
	url, err := buildPath(address, []string{version, "devices", resinUUID})
	if err != nil {
		return []error{err}
	}

	type dependantDeviceInfo struct {
		Status string `json:"status"`
		Online bool   `json:"is_online"`
	}

	content := &dependantDeviceInfo{
		Status: status,
		Online: online,
	}

	bytes, err := json.Marshal(content)
	if err != nil {
		return []error{err}
	}

	req := gorequest.New()
	req.Put(url)
	req.Query(key)
	req.Send((string)(bytes))

	log.WithFields(log.Fields{
		"URL":    req.Url,
		"Method": req.Method,
		"Query":  req.QueryData,
		"Body":   (string)(bytes),
	}).Debug("Transmitting dependant device info")

	resp, _, errs := req.End()
	return handleResp(resp, errs)
}

// DependantDeviceInfo returns a single dependant device assigned to the edge-node-manager
// NOT USED
func DependantDeviceInfo(resinUUID string) (device.Device, error) {
	return device.Device{}, fmt.Errorf("Not implemented")
}

// DependantDeviceProvision provisions a single dependant device to a specific application
func DependantDeviceProvision(appUUID int) (string, string, []error) {
	url, err := buildPath(address, []string{version, "devices"})
	if err != nil {
		return "", "", []error{err}
	}

	type dependantDeviceProvision struct {
		AppUUID int `json:"appId"`
	}

	content := &dependantDeviceProvision{
		AppUUID: appUUID,
	}

	bytes, err := json.Marshal(content)
	if err != nil {
		return "", "", []error{err}
	}

	req := gorequest.New()
	req.Post(url)
	req.Query(key)
	req.Send((string)(bytes))

	log.WithFields(log.Fields{
		"URL":    req.Url,
		"Method": req.Method,
		"Query":  req.QueryData,
		"Body":   (string)(bytes),
	}).Debug("Requesting dependant device provision")

	resp, body, errs := req.EndBytes()
	if errs = handleResp(resp, errs); errs != nil {
		return "", "", errs
	}

	var buffer map[string]interface{}
	if err := json.Unmarshal(body, &buffer); err != nil {
		return "", "", []error{err}
	}

	return buffer["uuid"].(string), buffer["name"].(string), nil
}

// DependantDevicesList returns all dependant devices assigned to the edge-node-manager
// NOT USED
func DependantDevicesList() (map[string]device.Device, error) {
	return nil, fmt.Errorf("Not implemented")
}

func init() {
	address = config.GetSuperAddr()
	version = config.GetSuperAPIVer()
	rawKey = config.GetSuperAPIKey()

	type apiKey struct {
		APIKey string `json:"apikey"`
	}

	content := &apiKey{
		APIKey: rawKey,
	}

	bytes, err := json.Marshal(content)
	if err != nil {
		log.WithFields(log.Fields{
			"Key":   rawKey,
			"Error": err,
		}).Fatal("Unable to marshall API key")
	}
	key = (string)(bytes)

	log.WithFields(log.Fields{
		"Address": address,
		"Version": version,
		"Key":     key,
		"Raw key": rawKey,
	}).Debug("Initialised outgoing supervisor API")
}

func buildPath(base string, paths []string) (string, error) {
	url, err := url.ParseRequestURI(address)
	if err != nil {
		return "", err
	}

	for _, p := range paths {
		url.Path = path.Join(url.Path, p)
	}

	return url.String(), nil
}

func handleResp(resp gorequest.Response, errs []error) []error {
	if errs != nil {
		return errs
	}

	if resp.StatusCode != 200 {
		return []error{fmt.Errorf("invalid response received: %s", resp.Status)}
	}

	log.WithFields(log.Fields{
		"Response": resp.Status,
	}).Debug("Valid response received")

	return nil
}
