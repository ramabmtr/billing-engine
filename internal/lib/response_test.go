package lib

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseSuccess(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		envelopes []string
		expected  Response
	}{
		{
			name: "Success with data, no envelope",
			data: map[string]string{
				"key": "value",
			},
			envelopes: nil,
			expected: Response{
				Status: "success",
				Data: map[string]string{
					"key": "value",
				},
			},
		},
		{
			name: "Success with data and envelope",
			data: []string{"item1", "item2"},
			envelopes: []string{"items"},
			expected: Response{
				Status: "success",
				Data: map[string]any{
					"items": []string{"item1", "item2"},
				},
			},
		},
		{
			name:      "Success with nil data",
			data:      nil,
			envelopes: nil,
			expected: Response{
				Status: "success",
				Data:   nil,
			},
		},
		{
			name:      "Success with string data",
			data:      "test data",
			envelopes: nil,
			expected: Response{
				Status: "success",
				Data:   "test data",
			},
		},
		{
			name:      "Success with integer data",
			data:      123,
			envelopes: nil,
			expected: Response{
				Status: "success",
				Data:   123,
			},
		},
		{
			name:      "Success with string data and envelope",
			data:      "test data",
			envelopes: []string{"message"},
			expected: Response{
				Status: "success",
				Data: map[string]any{
					"message": "test data",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResponseSuccess(tt.data, tt.envelopes...)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.Data, result.Data)
		})
	}
}

func TestResponseError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected Response
	}{
		{
			name: "Error with message",
			err:  errors.New("test error"),
			expected: Response{
				Status:  "error",
				Message: "test error",
			},
		},
		{
			name: "Error with empty message",
			err:  errors.New(""),
			expected: Response{
				Status:  "error",
				Message: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResponseError(tt.err)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.Message, result.Message)
			assert.Nil(t, result.Data)
		})
	}
}