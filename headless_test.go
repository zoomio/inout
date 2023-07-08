package inout

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	serverPort = 8655
	serverAddr = fmt.Sprintf(":%d", serverPort)
	source     = fmt.Sprintf("http://localhost:%d", serverPort)
)

// table driven tests
var headlessTests = []struct {
	name string

	httpStatus int
	input      *config

	expected *headlesResult
	err      error
}{
	{
		"headless_query",
		200,
		&config{source: source, query: "div p"},
		&headlesResult{htmlDoc: headlessExpectedHTML},
		nil,
	},
	{
		"headless_waitFor",
		200,
		&config{source: source, waitFor: "body"},
		&headlesResult{htmlDoc: headlessIndexHTML},
		nil,
	},
	{
		"headless_waitUntil",
		200,
		&config{source: source, waitUntil: 50 * time.Millisecond},
		&headlesResult{htmlDoc: headlessIndexHTML},
		nil,
	},
	{
		"headless_screenshot",
		200,
		&config{source: source, waitUntil: 50 * time.Millisecond, screenshot: true},
		&headlesResult{htmlDoc: headlessIndexHTML},
		nil,
	},
}

func Test_headless(t *testing.T) {
	ctx := context.Background()
	for _, tt := range headlessTests {
		t.Run(tt.name, func(t *testing.T) {
			defer stopServer(ctx, startServer(headlessIndexHTML, tt.httpStatus))
			res, err := headless(ctx, tt.input)

			if tt.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			if tt.expected != nil {
				assert.NotNil(t, res)
				assert.Equal(t, tt.expected.htmlDoc, res.htmlDoc)
				if len(tt.expected.imgBytes) > 0 {
					assert.Len(t, res.imgBytes, len(tt.expected.imgBytes))
				}
			}
		})
	}
}

func Benchmark_headless(b *testing.B) {
	ctx := context.Background()
	cfg := &config{source: source, waitUntil: 50 * time.Millisecond, screenshot: true}
	defer stopServer(ctx, startServer(headlessIndexHTML, 200))
	for i := 0; i < b.N; i++ {
		_, _ = headless(ctx, cfg)
	}
}

// ------------------------------- Setup -------------------------------

// startServer is a simple HTTP server that displays the passed headers in the html.
func startServer(pageHTML string, status int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, pageHTML)
	})
	srv := &http.Server{Addr: serverAddr, Handler: mux}
	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe(): %s\n", err)
		}
	}()
	return srv
}

func stopServer(c context.Context, srv *http.Server) {
	// close the server gracefully
	if err := srv.Shutdown(c); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
}

const (
	headlessExpectedHTML = "<p class=\"line\">There was a Boy whose name was Jim;</p>" +
		"<p class=\"line\">His Friends were very good to him.</p>" +
		"<p class=\"line\">They gave him Tea, and Cakes, and Jam,</p>" +
		"<p class=\"line\">And slices of delicious Ham,</p>" +
		"<p class=\"line\">And Chocolate with pink inside,</p>" +
		"<p class=\"line\">And little Tricycles to ride,</p>" +
		"<p class=\"line\">Andread him Stories through and through,</p>" +
		"<p class=\"line\">And even took him to the Zoo—</p>" +
		"<p class=\"line\">But there it was the dreadful Fate</p>" +
		"<p class=\"line\">Befell him, which I now relate.</p>"

	headlessIndexHTML = `<!DOCTYPE html><html><head>
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

</body></html>`
)
