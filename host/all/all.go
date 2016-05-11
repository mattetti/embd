// Package all conviniently loads all the inbuilt/supported host drivers.
package all

import (
	_ "github.com/mattetti/embd/host/bbb"
	_ "github.com/mattetti/embd/host/rpi"
)
