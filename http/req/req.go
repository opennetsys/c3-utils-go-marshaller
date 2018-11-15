package req

import (
	"bytes"
	"crypto/tls"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

// REQ_FILENAME is where the request should be written to
const REQ_FILENAME = "req_bytes.txt"

// TransformedRequest is used to marshall http requests
type TransformedRequest struct {
	Method           string
	URL              url.URL
	Proto            string // "HTTP/1.0"
	ProtoMajor       int    // 1
	ProtoMinor       int    // 0
	Header           http.Header
	BodyBytes        []byte
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Host             string
	Form             url.Values
	PostForm         url.Values // Go 1.1
	MultipartForm    multipart.Form
	Trailer          http.Header
	RemoteAddr       string
	RequestURI       string
	TLS              tls.ConnectionState
	Response         http.Response
}

// Marshal converts a TransformedRequest to a bytes array
func (tr *TransformedRequest) Marshal() ([]byte, error) {
	var reqBytes bytes.Buffer
	enc := gob.NewEncoder(&reqBytes)

	if err := enc.Encode(tr); err != nil {
		return nil, err
	}

	return reqBytes.Bytes(), nil
}

// Unmarshal converts a bytes array to a TransformedRequest
func (tr *TransformedRequest) Unmarshal(b []byte) error {
	reqBytes := bytes.NewBuffer(b)
	dec := gob.NewDecoder(reqBytes)

	return dec.Decode(tr)
}

// UnTransformRequest takes a TransformedRequest and turns it into an http request
func UnTransformRequest(tr *TransformedRequest) (*http.Request, error) {
	if tr == nil {
		log.Println("received a nil transformed request")
		return nil, errors.New("received a nil transformed request")
	}

	r := &http.Request{
		Method:           tr.Method,
		Proto:            tr.Proto,
		ProtoMajor:       tr.ProtoMajor,
		ProtoMinor:       tr.ProtoMinor,
		Header:           tr.Header,
		ContentLength:    tr.ContentLength,
		TransferEncoding: tr.TransferEncoding,
		Close:            tr.Close,
		Host:             tr.Host,
		Form:             tr.Form,
		PostForm:         tr.PostForm,
		Trailer:          tr.Trailer,
		RemoteAddr:       tr.RemoteAddr,
		// RequestURI:       tr.RequestURI,
	}

	if !reflect.DeepEqual(tr.URL, url.URL{}) {
		r.URL = &tr.URL
	}
	if !reflect.DeepEqual(tr.MultipartForm, multipart.Form{}) {
		r.MultipartForm = &tr.MultipartForm
	}
	if !reflect.DeepEqual(tr.TLS, tls.ConnectionState{}) {
		r.TLS = &tr.TLS
	}
	if !reflect.DeepEqual(tr.Response, http.Response{}) {
		r.Response = &tr.Response
	}
	if tr.BodyBytes != nil || len(tr.BodyBytes) != 0 {
		r.Body = ioutil.NopCloser(bytes.NewReader(tr.BodyBytes))
	}

	return r, nil
}

// TransformRequest takes an http request and turns it into a TransformedRequest for marshalling
func TransformRequest(r *http.Request) (*TransformedRequest, error) {
	if r == nil {
		log.Println("received a nil http req")
		return nil, errors.New("received a nil http req")
	}

	tr := &TransformedRequest{
		Method:           r.Method,
		Proto:            r.Proto,
		ProtoMajor:       r.ProtoMajor,
		ProtoMinor:       r.ProtoMinor,
		Header:           r.Header,
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncoding,
		Close:            r.Close,
		Host:             r.Host,
		Form:             r.Form,
		PostForm:         r.PostForm,
		Trailer:          r.Trailer,
		RemoteAddr:       r.RemoteAddr,
		// RequestURI:       r.RequestURI,
	}

	if r.URL != nil {
		tr.URL = *r.URL
	}
	if r.MultipartForm != nil {
		tr.MultipartForm = *r.MultipartForm
	}
	if r.TLS != nil {
		tr.TLS = *r.TLS
	}
	if r.Response != nil {
		tr.Response = *r.Response
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("err reading body\n%v", err)

		return nil, err
	}
	r.Body.Close()                                        //  must close
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // reading destroyed the body so we have to re-write it

	tr.BodyBytes = bodyBytes

	return tr, nil
}

// WriteReqToFile writes an http request to a file
func WriteReqToFile(r *http.Request, filename string) error {
	if r == nil {
		return errors.New("request is nil")
	}

	tr, err := TransformRequest(r)
	if err != nil {
		return err
	}

	b, err := tr.Marshal()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// ReadReqFromFile reads a req from file
func ReadReqFromFile(filename string) (*http.Request, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	tr := new(TransformedRequest)
	if err := tr.Unmarshal(b); err != nil {
		return nil, err
	}

	return UnTransformRequest(tr)
}
