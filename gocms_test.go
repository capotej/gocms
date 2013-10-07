package main

import (
	"testing"
)

func TestReplace(t *testing.T) {
	result := ProcessTemplateWithInput("test.html", "layout.html")
	reference := `<html><head>
  <title></title>
</head>
<body>
  <div id="main">
  replaced
</div>
  <div id="footer">
    standard footer
  </div>


</body></html>`

	if result != reference {
		panic("result did not equal reference")
	}
}

func TestReplace2(t *testing.T) {
	result := ProcessTemplateWithInput("test2.html", "layout.html")
	reference := `<html><head>
  <title></title>
</head>
<body>
  <div id="main">
  replaced with:
  <h1> other stuff</h1>
</div>
  <div id="footer">
  replaced footer
</div>


</body></html>`

	if result != reference {
		panic("result did not equal reference")
	}

}
