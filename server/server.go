package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kildevaeld/go-filecrypt"
	"github.com/kildevaeld/vault/vault"
	"github.com/mitchellh/mapstructure"
)

const (
	UploadEndpoint = "/upload"
	ListEndpoint   = "/list"
	ReadEndpoint   = "/read"
)

func errorJSONMessage(err error, status int) []byte {
	m := vault.Map{
		"error": err.Error(),
		"code":  status,
	}
	return m.ToJSON()
}

func writeJSON(w io.Writer, what interface{}) (err error) {

	if b, ok := what.([]byte); ok {
		_, err = w.Write(b)
	} else {
		b, e := json.Marshal(what)
		if e != nil {
			err = e
		} else {
			buf := bytes.NewBuffer(nil)
			err = json.Indent(buf, b, "", "  ")
			_, err = w.Write(buf.Bytes())
		}

	}
	return err
}

type VaultServerTCPConfig struct {
	Port uint16
}

type VaultServerUnixConfig struct {
	Path string
}

func GetServerConfig(conf *vault.Config) (interface{}, error) {
	typ := ""
	t := conf.Server["Type"]

	if t == nil {
		t = conf.Server["type"]
	}
	if t != nil {
		typ = t.(string)
	}

	var sConf interface{}
	if strings.ToLower(typ) == "tcp" {
		var config VaultServerTCPConfig
		err := mapstructure.Decode(conf.Server, &config)
		if err != nil {
			return nil, err
		}
		sConf = config
	} else if strings.ToLower(typ) == "unix" {
		var config VaultServerUnixConfig
		mapstructure.Decode(conf.Server, &config)
		sConf = config
	}
	return sConf, nil
}

type VaultServer struct {
	listener net.Listener
	mux      *mux.Router
	vault    *vault.Vault
}

func (s *VaultServer) Listen() error {

	s.mux.HandleFunc("/ping", s.handlePing).Methods("GET")
	s.mux.HandleFunc(UploadEndpoint, s.handleUpload).Methods("POST")
	s.mux.HandleFunc(ListEndpoint, s.handleList).Methods("GET")
	s.mux.HandleFunc(ReadEndpoint, s.handleRead).Methods("GET")
	s.mux.HandleFunc("/{item_id}", s.handleDelete).Methods("DELETE")
	//s.mux.HandleFunc("/{id}", f)
	return http.Serve(s.listener, s.mux)

}

func (s *VaultServer) Close() error {
	return s.listener.Close()
}

func (s *VaultServer) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func handleError(w http.ResponseWriter, e error, status int) {
	w.WriteHeader(status)
	w.Write(errorJSONMessage(e, status))
}

func (s *VaultServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	file, handler, err := r.FormFile("file")
	if err != nil {
		handleError(w, err, http.StatusBadRequest)
		return
	}

	defer file.Close()

	mime := handler.Header.Get("Content-Type")
	sizeStr := handler.Header.Get("Content-Length")

	if mime == "" || mime != "" {
		sample := make([]byte, 1024)
		_, e := file.Read(sample)
		if e != nil {
			handleError(w, e, http.StatusInternalServerError)
			return
		}
		mime, e = vault.DetectContentType(sample)
		_, e = file.Seek(0, 0)

		if e != nil {
			handleError(w, e, http.StatusInternalServerError)
			return
		}
	}
	size := uint64(1)
	if sizeStr == "" {

		i, e := file.Seek(0, 2)

		if e != nil {
			handleError(w, e, http.StatusInternalServerError)
			return
		}

		size = uint64(i)

		_, e = file.Seek(0, 0)

		if e != nil {
			handleError(w, e, http.StatusInternalServerError)
			return
		}

	} else {
		i, e := strconv.ParseInt(sizeStr, 10, 64)
		if e != nil {
			handleError(w, e, http.StatusInternalServerError)
			return
		}
		size = uint64(i)
	}

	name := r.FormValue("name")

	if name == "" {
		name = handler.Filename
	}

	key := r.FormValue("key")
	var k *[32]byte
	if key != "" {
		tmp := filecrypt.Key([]byte(key))
		k = &tmp
	}

	i, e := s.vault.Add(file, vault.ItemCreateOptions{
		Mime: mime,
		Name: name,
		Key:  k,
		Size: size,
	})

	if e != nil {

		handleError(w, e, http.StatusInternalServerError)
	} else {

		b, _ := json.Marshal(i)

		w.Write(b)

	}

}

func (s *VaultServer) handleDelete(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	vars := mux.Vars(r)

	itemId := vars["item_id"]
	fmt.Printf("item id %v", itemId)

	err := s.vault.Remove(itemId)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	m := vault.Map{
		"message": "ok",
		"code":    200,
	}
	b, _ := json.Marshal(m)
	w.Write(b)

}

func (s *VaultServer) handleList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	q := r.URL.Query()

	var items []*vault.Item
	if name := q.Get("query"); name != "" {

		items = s.vault.Find(name)
	} else {
		items = s.vault.List()
	}

	//var b []byte
	var err error
	if first := q.Get("first"); first != "" && len(items) > 0 {
		err = writeJSON(w, items[0])
		//b, _ = json.Marshal(items[0])
	} else {
		err = writeJSON(w, items)
		//b, _ = json.Marshal(items)
	}
	if err != nil {
		//err = writeJSON(w, errorJSONMessage(err, ))
	}

	//buf := bytes.NewBuffer(nil)
	//json.Indent(buf, b, "", "  ")

	//w.Write(buf.Bytes())
}

func (self *VaultServer) handleRead(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	if id == "" {
		handleError(w, errors.New("no id"), http.StatusBadRequest)
		return
	}

	item, err := self.vault.Get(id)

	if err != nil {
		handleError(w, err, http.StatusInternalServerError)
		return
	}

	reader, e := self.vault.Open(item)

	if e != nil {
		handleError(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", item.Mime)

	_, e = io.Copy(w, reader)
	if e != nil {
		handleError(w, e, http.StatusInternalServerError)
		return
	}

}

func NewVaultServer(v *vault.Vault, config interface{}) (*VaultServer, error) {

	var err error

	var listener net.Listener
	if unix, ok := config.(VaultServerUnixConfig); ok {

		listener, err = net.Listen("unix", unix.Path)

	} else if tcp, ok := config.(VaultServerTCPConfig); ok {
		var addr *net.TCPAddr
		addr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", tcp.Port))
		if err == nil {
			listener, err = net.ListenTCP("tcp", addr)
		}
	}

	if err != nil {
		return nil, err
	}

	return &VaultServer{
		listener: listener,
		mux:      mux.NewRouter(),
		vault:    v,
	}, nil

}
