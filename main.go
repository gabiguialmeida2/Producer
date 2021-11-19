package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gabiguialmeida2/rabbitConnection"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var ch *amqp.Channel

func apiResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		enfileirarMensagem(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}

func enfileirarMensagem(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Content-type 'application/json' esperado, mas foi recebido '%s'", ct)))
		return
	}

	var pessoa rabbitConnection.Pessoa
	err = json.Unmarshal(bodyBytes, &pessoa)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	rabbitConnection.PublishMessage(ch, pessoa)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Pessoa enfileirada com sucesso"}`))
}

func main() {

	conn, ch = rabbitConnection.CreateConnection()
	defer conn.Close()
	defer ch.Close()

	_ = rabbitConnection.DeclareQueue(conn, ch)

	http.HandleFunc("/pessoa", apiResponse)
	log.Fatal(http.ListenAndServe(":8089", nil))
}
