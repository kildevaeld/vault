package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"

	"github.com/kildevaeld/vault/vault"
)

type ClientConfig struct {
	ServerConfig interface{}
}

func (self *ClientConfig) Endpoint(path string) string {
	domain := ""

	if tcp, ok := self.ServerConfig.(VaultServerTCPConfig); ok {
		domain = fmt.Sprintf("localhost:%d", tcp.Port)
	} else {
		domain = "vault.socket"
	}
	str := "http://"
	if domain != "" {
		str += domain + "/"
	}

	if path[0] == '/' {
		path = path[1:]
	}
	str += path

	return str

}

func (self *ClientConfig) IsUnix() bool {
	if _, ok := self.ServerConfig.(VaultServerUnixConfig); ok {
		return true
	}
	return false
}

type Client struct {
	client *http.Client
	config *ClientConfig
}

type Query struct {
	Id   string
	Name string
	Tags []string
}

type ErrorMessage struct {
	Error string
	Code  int
}

func (self *Client) doGet(path string) (*http.Response, error) {
	if path[0] == '/' {
		path = path[1:]
	}

	return self.client.Get(self.config.Endpoint(path))

}

func (self *Client) Ping() error {

	res, err := self.doGet("ping")

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	return nil
}

func (self *Client) List() ([]*vault.Item, error) {

	r, e := self.doGet(ListEndpoint)

	if e != nil {
		return nil, e
	}

	if r.StatusCode != 200 {
		return nil, errors.New(r.Status)
	}

	var out []*vault.Item
	buf := bytes.NewBuffer(nil)
	_, e = io.Copy(buf, r.Body)
	defer r.Body.Close()

	if e != nil {
		return nil, e
	}

	e = json.Unmarshal(buf.Bytes(), &out)

	return out, e
}

func (self *Client) Find(glob string) ([]*vault.Item, error) {

	req, e := http.NewRequest("GET", self.config.Endpoint("list"), nil)

	if e != nil {
		fmt.Printf("error")
		return nil, e
	}

	q := req.URL.Query()
	q.Set("query", glob)
	req.URL.RawQuery = q.Encode()

	r, err := self.client.Do(req)

	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, errors.New(r.Status)
	}

	var out []*vault.Item
	buf := bytes.NewBuffer(nil)
	_, e = io.Copy(buf, r.Body)
	defer r.Body.Close()

	if e != nil {
		return nil, e
	}

	e = json.Unmarshal(buf.Bytes(), &out)

	return out, e

}

func (self *Client) Get(id string) (*vault.Item, error) {

	return nil, nil
}

func (self *Client) Reader(id string) (io.ReadCloser, error) {

	req, e := http.NewRequest("GET", self.config.Endpoint("read"), nil)

	if e != nil {
		return nil, e
	}

	q := req.URL.Query()
	q.Set("id", id)
	req.URL.RawQuery = q.Encode()

	res, er := self.client.Do(req)

	if er != nil {
		return nil, er
	}

	if res.StatusCode != http.StatusOK {
		var msg ErrorMessage
		decodeBody(readBody(res.Body), &msg)
		defer res.Body.Close()
		err := fmt.Errorf(msg.Error)
		return nil, err
	}

	return res.Body, nil
}

func (self *Client) Remove(id string) error {

	req, e := http.NewRequest("DELETE", self.config.Endpoint(id), nil)

	if e != nil {
		return e
	}

	res, err := self.client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		var msg ErrorMessage
		decodeBody(readBody(res.Body), &msg)
		defer res.Body.Close()
		err := fmt.Errorf(msg.Error)
		return err
	}
	return nil
}

func (self *Client) Create(r io.Reader, options vault.ItemCreateOptions) (*vault.Item, error) {

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", options.Name)
	if err != nil {
		return nil, err
	}

	if _, e := io.Copy(fw, r); e != nil {
		return nil, e
	}

	if e := addField(w, "name", []byte(options.Name)); e != nil {
		return nil, e
	}

	if options.Key != nil {
		if e := addField(w, "key", options.Key[:]); e != nil {
			return nil, e
		}
	}

	w.Close()

	req, err := http.NewRequest("POST", self.config.Endpoint("upload"), &b)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	//req.Header.Set("Content-Length", fmt.Sprintf("%d", options.Size))

	res, e := self.client.Do(req)

	if e != nil {
		return nil, e
	}

	if res.StatusCode != http.StatusOK {
		var msg ErrorMessage
		decodeBody(readBody(res.Body), &msg)
		defer res.Body.Close()
		err = fmt.Errorf(msg.Error)
		return nil, err
	}

	var item vault.Item
	err = decodeBody(readBody(res.Body), &item)
	defer res.Body.Close()

	return &item, err
}

func decodeBody(b []byte, t interface{}) error {
	return json.Unmarshal(b, t)
}

func readBody(r io.Reader) []byte {

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r)
	return buf.Bytes()
}

func addField(w *multipart.Writer, fieldName string, value []byte) error {

	fw, err := w.CreateFormField(fieldName)
	if err != nil {
		return err
	}

	_, err = fw.Write(value)

	return err

}

func NewVaultClient(config *ClientConfig) (*Client, error) {

	var client *http.Client
	if config.IsUnix() {

		path := config.ServerConfig.(VaultServerUnixConfig).Path

		trans := http.Transport{
			Dial: func(proto, addr string) (conn net.Conn, err error) {
				return net.Dial("unix", path)
			},
		}
		client = &http.Client{
			Transport: &trans,
		}
	} else {
		client = &http.Client{}
	}

	return &Client{client, config}, nil
}
