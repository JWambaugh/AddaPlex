package main

import (
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/martamius/AddaPlex/pluginarch"
)

func youTubeWorker(jobChan <-chan youTubeJob) {
	log.Printf("YouTube: Job worker waiting for new jobs")

	for job := range jobChan {
		log.Printf("YouTube:Starting job for :" + job.url)
		download(job)
		log.Printf("YouTube: Job Completed")

	}
}

func download(job youTubeJob) {
	log.Print(config.Options["outputDir"])
	out := config.Options["videoOutDir"]
	audio := ""
	opts := config.Options["videoOpts"]

	if job.audioOnly {
		out = config.Options["audioOutDir"]
		audio = "-x "
		opts = config.Options["audioOpts"]
	}

	line := "youtube-dl " + audio + " " + opts + " -o \"" + out + "\" " + job.url
	//log.Print(line)
	cmd := exec.Command("bash", "-lic", line)
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

type youTubePlugin struct {
}

type youTubeJob struct {
	url       string
	audioOnly bool
}

var config pluginarch.PluginConfig
var jobChannel chan youTubeJob

func (y youTubePlugin) Init(conf pluginarch.PluginConfig) {
	config = conf
	//log.Printf("%v", config)
	log.Print("YouTube: Checking for youtube-dl")
	cmd := exec.Command("youtube-dl", "--version")
	err := cmd.Run()
	if err != nil {
		log.Print("YouTube: youtube-dl appears to be missing or not working. Please install it and make sure its in your PATH: https://github.com/ytdl-org/youtube-dl/")
		log.Printf("YouTube: %v", err)
	} else {
		log.Printf("YouTube: youtube-dl looks ok!")
	}

	// make a channel with a capacity of 100.
	jobChannel = make(chan youTubeJob, 100)

	// start the worker
	go youTubeWorker(jobChannel)
}

func (y youTubePlugin) Name() string {
	return "Youtube"
}

func (y youTubePlugin) Identifier() string {
	return "youtube"
}

func (y youTubePlugin) ActionDefinitions() []pluginarch.PluginAction {

	s := make([]pluginarch.PluginAction, 2)
	s[0] = pluginarch.PluginAction{
		Name: "Download Video",
		Type: "url",
		Options: map[string]string{
			"regex": "youtube.com/watch",
		},
	}
	s[1] = pluginarch.PluginAction{
		Name: "Download Audio",
		Type: "url",
		Options: map[string]string{
			"regex": "youtube.com/watch",
		},
	}
	return s
}

func (y youTubePlugin) PerformAction(action string, options map[string]string) (string, bool) {
	switch action {
	case "Download Video":
		job := youTubeJob{
			url:       options["url"],
			audioOnly: false,
		}
		log.Printf("YouTube: queueing job!")
		log.Printf("%v", y)

		//queue a new job
		jobChannel <- job
		return "Ok, Ill download that video", true
	case "Download Audio":
		job := youTubeJob{
			url:       options["url"],
			audioOnly: true,
		}
		//queue a new job
		jobChannel <- job
		return "Ok, Ill download that audio", true
	default:
		return "Unkown action", false
	}
}

// Plugin the plugin object to export
var Plugin youTubePlugin
