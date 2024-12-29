package dns

import (
	"log"

	dnstap "github.com/dnstap/golang-dnstap"

	"google.golang.org/protobuf/proto"
)

const outputChannelSize = 32

type callbackDNSTapOutput struct {
	callback func(*dnstap.Dnstap)
	outputs  []dnstap.Output
	data     chan []byte
	done     chan struct{}
}

func NewCallbackOutput(callback func(*dnstap.Dnstap)) dnstap.Output {
	return &callbackDNSTapOutput{
		callback: callback,
		data:     make(chan []byte, outputChannelSize),
		done:     make(chan struct{}),
	}
}

func (o *callbackDNSTapOutput) Add(output dnstap.Output) {
	o.outputs = append(o.outputs, output)
}

func (o *callbackDNSTapOutput) Close() {
	close(o.data)
	<-o.done
}

func (o *callbackDNSTapOutput) GetOutputChannel() chan []byte {
	return o.data
}

func (o *callbackDNSTapOutput) RunOutputLoop() {
	for payload := range o.data {
		// Mirror the payload to all outputs
		for _, output := range o.outputs {
			output.GetOutputChannel() <- payload
		}

		o.processFrame(payload)
	}

	for _, output := range o.outputs {
		output.Close()
	}
	close(o.done)
}

func (o *callbackDNSTapOutput) processFrame(frame []byte) {
	dt := &dnstap.Dnstap{}
	if err := proto.Unmarshal(frame, dt); err != nil {
		log.Printf("Error unmarshaling Dnstap message: %v", err)
		return
	}

	o.callback(dt)
}
