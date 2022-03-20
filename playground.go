package main

import (
	"fmt"
	"log"
	"strings"
	"syscall/js"

	"github.com/hashicorp/hcl/v2/hclparse"
)

func main() {
	log.Println("starting playground")

	document := js.Global().Get("document")
	codeInput := document.Call("getElementById", "code")
	formattedOutput := document.Call("getElementById", "formatted")
	errorsOutput := document.Call("getElementById", "errors")
	validateButton := document.Call("getElementById", "validate")

	fmt.Println(codeInput.Get("value"))
	fmt.Println(formattedOutput.Get("value"))
	formattedOutput.Set("value", "hallo")

	buttonClicks := make(chan struct{})
	validateButton.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		log.Println("button clicked")
		buttonClicks <- struct{}{}
		return nil
	}))

	for range buttonClicks {
		log.Println("received button clicked msg")

		hclInput := []byte(codeInput.Get("value").String())
		p := hclparse.NewParser()
		_, diags := p.ParseHCL(hclInput, "playground")

		errMsgs := make([]string, len(diags))

		for i, diag := range diags {
			errMsgs[i] = diag.Error()
		}

		errorsOutput.Set("value", strings.Join(errMsgs, "\n"))
	}
}
