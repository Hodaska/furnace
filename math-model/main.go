package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/server"
	"github.com/gopcua/opcua/ua"
)

const (
	maxTemp     = 300.0
	ambientTemp = 20.0
	heatingTau  = 120.0
	coolingTau  = 30.0
)

func main() {
	log.SetFlags(log.LstdFlags)

	hostname, _ := os.Hostname()

	opts := []server.Option{
		server.EnableSecurity("None", ua.MessageSecurityModeNone),
		server.EnableAuthMode(ua.UserTokenTypeAnonymous),
		server.EndPoint("0.0.0.0", 4840),
		server.EndPoint("localhost", 4840),
		server.EndPoint("math-model", 4840),
		server.EndPoint(hostname, 4840),
	}

	s := server.New(opts...)

	ns := server.NewMapNamespace(s, "Furnace")
	ns.Data["Power"] = float64(0)
	ns.Data["Temperature"] = float64(ambientTemp)

	root_ns, _ := s.Namespace(0)
	root_ns.Objects().AddRef(ns.Objects(), id.HasComponent, true)

	if err := s.Start(context.Background()); err != nil {
		log.Fatalf("Error starting OPC UA server: %s", err)
	}
	defer s.Close()
	log.Println("OPC UA server started at opc.tcp://0.0.0.0:4840")

	temperature := ambientTemp

	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			ns.Mu.Lock()
			power := ns.Data["Power"].(float64)
			ns.Mu.Unlock()

			target := ambientTemp
			if power > 0 {
				target = ambientTemp + (maxTemp-ambientTemp)*(power/100.0)
			}

			tau := coolingTau
			if temperature < target {
				tau = heatingTau
			}

			temperature += (target - temperature) * (1 - math.Exp(-1.0/tau))

			ns.SetValue("Temperature", temperature)
			log.Printf("Power=%.1f%% Temp=%.2f°C", power, temperature)
		}
	}()

	// POST /power?value=100  — Node-RED sets power (replaces OPC UA write)
	http.HandleFunc("/power", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		val, err := strconv.ParseFloat(r.URL.Query().Get("value"), 64)
		if err != nil || val < 0 || val > 100 {
			http.Error(w, "invalid value", http.StatusBadRequest)
			return
		}
		ns.Mu.Lock()
		ns.Data["Power"] = val
		ns.Mu.Unlock()
		log.Printf("Power set to %.1f%% via HTTP", val)
		fmt.Fprintln(w, "OK")
	})

	// GET /temperature  — Node-RED reads temperature (replaces OPC UA read)
	http.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) {
		ns.Mu.Lock()
		temp := ns.Data["Temperature"].(float64)
		ns.Mu.Unlock()
		fmt.Fprintf(w, "%.4f", temp)
	})

	go func() {
		log.Println("HTTP server started at :8888")
		if err := http.ListenAndServe(":8888", nil); err != nil {
			log.Fatalf("HTTP server error: %s", err)
		}
	}()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
	log.Println("Shutting down...")
}
