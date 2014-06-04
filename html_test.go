package ronn

import (
  "testing"
  "io/ioutil"
  "fmt"
)

func TestHTML(t *testing.T) {
  files := []string{
    "custom_title_document",
    "basic_document",
  }

  for _, file := range files {
    testHTML(t, file)
  }
}

func testHTML(t *testing.T, file string) {
  ronn := example(file+".ronn")
  html := example(file+".html")

  by, err := ioutil.ReadFile(html)
  if err != nil {
    t.Fatalf("Unable to read %s: %s", html, err)
  }
  expected := string(by)

  input, err := ioutil.ReadFile(ronn)
  if err != nil {
    t.Fatalf("Unable to read %s: %s", ronn, err)
  }

  actual := HTML(input)
  if actual != expected {
    t.Errorf("%s output does not match", file)
    fmt.Println(file)
    fmt.Println("EXPECTED:")
    fmt.Println(expected)
    fmt.Println("ACTUAL:")
    fmt.Println(actual)
  }
}
