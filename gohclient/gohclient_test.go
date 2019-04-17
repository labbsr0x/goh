package gohclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	defaultUrl, e := url.Parse("http://localhost:7070")
	if e != nil {
		t.Fatal(e)
	}

	type args struct {
		httpClient *http.Client
		baseURL    string
	}
	tests := []struct {
		name    string
		args    args
		want    *Default
		wantErr bool
	}{
		{
			name:    "use http.DefaultClient when nil httpClient parameter",
			args:    args{nil, defaultUrl.String()},
			want:    &Default{BaseURL: defaultUrl, HTTPClient: http.DefaultClient},
			wantErr: false,
		},
		{
			name:    "use http.Client passed as parameter",
			args:    args{&http.Client{Timeout: time.Hour}, defaultUrl.String()},
			want:    &Default{BaseURL: defaultUrl, HTTPClient: &http.Client{Timeout: time.Hour}},
			wantErr: false,
		},
		{
			name:    "error when baseURL is a empty string ",
			args:    args{&http.Client{Timeout: time.Hour}, ""},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when baseURL is a empty space",
			args:    args{&http.Client{Timeout: time.Hour}, " "},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error when invalid baseURL parameter",
			args:    args{&http.Client{Timeout: time.Hour}, "::"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.httpClient, tt.args.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_request(t *testing.T) {
	type fields struct {
		UserAgent   string
		ContentType string
		Accept      string
		baseURL     *url.URL
		httpClient  *http.Client
	}
	type args struct {
		path   string
		method string
		body   []byte
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantData   []byte
		wantErr    bool
	}{
		{
			name: "error parse url",
			fields: fields{
				baseURL: &url.URL{Host: "localhost"},
			},
			args:    args{path: "::"},
			wantErr: true,
		},
		{
			name: "error creating request, invalid method",
			fields: fields{
				baseURL: &url.URL{Host: "localhost"},
			},
			args:    args{method: "INVALID METHOD"},
			wantErr: true,
		},
		{
			name: "error executing request",
			fields: fields{
				baseURL:    &url.URL{Scheme: "http", Host: "127.0.0.1:0"},
				httpClient: http.DefaultClient,
			},
			wantErr: true,
		},
		{
			name: "error executing request - unsupported protocol scheme",
			fields: fields{
				baseURL:    &url.URL{Scheme: "unsupported", Host: "127.0.0.1:7070"},
				httpClient: http.DefaultClient,
			},
			wantErr: true,
		},
		{
			name: "empty http method must do a GET",
			fields: fields{
				httpClient:  http.DefaultClient,
				UserAgent:   "gohclient",
				Accept:      "application/json",
				ContentType: "application/json",
			},
			args:       args{body: []byte(`{"hello": "world"}`)},
			wantStatus: http.StatusPaymentRequired,
			wantData:   []byte(`{"world": "hello"}`),
			wantErr:    false,
		},
		{
			name: "post json",
			fields: fields{
				httpClient:  http.DefaultClient,
				UserAgent:   "gohclient",
				Accept:      "application/json",
				ContentType: "application/json",
			},
			args:       args{path: "post", method: "POST", body: []byte(`{"hello": "world"}`)},
			wantStatus: http.StatusPaymentRequired,
			wantData:   []byte(`{"world": "hello"}`),
			wantErr:    false,
		},
		{
			name: "PUT text",
			fields: fields{
				httpClient:  http.DefaultClient,
				UserAgent:   "gohclient",
				Accept:      "text/plain",
				ContentType: "text/plain",
			},
			args:       args{path: "post", method: "PUT", body: []byte(`{"hello": "world"}`)},
			wantStatus: http.StatusPaymentRequired,
			wantData:   []byte(`{"world": "hello"}`),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Default{
				UserAgent:   tt.fields.UserAgent,
				ContentType: tt.fields.ContentType,
				Accept:      tt.fields.Accept,
				BaseURL:     tt.fields.baseURL,
				HTTPClient:  tt.fields.httpClient,
			}

			// test do request if no error expected
			if !tt.wantErr {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// check path
					if r.URL.Path != "/"+tt.args.path {
						t.Errorf("path on server must be the same as the client, want %s got %s", tt.args.path, r.URL.Path)
					}
					//a GET method is used when an empty string is passed by method value
					if tt.args.method == "" {
						if r.Method != "GET" {
							t.Errorf("expected method 'GET', got %s", r.Method)
						}
					} else if r.Method != tt.args.method {
						t.Errorf("expected method %s, got %s", tt.args.method, r.Method)
					}
					//check content type
					if tt.args.body != nil {
						if r.Header.Get("Content-Type") != tt.fields.ContentType {
							t.Errorf("expected Content-Type header value %s, got %s", tt.fields.ContentType, r.Header.Get("Content-Type"))
						}
					}

					//check UserAgent
					if r.Header.Get("User-Agent") != tt.fields.UserAgent {
						t.Errorf("expected User-Agent header value %s, got %s", tt.fields.UserAgent, r.Header.Get("User-Agent"))
					}

					//check Accept
					if r.Header.Get("Accept") != tt.fields.Accept {
						t.Errorf("expected Accept header value %s, got %s", tt.fields.Accept, r.Header.Get("Accept"))
					}

					//check body
					if body, _ := ioutil.ReadAll(r.Body); string(body) != string(tt.args.body) {
						t.Errorf("request body = %v, want %v", string(body), string(tt.args.body))
					}

					//write response status code and body
					w.WriteHeader(tt.wantStatus)
					_, err := fmt.Fprint(w, string(tt.wantData))
					if err != nil {
						t.Fatal(err)
					}
				}))
				defer server.Close()

				parsed, err := url.Parse(server.URL)
				if err != nil {
					t.Fatal(err)
				}
				c.BaseURL = parsed
			}

			gotResp, gotData, err := c.request(tt.args.path, tt.args.method, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if gotResp != nil {
					t.Error("the response must be nil")
				}
				if gotData != nil {
					t.Error("the data value must be nil")
				}
			} else {
				if gotResp.StatusCode != tt.wantStatus {
					t.Errorf("Client.request() gotResp = %v, want %v", gotResp.StatusCode, tt.wantStatus)
				}
			}

			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Client.request() gotData = %v, want %v", gotData, tt.wantData)
			}
		})
	}
}

func TestClient_Put(t *testing.T) {
	data := []byte(`{"hello": "world"}`)
	path := "test/put"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected method PUT, got %s", r.Method)
		}
		if body, _ := ioutil.ReadAll(r.Body); string(body) != string(data) {
			t.Errorf("request body = %v, want %v", string(body), string(data))
		}
		if r.URL.Path != "/"+path {
			t.Errorf("path on server must be the same as the client, want %s got %s", path, r.URL.Path)
		}
	}))
	defer server.Close()

	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Default{
		BaseURL:    parsed,
		HTTPClient: http.DefaultClient,
	}

	if _, _, err := c.Put(path, data); err != nil {
		t.Fatal(err)
	}
}

func TestClient_POST(t *testing.T) {
	data := []byte(`{"hello": "world"}`)
	path := "test/post"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		if body, _ := ioutil.ReadAll(r.Body); string(body) != string(data) {
			t.Errorf("request body = %v, want %v", string(body), string(data))
		}
		if r.URL.Path != "/"+path {
			t.Errorf("path on server must be the same as the client, want %s got %s", path, r.URL.Path)
		}
	}))
	defer server.Close()

	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Default{
		BaseURL:    parsed,
		HTTPClient: http.DefaultClient,
	}

	if _, _, err := c.Post(path, data); err != nil {
		t.Fatal(err)
	}
}

func TestClient_GET(t *testing.T) {
	path := "test/get"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/"+path {
			t.Errorf("path on server must be the same as the client, want %s got %s", path, r.URL.Path)
		}
	}))
	defer server.Close()

	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Default{
		BaseURL:    parsed,
		HTTPClient: http.DefaultClient,
	}

	if _, _, err := c.Get(path); err != nil {
		t.Fatal(err)
	}
}

func TestClient_DELETE(t *testing.T) {
	path := "test/delete"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected method DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/"+path {
			t.Errorf("path on server must be the same as the client, want %s got %s", path, r.URL.Path)
		}
	}))
	defer server.Close()

	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := &Default{
		BaseURL:    parsed,
		HTTPClient: http.DefaultClient,
	}

	if _, _, err := c.Delete(path); err != nil {
		t.Fatal(err)
	}
}
