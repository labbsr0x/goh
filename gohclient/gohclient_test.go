package gohclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	c := "application/json"
	a := New(c)
	if a.ContentType != c {
		t.Errorf("Expecting API to have %v as the content type; got %v", c, a.ContentType)
	}
}

func TestPut(t *testing.T) {
	a := New("application/json")
	testRequestWithPayload(a.Put, t)
}

func TestPost(t *testing.T) {
	a := New("application/json")
	testRequestWithPayload(a.Post, t)
}

func TestGet(t *testing.T) {
	a := New("application/json")
	testRequestWithNoPayload(a.Get, t)
}

func TestDelete(t *testing.T) {
	a := New("application/json")
	testRequestWithNoPayload(a.Delete, t)
}

func testRequestWithNoPayload(do func(url string) (*http.Response, []byte, error), t *testing.T) (*http.Response, []byte, error) {
	httpDo = func(req *http.Request) (*http.Response, error) { return getMockResponse(true), nil }

	resp, data, err := do("/test")

	if resp == nil || data == nil || err != nil {
		t.Errorf("Expecting non nil http.Response, non nil data and nil err. Got %v http.Response, %v data and %v err", resp, data, err)
	}

	var r MockResp
	err = json.Unmarshal(data, &r)

	if err != nil {
		t.Errorf("Expecting a MockResp as response. Got the following error instead: %v", err)
	}

	if !r.Success {
		t.Errorf("Expecting the response do be successfull")
	}

	return resp, data, err
}

func testRequestWithPayload(do func(url string, data []byte) (*http.Response, []byte, error), t *testing.T) (*http.Response, []byte, error) {
	httpDo = func(req *http.Request) (*http.Response, error) {
		d, _ := ioutil.ReadAll(req.Body)
		return getMockResponseWithPayload(d), nil
	}

	m := MockData{Info: "just a simple test", Ok: true}
	d, _ := json.Marshal(m)
	resp, data, err := do("/test", d)

	if resp == nil || data == nil || err != nil {
		t.Errorf("Expecting non nil http.Response, non nil data and nil err. Got %v http.Response, %v data and %v err", resp, data, err)
	}

	var r MockData
	err = json.Unmarshal(data, &r)

	if err != nil {
		t.Errorf("Expecting a MockData as response. Got the following error instead: %v", err)
	}

	if r != m {
		t.Errorf("Expecting response to be exactly the request payload. Got '%v' instead", r)
	}

	return resp, data, err
}

func getMockResponseWithPayload(payload []byte) *http.Response {
	w := httptest.NewRecorder()
	w.Write(payload)
	return w.Result()
}

func getMockResponse(success bool) *http.Response {
	w := httptest.NewRecorder()
	d, _ := json.Marshal(MockResp{Success: success})
	w.Write(d)
	return w.Result()
}

type MockData struct {
	Info string `json:"info"`
	Ok   bool   `json:"ok"`
}

type MockResp struct {
	Success bool `json:"success"`
}
