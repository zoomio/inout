package inout

import (
	"bufio"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handleHTTP(t *testing.T) {
	ts := startServer(readerIndexHTML, 200)
	defer ts.Close()
	reader, _, err := handleHTTP(context.Background(), &config{source: source})
	assert.Nil(t, err)
	assert.NotNil(t, reader)

	r := &Reader{
		reader:  reader,
		scanner: bufio.NewScanner(reader),
	}
	lines, err := r.ReadLines()
	assert.Nil(t, err)
	assert.Equal(t, readerIndexHTML, strings.Join(lines, "\n"))
}

func Test_handleFS(t *testing.T) {
	reader, err := handleFS("file_reader_test.txt")
	assert.Nil(t, err)
	assert.NotNil(t, reader)

	r := &Reader{
		reader:  reader,
		scanner: bufio.NewScanner(reader),
	}
	lines, err := r.ReadLines()
	assert.Nil(t, err)
	assert.Equal(t, fileReaderText, strings.Join(lines, "\n"))
}

const (
	readerIndexHTML = `<!doctype html>
<html>
<head>
  <title>Test</title>
</head>
<body>
  <div id="box1">
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
</body>
</html>`

	fileReaderText = `There was a Boy whose name was Jim;
His Friends were very good to him.
They gave him Tea, and Cakes, and Jam,
And slices of delicious Ham,
And Chocolate with pink inside,
And little Tricycles to ride,
Andread him Stories through and through,
And even took him to the Zoo—
But there it was the dreadful Fate
Befell him, which I now relate.`
)
