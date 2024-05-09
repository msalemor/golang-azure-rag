package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CallApiWithApiKey(url, key string, payload any) ([]byte, error) {

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	fmt.Println("POST", url, key, string(payloadBytes))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("api-key", key)
	req.Header.Add("Content-Type", "application/json")

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API call failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func CallApi[T any](url, key string, payload any) (T, error) {
	bytes, err := CallApiWithApiKey(url, key, payload)
	var gptResponse T

	if err != nil {
		return gptResponse, err
	}

	// Convert bytes to string
	jsonString := string(bytes)

	// Unmarshal the JSON string into a GPTResponse struct

	err = json.Unmarshal([]byte(jsonString), &gptResponse)
	if err != nil {
		return gptResponse, err
	}

	return gptResponse, nil
}

func CallGPT(payload GPTRequest) (GPTResponse, error) {
	resp, err := CallApi[GPTResponse](gptFullEndpoint, gptApiKey, payload)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func CallAISearch(payload AISearchRequest) (AISearchResponse, error) {
	resp, err := CallApi[AISearchResponse](aiSearchEndpoint, aiSearchApiKey, payload)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func Intent(question string) (string, error) {
	template := `system:
You can determine intent from the following list:

RAG: Questions related to product infomrmation.
Other: Anything else.

user:
<Question>

Output in the one word intent only.
`
	message := Message{Role: "assistant", Content: template}
	messages := []Message{message}
	var MaxTokens = 2
	payload := GPTRequest{Messages: messages, MaxTokens: &MaxTokens, Temperature: 0.2}
	resp, err := CallGPT(payload)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil

}
