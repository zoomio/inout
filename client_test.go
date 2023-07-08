package inout

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fetch(t *testing.T) {
	defer stopServer(context.Background(), startServer(clientHTML, 200))
	res := fetch(source)

	assert.NotNil(t, res)

	assert.Nil(t, res.err)
	assert.Equal(t, 1, res.attempts)
	assert.NotNil(t, res.resp)

	defer res.resp.Body.Close()
	bs, _ := io.ReadAll(res.resp.Body)
	assert.Equal(t, clientHTML, string(bs))
}

func Test_fetch_retries(t *testing.T) {
	ts := startServer("", 429)
	defer ts.Close()
	res := fetch(source)
	assert.NotNil(t, res)
	assert.NotNil(t, res.err)
	assert.Equal(t, 3, res.attempts)
}

const (
	clientHTML = `<!doctype html>
<html>
<head>
  <title>Test</title>
</head>
<body>
  <div id="box1">Test</div>
</body>
</html>`
)
