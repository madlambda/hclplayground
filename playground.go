package main

import (
	"log"
	"strings"
	"syscall/js"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func main() {
	log.Println("starting playground")

	document := js.Global().Get("document")
	codeInput := document.Call("getElementById", "code")
	formattedOutput := document.Call("getElementById", "formatted")
	errorsOutput := document.Call("getElementById", "errors")
	validateButton := document.Call("getElementById", "validate")

	buttonClicks := make(chan struct{})
	validateButton.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		log.Println("button clicked")
		buttonClicks <- struct{}{}
		return nil
	}))

	for range buttonClicks {
		log.Println("validating")

		hclInput := []byte(codeInput.Get("value").String())
		p := hclparse.NewParser()
		_, diags := p.ParseHCL(hclInput, "playground")

		log.Println("parsed code, checking errors")

		errMsgs := make([]string, len(diags))

		for i, diag := range diags {
			errMsgs[i] = diag.Error()
		}

		errorsOutput.Set("value", strings.Join(errMsgs, "\n"))

		if diags.HasErrors() {
			continue
		}

		log.Println("no errors found, formatting code")

		formattedCode := hclwrite.Format(hclInput)
		formattedOutput.Set("value", string(formattedCode))
	}
}
