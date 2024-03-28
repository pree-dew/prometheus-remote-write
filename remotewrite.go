package remotewrite

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/prometheus/prompb"
)

func sendToRemoteWrite(data *bytes.Buffer, remoteWriteURL string) (*http.Response, error) {
	req, err := http.NewRequest("POST", remoteWriteURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Encoding", "snappy")
	req.Header.Set("Content-Type", "application/x-protobuf")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	return resp, nil
}

func createSnappyWithMetricFamily(mfs []*io_prometheus_client.MetricFamily) (*prompb.WriteRequest, error) {
	var ts []prompb.TimeSeries
	tStamp := time.Now().UnixNano() / int64(time.Millisecond)

	for _, mf := range mfs {
		for _, m := range mf.Metric {
			labels := []prompb.Label{
				{Name: "__name__", Value: mf.GetName()},
			}
			for _, lp := range m.Label {
				labels = append(labels, prompb.Label{
					Name:  lp.GetName(),
					Value: lp.GetValue(),
				})
			}

			var samples []prompb.Sample
			value := 0.0
			switch *mf.Type {
			case io_prometheus_client.MetricType_COUNTER:
				value = m.GetCounter().GetValue()
			case io_prometheus_client.MetricType_GAUGE:
				value = m.GetGauge().GetValue()
			case io_prometheus_client.MetricType_UNTYPED:
				value = m.GetUntyped().GetValue()
			case io_prometheus_client.MetricType_SUMMARY:
				value = m.GetSummary().GetSampleSum()
			case io_prometheus_client.MetricType_HISTOGRAM:
				value = m.GetHistogram().GetSampleSum()

			default:
				log.Fatalf("Unknown metric type: %v", *mf.Type)
			}

			samples = append(samples, prompb.Sample{
				Value:     value,
				Timestamp: tStamp,
			})

			ts = append(ts, prompb.TimeSeries{
				Labels:  labels,
				Samples: samples,
			})
		}
	}

	fmt.Printf("Writing %v metrics at time: %v\n", len(ts), tStamp)
	return &prompb.WriteRequest{Timeseries: ts}, nil
}

func RemoteWrite(remoteWriteURL string, frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for range ticker.C {
		m, err := prometheus.DefaultGatherer.Gather()
		if err != nil {
			log.Fatalf("Failed to gather metrics: %v", err)
		}

		wr, err := createSnappyWithMetricFamily(m)
		if err != nil {
			log.Fatalf("Failed to create snappy with metric family: %v", err)
		}

		data, err := proto.Marshal(wr)
		if err != nil {
			log.Fatalf("unable to marshal protobuf: %v", err)
		}

		buf := &bytes.Buffer{}
		buf.Write(snappy.Encode(nil, data))

		resp, err := sendToRemoteWrite(buf, remoteWriteURL)
		if err != nil {
			log.Fatalf("Failed to send data to remote write endpoint: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			log.Fatalf("Unexpected response status: %s", resp.Status)
		}

		resp.Body.Close()
		log.Println("Data written successfully to Prometheus remote storage")

	}
}
