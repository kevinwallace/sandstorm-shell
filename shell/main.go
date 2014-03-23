package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/cmd", handleCommand)

	panic(http.ListenAndServe("127.0.0.1:8080", nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	fmt.Fprint(w, `
    <html>
      <head>
        <style>
          #cmd {
            width: 100%;
          }
        </style>
      </head>
      <body>
        <form id="cmdform">
          <input id="cmd" />
        </form>
        <pre id="output"></pre>
        <script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.0/jquery.min.js"></script>
        <script>
          $("#cmdform").submit(function(e) {
            e.preventDefault();
            var cmd = $("#cmd").val();
            $("#cmd").val("");
            $("#output").text("");
            $.post("/cmd", {"cmd": cmd}, function(response) {
              $("#output").text("$ " + cmd + "\n" + response);
            });
          });
          $("#cmd").focus();
        </script>
      </body>
    </html>
  `)
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/plain")
	cmd := exec.Command("/bin/sh", "-c", r.FormValue("cmd"))
	output, err := cmd.CombinedOutput()
	w.Write(output)
	if err != nil {
		fmt.Fprintf(w, "\n%s", err)
	}
}
