package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"sort"

	"github.com/martamius/AddaPlex/pluginarch"
	"github.com/martamius/AddaPlex/pluginmanager"
	"github.com/rs/cors"
)

type statusResponsePlugin struct {
	Name       string
	Identifier string
	Actions    []pluginarch.PluginAction
}
type statusResponse struct {
	Plugins    []statusResponsePlugin
	ServerName string
	Status     string
	Message    string
}

var usedSignatures []string

func getStatus(w http.ResponseWriter, r *http.Request) {
	if !checkSignature(w, r) {
		return
	}
	resp := statusResponse{}

	plugins := pluginmanager.Plugins()

	for _, plugin := range *plugins {
		p := statusResponsePlugin{}
		p.Name = plugin.Name()
		p.Identifier = plugin.Identifier()
		p.Actions = plugin.ActionDefinitions()

		resp.Plugins = append(resp.Plugins, p)
	}
	resp.ServerName = configData.ServerName
	resp.Status = "ok"
	js, _ := json.Marshal(resp)
	fmt.Fprintf(w, "%s", js)
}

type basicResponse struct {
	Status  string
	Message string
}

func performAction(w http.ResponseWriter, r *http.Request) {
	if !checkSignature(w, r) {
		return
	}

	plugins := pluginmanager.Plugins()
	for _, plugin := range *plugins {
		if plugin.Identifier() == r.FormValue("pluginIdentifier") {
			response := basicResponse{}
			var ok bool
			options := make(map[string]string)
			options["url"] = r.FormValue("url")
			response.Message, ok = plugin.PerformAction(r.FormValue("Name"), options)
			if ok {
				response.Status = "ok"
			} else {
				response.Status = "error"
			}
			js, _ := json.Marshal(response)
			fmt.Fprintf(w, "%s", js)
			return
		}
	}
	response := basicResponse{}
	response.Message = "Invalid plugin identifier"
	response.Status = "error"
	js, _ := json.Marshal(response)
	fmt.Fprintf(w, "%s", js)
}

func getMac(message []byte) string {
	key := configData.SharedKey
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
}

func checkSignature(w http.ResponseWriter, r *http.Request) bool {
	providedSig := r.Header.Get("X-Signature")
	sigValid := true

	//check if signature has been used
	for _, v := range usedSignatures {
		if v == providedSig {
			sigValid = false
			log.Printf("Signatue '%s' has already been used. Access denied", providedSig)
		}
	}

	if sigValid {
		// Get form keys in a list and sort them
		var keys []string

		r.ParseForm()
		// log.Printf("%v", r.Form)
		for k := range r.Form {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		key := configData.SharedKey
		mac := hmac.New(sha256.New, []byte(key))

		// loop through sorted keys and add keys/values to hash
		for _, k := range keys {
			val := r.Form[k][0]
			mac.Write([]byte(k))
			// log.Print(k)
			if val != "" {
				// log.Print(val)
				mac.Write([]byte(val))
			}
		}

		calculatedHash := hex.EncodeToString(mac.Sum(nil))

		match := providedSig == calculatedHash
		if !match {
			log.Printf("Signature mismatch:\nProvided:   %s\nCalculated: %s", providedSig, calculatedHash)
		}

		sigValid = sigValid && match
	}

	if !sigValid {
		response := basicResponse{}
		response.Message = "Forbidden"
		response.Status = "error"
		js, _ := json.Marshal(response)
		w.WriteHeader(401)
		fmt.Fprintf(w, "%s", js)

	} else {
		//add signature to list of used signatures
		usedSignatures = append(usedSignatures, providedSig)
	}

	return sigValid
}

func startHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	mux.HandleFunc("/status", getStatus)
	mux.HandleFunc("/action", performAction)

	log.Print("Listening on port " + configData.ListenPort + " as " + configData.ServerName)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"chrome-extension://dgnlchffpfppiohcconpfoogcanhiafj"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-Signature"},
		// Enable Debugging for testing, consider disabling in production
		//Debug: true,
	})
	handler := c.Handler(mux)
	log.Fatal(http.ListenAndServe(":"+configData.ListenPort, handler))
}
