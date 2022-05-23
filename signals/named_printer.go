package signals

import (
	"log"

)

var _ Output = NamedPrinterDecorator{}

type NamedPrinterDecorator struct {
	Name     string
	Internal Output
}

// Set implements signals.IO
func (np NamedPrinterDecorator) Set(v bool) {
	if v {
		log.Printf("%s on", np.Name)
	} else{
		log.Printf("%s off", np.Name)
	}
	if np.Internal != nil{
		np.Internal.Set(v)
	}
}
