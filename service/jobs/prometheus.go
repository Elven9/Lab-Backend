package jobs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// PrometheusValue ,
type PrometheusValue struct {
	UnixTime interface{} `json:"unix"`
	Value    interface{} `json:"value"`
}

func (j *Job) query(str string) ([]byte, error) {
	response, err := http.Get(fmt.Sprintf("http://prometheus-port-fowarding:9090/api/v1/query?query=%s", str))

	if err != nil {
		return nil, fmt.Errorf("[Prometheus.query] Access Prometheus Server Timeout")
	}

	defer response.Body.Close()

	if response.StatusCode == 400 {
		return nil, fmt.Errorf("[Prometheus.query] Bad Request")
	}

	// Get Information
	respByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("[Prometheus.query] Read From Response Body Failed")
	}

	return respByte, nil
}

// GetCPUUsageOvertime .
func (j *Job) GetCPUUsageOvertime() ([]PrometheusValue, error) {
	// Parse All Information
	bytes, err := j.query("kube_pod_container_resource_requests{namespace!=\"monitoring\",namespace!=\"kube-system\",resource=\"cpu\"}[2h]")

	if err != nil {
		return nil, err
	}

	var info struct {
		Data struct {
			Result []struct {
				Metric struct {
					Pod string `json:"pod"`
				} `json:"metric"`
				Values [][2]interface{} `json:"values"`
			} `json:"result"`
		} `json:"data"`
	}

	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return nil, fmt.Errorf("[Prometheus.GetCPUUsage] Unmarshal Value Failed %v", err)
	}

	var result []PrometheusValue
	aggregate := make(map[float64]int64)

	for _, pod := range info.Data.Result {
		if strings.Split(pod.Metric.Pod, "-")[0] != strings.Split(j.Name, "/")[1] {
			continue
		}

		for _, v := range pod.Values {
			// Convert Type
			time := v[0].(float64)
			cores, _ := strconv.ParseInt(v[1].(string), 10, 64)

			aggregate[time] += cores
		}
	}

	// Construct Result
	for k, v := range aggregate {
		result = append(result, PrometheusValue{
			UnixTime: k,
			Value:    v,
		})
	}

	return result, nil
}
