package inout

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_waitForDomElement(t *testing.T) {
	defer stopServer(startServer(fmt.Sprintf(":%d", port), headlessIndexHTML, 200))
	innerContents, err := waitForDomElement(context.TODO(), "div p", fmt.Sprintf("http://localhost:%d", port), false)
	assert.Nil(t, err)
	assert.Equal(t, headlessExpectedHTML, innerContents)
}

// startServer is a simple HTTP server that displays the passed headers in the html.
func startServer(addr, document string, status int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(status)
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(res, document)
	})
	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	return srv
}

func stopServer(srv *http.Server) {
	// close the server gracefully
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

const (
	port = 8655

	headlessExpectedHTML = "<p>There was a Boy whose name was Jim;</p>" +
		"<p>His Friends were very good to him.</p>" +
		"<p>They gave him Tea, and Cakes, and Jam,</p>" +
		"<p>And slices of delicious Ham,</p>" +
		"<p>And Chocolate with pink inside,</p>" +
		"<p>And little Tricycles to ride,</p>" +
		"<p>Andread him Stories through and through,</p>" +
		"<p>And even took him to the Zoo—</p>" +
		"<p>But there it was the dreadful Fate</p>" +
		"<p>Befell him, which I now relate.</p>"

	headlessIndexHTML = `<!doctype html>
<html>
<head>
  <title>Test</title>
</head>
<body>
  <div id="box1" style="display:none">
    <div id="box2">
      <p class="line">There was a Boy whose name was Jim;</p>
	  <p class="line">His Friends were very good to him.</p>
	  <p class="line">They gave him Tea, and Cakes, and Jam,</p>
	  <p class="line">And slices of delicious Ham,</p>
	  <p class="line">And Chocolate with pink inside,</p>
	  <p class="line">And little Tricycles to ride,</p>
	  <p class="line">Andread him Stories through and through,</p>
	  <p class="line">And even took him to the Zoo—</p>
	  <p class="line">But there it was the dreadful Fate</p>
	  <p class="line">Befell him, which I now relate.</p>
    </div>
  </div>
  <script>
  	setTimeout(function() {
		document.querySelector('#box1').style.display = '';
	}, 3000);
  </script>
</body>
</html>`
)
