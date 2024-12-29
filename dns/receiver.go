package dns

import (
	"net"

	dnstap "github.com/dnstap/golang-dnstap"
)

type DNSTapReceiver struct {
	input   dnstap.Input
	output  dnstap.Output
	running bool
}

func NewDNSTapReceiver(listener net.Listener, callback func(*DNSResponse)) *DNSTapReceiver {
	input := dnstap.NewFrameStreamSockInput(listener)
	adaptedCallback := dnsTapCallbackAdapter(callback)
	output := NewCallbackOutput(adaptedCallback)

	return &DNSTapReceiver{
		input:   input,
		output:  output,
		running: false,
	}
}

func (r *DNSTapReceiver) Start() {
	if r.running {
		return
	}

	r.running = true
	go r.output.RunOutputLoop()
	go r.input.ReadInto(r.output.GetOutputChannel())

	r.input.Wait()
	r.running = false
}
