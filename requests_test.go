package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
)

type CommonResponse struct {
	Result interface{} `json:"result"`
	Status bool        `json:"status"`
}

func TestLoadGoods(t *testing.T) {

	file := openFile("assets/book1.xlsx")

	params := map[string]io.Reader{
		"goods_file": file,
		"seller_id":  strings.NewReader("1"),
	}

	client := &http.Client{}

	err, res := makeRequest(client, "http://127.0.0.1:4560/load_goods", params)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestRetrieveGoodsAllParamsPresent(t *testing.T) {
	client := &http.Client{}

	url := "http://127.0.0.1:4560/retrieve_goods?query=prod&seller_id=1&offer_id=1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("New Request: %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Client Do: %v", err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.True(t, res.Body != nil, "Response body is null")
}

func TestRetrieveGoodsNotAllParamsPresent(t *testing.T) {
	client := &http.Client{}

	url := "http://127.0.0.1:4560/retrieve_goods?query=product&offer_id=2"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("New Request: %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Client Do: %v", err)
	}
	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.True(t, res.Body != nil, "Response body is null")
}

func makeRequest(client *http.Client, url string, params map[string]io.Reader) (err error, res *http.Response) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range params {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err, nil
		}
	}
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	return nil, res
}

func openFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return file
}
