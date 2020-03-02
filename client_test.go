package inout

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fetch(t *testing.T) {
	ts := startServer(clientHTML, 200)
	defer ts.Close()
	res := fetch(ts.URL)

	assert.NotNil(t, res)

	assert.Nil(t, res.err)
	assert.Equal(t, 1, res.attempts)
	assert.NotNil(t, res.resp)

	defer res.resp.Body.Close()
	bs, _ := ioutil.ReadAll(res.resp.Body)
	assert.Equal(t, clientHTML, string(bs))
}

func Test_fetch_retries(t *testing.T) {
	ts := startServer("", 429)
	defer ts.Close()
	res := fetch(ts.URL)
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
