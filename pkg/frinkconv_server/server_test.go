package frinkconv_server

import (
	"bytes"
	"fmt"
	"github.com/initialed85/frinkconv-api/internal/helpers"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	port, err := helpers.GetRandomPort()
	if err != nil {
		t.Fatal(err)
	}

	server, err := New(port, 2)
	assert.Nil(t, err)

	defer server.Close()

	response, err := client.Post(
		fmt.Sprintf("http://localhost:%v/convert/", port),
		"application/json",
		bytes.NewBuffer([]byte(`{"source_value": 10, "source_units": "apples", "destination_units": "oranges"}`)),
	)
	assert.Nil(t, err)

	response, err = client.Post(
		fmt.Sprintf("http://localhost:%v/convert/", port),
		"application/json",
		bytes.NewBuffer([]byte(`{"source_value": 10, "source_units": "feet", "destination_units": "inches"}`)),
	)
	assert.Nil(t, err)

	responseBody, err := ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(
		t,
		"{\"destination_value\":120}",
		string(responseBody),
	)
	_ = response.Body.Close()

	response, err = client.Post(
		fmt.Sprintf("http://localhost:%v/batch_convert/", port),
		"application/json",
		bytes.NewBuffer([]byte(`[{"source_value": 10, "source_units": "apples", "destination_units": "oranges"}, {"source_value": 10, "source_units": "feet", "destination_units": "inches"}]`)),
	)
	assert.Nil(t, err)

	responseBody, err = ioutil.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(
		t,
		"[{\"error\":\"Warning: undefined symbol \\\"apples\\\".\\nUnknown symbol \\\"oranges\\\"\\nWarning: undefined symbol \\\"apples\\\".\\nWarning: undefined symbol \\\"oranges\\\".\\nUnconvertable expression:\\n  10 apples (undefined symbol) -\\u003e oranges (undefined symbol)\"},{\"destination_value\":120}]",
		string(responseBody),
	)

	defer func() {
		_ = response.Body.Close()
	}()
}
