package utils

import (
	"fmt"
	"ledfx/api"
	"ledfx/audio"
	"ledfx/bridgeapi"
	log "ledfx/logger"
	"net/http"
	"regexp"
	"runtime"

	pretty "github.com/fatih/color"
	"github.com/rs/cors"
)

func ServeHttp() {
	DownloadFrontend()
	serveFrontend := http.FileServer(http.Dir("frontend"))
	fileMatcher := regexp.MustCompile(`\.[a-zA-Z]*$`)
	api.HandleApi()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !fileMatcher.MatchString(r.URL.Path) {
			http.ServeFile(w, r, "frontend/index.html")
		} else {
			serveFrontend.ServeHTTP(w, r)
		}
	})
}

func InitFrontend(ip string, port int) {
	pretty.Set(pretty.BgBlack, pretty.FgRed).Print("╭───────────────────────────────────────────────────────╮\n│               ")
	pretty.Set(pretty.BgBlack, pretty.FgRed, pretty.Bold).Print(" LedFx-Frontend")
	pretty.Set(pretty.BgBlack, pretty.FgWhite, pretty.Faint).Print(" by Blade ")
	pretty.Set(pretty.BgBlack, pretty.FgRed).Print("               │\n├───────────────────────────────────────────────────────┤\n│                                                       │\n│   ")
	switch runtime.GOOS {
	case "darwin":
		pretty.Set(pretty.BgBlack, pretty.FgHiYellow).Print("[CMD]+LMB: ")
	default:
		pretty.Set(pretty.BgBlack, pretty.FgHiYellow).Print("[CTRL]+Click: ")
	}
	pretty.Set(pretty.BgBlack, pretty.FgHiBlue, pretty.Bold, pretty.Underline).Print("http://localhost:8080/#/?newCore=1")
	switch runtime.GOOS {
	case "darwin":
		pretty.Set(pretty.BgBlack, pretty.FgRed).Print("       │\n")

	default:
		pretty.Set(pretty.BgBlack, pretty.FgRed).Print("    │\n")

	}
	pretty.Set(pretty.BgBlack, pretty.FgRed).Print("│                                                       │\n╰───────────────────────────────────────────────────────╯\n")
	pretty.Unset()

	go func() {
		mux := http.DefaultServeMux
		err := bridgeapi.NewServer(func(buf audio.Buffer) {
			// No callback for now
		}, mux)
		if err != nil {
			log.Logger.Fatal(err)
		}

		if err = http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), cors.AllowAll().Handler(http.DefaultServeMux)); err != nil {
			log.Logger.Fatal(err)
		}
	}()
}
