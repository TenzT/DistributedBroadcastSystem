package server

import (
	"BroadcastService/cluster"
	"BroadcastService/cluster/echo"
	"BroadcastService/common"
	"BroadcastService/eventbus"
	"BroadcastService/eventbus/native"
	"BroadcastService/server/dbshttp"
	"BroadcastService/storage"
	"BroadcastService/storage/gocache"
	"BroadcastService/validator"
	"BroadcastService/validator/md5_validator"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	storage    storage.Storage
	validator  validator.Validator
	core       *Core
	cancelFunc context.CancelFunc
	ctx        context.Context
	cluster    cluster.Cluster
	eventBus   eventbus.EventBus
}

func (s *Server) Run() {
	log.Print("Running server.")

	// start http server
	go func() {
		log.Printf("Start http server on address: {%s}", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()
		if err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				log.Println("Shutting down http server gracefully.")
			} else {
				log.Panic(err.Error())
			}
		}
	}()

	// run main loop
	go s.loop()

	// run core
	go s.core.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Wating for signals
	<-sigs
	log.Println("notify sigs")
	s.httpServer.Shutdown(context.Background())
	s.cancelFunc()
	log.Println("server shutdown")

	time.Sleep(1 * time.Second)
}

func (s *Server) loop() {
	log.Println("server loop")
	for {
		select {
		case event := <-s.eventBus.Subscribe():
			// TODO: implements event loop
			log.Println(event)
			switch event.EventType {
			case eventbus.EVENT_TYPE_BROADCAST_DATA_BATCH:
				dataList, ok := event.Payload.([]common.Data)
				if !ok {
					log.Println("Error parsing data to be broadcast")
				} else {
					s.cluster.Broadcast(dataList)
				}
			case eventbus.EVENT_TYPE_RECEIVE_DATA:
				data, ok := event.Payload.(common.Data)
				if !ok {
					log.Println("Error parsing data from cluster")
				} else {
					s.core.Deliver(&data)
				}
			}
		case data := <-s.cluster.Receive():
			event := eventbus.Event{
				EventType: eventbus.EVENT_TYPE_RECEIVE_DATA,
				Payload:   data,
			}
			s.eventBus.Publish(event)
		case event := <-s.core.Events():
			s.eventBus.Publish(event)
		case <-s.ctx.Done():
			log.Println("stop server")
			return
		}
	}
}

func (s *Server) HandleNewPosts(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte("Error reading data"))
		return
	}

	time.Now().String()
	data := &common.Data{}
	err = json.Unmarshal(b, data)
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte("Error parsing data"))
		return
	}
	data.TimeStamp = strconv.Itoa(int(time.Now().Unix()))

	err = s.core.Deliver(data)

	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte(err.Error()))
		return
	}

	result := &dbshttp.PostNewDataRsp{}

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

	list, err := s.core.GetAllData()
	if err != nil {
		w.WriteHeader(200)
		w.Write([]byte("Errror getting all data"))
		return
	}

	resList := make([]*dbshttp.HttpData, 0)
	for _, v := range list {
		resList = append(resList, &dbshttp.HttpData{
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

func New(config *ServerConfig) *Server {
	s := &Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/postNewData", s.HandleNewPosts)
	mux.HandleFunc("/listAllRecentData", s.HandleListAllRecentData)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HttpPort),
		Handler: mux,
	}
	s.httpServer = httpServer

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	s.ctx = ctx

	kvstore := gocache.New()
	validator := md5_validator.New()
	cluster := echo.NewEchoCluster()
	eventbus := native.NewNativeChannelEventBus()

	s.validator = validator
	s.storage = kvstore
	s.cluster = cluster
	s.eventBus = eventbus

	// create Core
	core := NewCore(kvstore, ctx, validator)

	s.core = core

	return s
}
