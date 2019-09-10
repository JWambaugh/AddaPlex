package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/martamius/AddaPlex/pluginarch"
)

func deezerWorker(jobChan <-chan deezerJob) {
	log.Printf("Deezer: Job worker waiting for new jobs")

	for job := range jobChan {
		log.Printf("Deezer:Starting job for :" + job.url)
		download(job)
		log.Printf("Deezer: Job Completed")

	}
}

func download(job deezerJob) {

	line := "node "
	log.Print(line)

	os.Chdir(config.Options["SMLoadrPath"])
	cmd := exec.Command("node", path.Join(config.Options["SMLoadrPath"], "SMLoadr.js"), "-p", config.Options["outputPath"], "-q", config.Options["quality"], "-u", job.url)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("%s", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Print(err)
	}

	if err := cmd.Start(); err != nil {
		log.Print(err)
	}

	stde, _ := ioutil.ReadAll(stderr)
	stdo, _ := ioutil.ReadAll(stdout)
	log.Printf("%s", stdo)

	log.Printf("%s", stde)
}

type deezerPlugin struct {
}

type deezerJob struct {
	url string
}

var config pluginarch.PluginConfig
var jobChannel chan deezerJob

func (y deezerPlugin) Init(conf pluginarch.PluginConfig) {
	config = conf
	// log.Printf("%v", config)
	// log.Print("Deezer: Checking for deezer-dl")
	// cmd := exec.Command("deezer-dl", "--version")
	// err := cmd.Run()
	// if err != nil {
	// 	log.Print("Deezer: deezer-dl appears to be missing or not working. Please install it and make sure its in your PATH: https://github.com/ytdl-org/deezer-dl/")
	// 	log.Printf("Deezer: %v", err)
	// } else {
	// 	log.Printf("Deezer: deezer-dl looks ok!")
	// }

	// make a channel with a capacity of 100.
	jobChannel = make(chan deezerJob, 100)

	// start the worker
	go deezerWorker(jobChannel)
}

func (y deezerPlugin) Name() string {
	return "Deezer"
}

func (y deezerPlugin) Identifier() string {
	return "deezer"
}

func (y deezerPlugin) ActionDefinitions() []pluginarch.PluginAction {

	s := make([]pluginarch.PluginAction, 1)
	s[0] = pluginarch.PluginAction{
		Name: "Download Music",
		Type: "url",
		Options: map[string]string{
			"regex": "deezer.com/us/(artist|album|track|playlist)",
		},
	}

	return s
}

func (y deezerPlugin) PerformAction(action string, options map[string]string) (string, bool) {
	switch action {
	case "Download Music":
		job := deezerJob{
			url: options["url"],
		}
		log.Printf("Deezer: queueing job!")
		log.Printf("%v", y)

		//queue a new job
		jobChannel <- job
		return "Ok, Ill download that", true

	default:
		return "Unkown action", false
	}
}

// Plugin the plugin object to export
var Plugin deezerPlugin
