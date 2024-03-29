package superstream

import (
	"fmt"

	"github.com/IBM/sarama"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func startInterceptors(config interface{}, client *Client) {
	if config == nil {
		client.handleError(fmt.Sprintf(" startInterceptors: config is nil"))
		return
	}

	if config, ok := config.(*sarama.Config); ok {
		ConfigSaramaInterceptor(config, client)
		return
	} else {
		client.handleError(fmt.Sprintf(" startInterceptors: unsupported sdk"))
		fmt.Println("superstream: unsupported sdk")
		return
	}
}

func protoToJson(msgBytes []byte, desc protoreflect.MessageDescriptor) ([]byte, error) {
	newMsg := dynamicpb.NewMessage(desc)
	err := proto.Unmarshal(msgBytes, newMsg)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := protojson.Marshal(newMsg)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

func jsonToProto(msgBytes []byte, desc protoreflect.MessageDescriptor) ([]byte, error) {
	newMsg := dynamicpb.NewMessage(desc)
	err := protojson.Unmarshal(msgBytes, newMsg)
	if err != nil {
		return nil, err
	}

	protoBytes, err := proto.Marshal(newMsg)
	if err != nil {
		return nil, err
	}

	return protoBytes, nil
}
