package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2/google"
)

type Prediction struct {
	Prediction int       `json:"prediction"`
	Key        string    `json:"key"`
	Scores     []float64 `json:"scores"`
}

type MLResponseBody struct {
	Predictions []Prediction `json:"predictions"`
}

type ImageBytes struct {
	B64 []byte `json:"b64"`
}

type Instance struct {
	ImageBytes ImageBytes `json:"image_bytes"`
	Key        string     `json:"key"`
}

type MLRequestBody struct {
	Instances []Instance `json:"instances"`
}

const (
	// Replace this project ID and model name with your configuration.
	PROJECT = "around-250315"
	MODEL   = "face"
	URL     = "https://ml.googleapis.com/v1/projects/" + PROJECT + "/models/" + MODEL + ":predict"
	SCOPE   = "https://www.googleapis.com/auth/cloud-platform"
)

// Annotate a image file based on ml model, return score and error if exists.
func annotate(r io.Reader) (float64, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Printf("Cannot read image data %v\n", err)
		return 0.0, err
	}

	client, err := google.DefaultClient(context.Background(), SCOPE)
	if err != nil {
		fmt.Printf("Failed to create HTTP client %v\n", err)
		return 0.0, err
	}

	// Construct a ML request
	requestBody := &MLRequestBody{
		Instances: []Instance{
			{
				ImageBytes: ImageBytes{
					B64: buf,
				},
				Key: "1",
			},
		},
	}
	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Failed to create ML request body %v\n", err)
		return 0.0, err
	}

	request, err := http.NewRequest("POST", URL, strings.NewReader(string(jsonRequestBody)))

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Failed to send ML request %v\n", err)
		return 0.0, err
	}

	jsonResponseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to get ML response body %v\n", err)
		return 0.0, err
	}

	if len(jsonResponseBody) == 0 {
		fmt.Println("Empty prediction response body")
		return 0.0, errors.New("Empty prediction response body")
	}

	var responseBody MLResponseBody
	if err := json.Unmarshal(jsonResponseBody, &responseBody); err != nil {
		fmt.Printf("Failed to decode ML response %v\n", err)
		return 0.0, err
	}

	if len(responseBody.Predictions) == 0 {
		fmt.Println("Empty prediction result")
		return 0.0, errors.New("Empty prediction result")
	}

	results := responseBody.Predictions[0]
	fmt.Printf("Received a prediction result %f\n", results.Scores[0])
	return results.Scores[0], nil
}
