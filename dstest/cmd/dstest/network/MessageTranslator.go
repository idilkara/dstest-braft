package network

import (
	"log"
	"strings"
	"os"
	"bytes" 
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

type MessageTranslator interface {
	Translate(*Message) *Message
	TranslateBRPC(*Message) *Message
}

func NewMessageTranslator(messageType MessageType) MessageTranslator {
	switch messageType {
	case GRPC:
		return newGRPCTranslator()
	default:
		return nil
	}
}

type GRPCTranslator struct {
	Log *log.Logger
}

func newGRPCTranslator() *GRPCTranslator {
	translator := new(GRPCTranslator)
	translator.Log = log.New(os.Stdout, "[GRPCTranslator]", log.LstdFlags)
	return translator
}

func (t *GRPCTranslator) Translate(message *Message) *Message {
	f := message.Payload.(*http2.Framer)

	message.Type = GRPC

	decoder := hpack.NewDecoder(1024, nil)
	for message.Name == "" {
		frame, err := f.ReadFrame()
		if err != nil {
			t.Log.Printf("Error while reading frames: %s\n", err)
			break
		}

		switch ty := frame.(type) {
		case *http2.HeadersFrame:
			out, err := decoder.DecodeFull(ty.HeaderBlockFragment())
			if err != nil {
				t.Log.Printf("Error while decoding frames: %s\n", err)
			}
			for _, v := range out {
				if v.Name == ":path" {
					path := strings.Split(v.Value, "/")
					message.Name = path[len(path)-1]
				}
				
			}
		}
	}
	return message
}

//special translation function for baidu
func (t *GRPCTranslator) TranslateBRPC(message *Message) *Message {
	f, ok := message.Payload.(*bytes.Buffer)
	if !ok {
	    t.Log.Printf("Error: Payload is not of expected type, it is %T\n", message.Payload)
	    return message
	}

	payloadStr := f.String()

	    // BRPC has a pattern
	    startMarker := "RaftService/"
	    endMarker := "A"
	
	    // indices
	    startIndex := strings.Index(payloadStr, startMarker)
	    endIndex := strings.Index(payloadStr, endMarker)
	
	    // check if markers exist
	    if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
	        
	        extractedStr := payloadStr[startIndex+len(startMarker) : endIndex]
	        //t.Log.Printf("Extracted String: %s\n", extractedStr)
	        message.Name =  extractedStr

	    } else {

	        t.Log.Printf("Couldn't extract the message information")

	    }

	message.Type = GRPC
	
	return message
}