package main

import (
	"fmt"

	"github.com/avlitman/operator-observability/examples/rules"
	"github.com/avlitman/operator-observability/pkg/docs"
)

func main() {
	rules.SetupRules()
	docsString := docs.BuildAlertsDocs(rules.ListAlerts())
	fmt.Println(docsString)
}
