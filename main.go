package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("<html><body><form action='http://localhost:8080/pdf' method='post'><input type='submit' value='Generate PDF'></form></body></html>")))
}

func PDF(w http.ResponseWriter, r *http.Request) {

	cmd := exec.Command("php", "./dompdf/dompdf.php", "./dompdf/receipt.html")
	log.Println("&&&&&&&&&")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err)
		return
	}
	log.Println("**********")
	fmt.Print(string(out))

	// NOTE : In real world application, the output filename should be dynamic
	// and not static like in this example. For the sake of simplicity, we just
	// keep the output name to receipt.pdf
	// to change the output filename, see the domPDF command line instruction.

	// grab the generated receipt.pdf file and stream it to browser
	streamPDFbytes, err := ioutil.ReadFile("./dompdf/receipt.pdf")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b := bytes.NewBuffer(streamPDFbytes)

	// stream straight to client(browser)
	w.Header().Set("Content-type", "application/pdf")

	if _, err := b.WriteTo(w); err != nil {
		fmt.Fprintf(w, "%s", err)
	}

	w.Write([]byte("PDF Generated"))
}

func main() {
	// http.Handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/pdf", PDF)

	http.ListenAndServe(":8080", mux)
}
