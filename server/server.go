package server

import (
	"BroadcastService/common"
	"BroadcastService/server/dbshttp"
	"BroadcastService/storage"
	"BroadcastService/storage/gocache"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	storage    storage.Storage
}

func (s *Server) Run() {
	log.Print("Running server.")

	go func() {
		s.httpServer.ListenAndServe()
		log.Println("Shutting down http server gracefully.")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Wating for signals
	<-sigs
	fmt.Println("notify sigs")
	s.httpServer.Shutdown(context.Background())
	fmt.Println("server shutdown")

	time.Sleep(1 * time.Second)
}

func (s *Server) HandleNewPosts(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte("Error reading data"))
		return
	}

	data := &common.Data{}
	err = json.Unmarshal(b, data)
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte("Error parsing data"))
		return
	}
	data.TimeStamp = time.Now().String()

	result := &dbshttp.PostNewDataRsp{}
	err = s.storage.SaveData(data)
	if err != nil {
		log.Println(err.Error())
		result.Result = false
		result.Msg = err.Error()
	} else {
		result.Result = true
		result.Msg = "ok"
	}
	b, err = json.Marshal(result)
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte("Errror marshalling data"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(b)
	return

}

func (s *Server) HandleListAllRecentData(w http.ResponseWriter, r *http.Request) {

	list, err := s.storage.GetAllData()
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte("Errror getting all data"))
		return
	}

	resList := make([]*common.Data, 0)
	for _, v := range list {
		resList = append(resList, &common.Data{
			Id:        v.GetId(),
			Raw:       v.GetRawData(),
			Signature: v.GetSignature(),
		})
	}

	result := &dbshttp.ListAllRecentDataRsp{
		DataList: resList,
	}

	b, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte("Errror marshalling data"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write(b)
	return
}

func New() *Server {
	s := &Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/postNewData", s.HandleNewPosts)
	mux.HandleFunc("/listAllRecentData", s.HandleListAllRecentData)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	s.httpServer = httpServer

	s.storage = gocache.New()

	return s
}