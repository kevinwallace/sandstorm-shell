package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
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
        <pre id="output">Type a shell command in the above textbox and press enter to execute.

Yes, DoS attacks are possible.  Please be nice while Sandstorm is still in alpha.</pre>
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
	// limit to 10 KiB of output
	out := &limitWriter{w, 10 * 1024, func() {
		fmt.Fprint(w, "\n...output truncated")
		cmd.Process.Kill()
	}}
	cmd.Stdout = out
	cmd.Stderr = out
	if err := runWithTimeout(cmd, 1*time.Second); err != nil {
		fmt.Fprintf(w, "\n%s", err)
	}
}

func runWithTimeout(cmd *exec.Cmd, d time.Duration) error {
	errch := make(chan error, 1)
	go func() {
		errch <- cmd.Run()
	}()
	select {
	case err := <-errch:
		return err
	case <-time.After(d):
		cmd.Process.Kill()
		return fmt.Errorf("timed out after %s", d)
	}
}

// limitWriter wraps an io.Writer, only allowing a limited number of bytes to be
// written before returning errors.
type limitWriter struct {
	w io.Writer // underlying writer
	n int       // number of bytes remaining
	f func()    // optional function to call when output is truncated
}

func (w *limitWriter) Write(p []byte) (n int, err error) {
	var trunc bool
	if len(p) > w.n {
		p = p[:w.n]
		trunc = true
	}
	if len(p) > 0 {
		n, err = w.w.Write(p)
		w.n -= n
		if w.n == 0 && w.f != nil {
			w.f()
		}
	}
	if trunc && err != nil {
		err = errors.New("output truncated")
	}
	return
}
