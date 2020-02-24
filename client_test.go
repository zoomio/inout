package inout

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fetch(t *testing.T) {
	defer stopServer(startServer(fmt.Sprintf(":%d", 8666), clientHTML, 200))
	res := fetch(context.TODO(), "http://localhost:8666")

	assert.NotNil(t, res)
	assert.Nil(t, res.err)

	assert.NotNil(t, res.req)
	assert.Equal(t, "http://localhost:8666", res.req.URL.String())

	assert.NotNil(t, res.resp)

	defer res.resp.Body.Close()
	bs, _ := ioutil.ReadAll(res.resp.Body)
	assert.Equal(t, clientHTML, string(bs))

	assert.Equal(t, 1, res.attempts)
}

func Test_fetch_retries(t *testing.T) {
	defer stopServer(startServer(fmt.Sprintf(":%d", 8666), "", 429))
	res := fetch(context.TODO(), "http://localhost:8666")
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
