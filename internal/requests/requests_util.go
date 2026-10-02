package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func PostRequest(ctx context.Context, to string, body any) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return []byte{}, errors.New("failed to marshal request body")
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		to,
		bytes.NewReader(jsonBody),
	)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return []byte{}, fmt.Errorf("request failed with status %s", res.Status)
	}

	return responseBody, nil
}

func checkURLValuesError(values url.Values) error {
	if responseError := values.Get("error"); responseError != "" {
		if description := values.Get("error_description"); description != "" {
			return fmt.Errorf("request failed: %s; description: %s", responseError, description)
		}
		return fmt.Errorf("request failed: %s", responseError)
	}
	return nil
}

func requireURLValues(values url.Values, fields ...string) error {
	for _, field := range fields {
		if !values.Has(field) {
			return fmt.Errorf("response missing required field %q", field)
		}
	}
	return nil
}
