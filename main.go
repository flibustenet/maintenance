package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/hello/{x}", HelloServer)
	http.HandleFunc("/probe", Probe)
	http.HandleFunc("/", HelloServer)
	log.Println("Listen :" + port)
	http.ListenAndServe(":"+port, nil)
}

func Probe(w http.ResponseWriter, r *http.Request) {
	max, _ := strconv.Atoi(r.FormValue("max"))
	rd := 0
	if max > 0 {
		rd = rand.Intn(max)
	}
	s := fmt.Sprintf("%s sleep %ds", r.URL.RequestURI(), rd)
	fmt.Fprint(w, s)
	log.Println(s)
	time.Sleep(time.Duration(rd) * time.Second)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	redir := r.Header.Get("x-redirect")
	if redir != "" {
		fmt.Fprintf(w, `<!DOCTYPE html><html lang="fr"><body>`)
		fmt.Fprintf(w, "Nouvelle adresse : <a href='https://%s'>%s</a>", redir, redir)
		return
	}

	fmt.Fprintf(w, `
<!doctype html>
<title>Maintenance</title>
<style>
  body { text-align: center; padding: 150px; }
  h1 { font-size: 50px; }
  body { font: 20px Helvetica, sans-serif; color: #333; }
  article { display: block; text-align: left; width: 650px; margin: 0 auto; }
  a { color: #dc8100; text-decoration: none; }
  a:hover { color: #333; text-decoration: none; }
</style>

<article>
        <div>Travaux en cours, merci de revenir un peu plus tard %s...</div>
</article>
`, r.PathValue("x"))
}
