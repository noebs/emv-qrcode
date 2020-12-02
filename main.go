package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/noebs/emv-qrcode/emv/mpm"
)

var tmpl = template.Must(template.ParseFiles("static/qr.html"))

func main() {
	var qr QRCheck

	http.HandleFunc("/encode", qr.DecodeHandler)
	http.HandleFunc("/decode", qr.DecodeHandler)
	http.HandleFunc("/", qr.DecodeHtml)
	http.ListenAndServe(":8012", nil)

}

func (q QRCheck) DecodeHandler(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &q); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if q.Qr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qrencoded, err := mpm.Decode(q.Qr)
	if err != nil {
		log.Printf("Error in parsing QR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qr := QR{}
	if err := qr.init(qrencoded); err != nil {
		log.Printf("Error in parsing QR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(qr)
	if err != nil {
		log.Printf("Error in parsing QR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func (q QRCheck) DecodeHtml(w http.ResponseWriter, r *http.Request) {

	data := r.URL.Query().Get("qr")
	q.Qr = data
	if q.Qr == "" {
		log.Printf("Empty qr")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qrencoded, err := mpm.Decode(q.Qr)
	if err != nil {
		log.Printf("Error in parsing QR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	qr := QR{}
	if err := qr.init(qrencoded); err != nil {
		log.Printf("Error in parsing QR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tmpl.Execute(w, qr)

}

//QRCheck Returns a payment token
type QRCheck struct {
	Qr        string  `json:"qr,omitempty"`
	Latitude  float32 `json:"latitude,omitempty"`
	Longitude float32 `json:"longitude,omitempty"`
}

type validationErrors struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}
