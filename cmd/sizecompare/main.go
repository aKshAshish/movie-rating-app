package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"movie-rating-app/gen"
	model "movie-rating-app/metadata/pkg"

	"google.golang.org/protobuf/proto"
)

func serializeToJSON(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}

func serializeToXML(m *model.Metadata) ([]byte, error) {
	return xml.Marshal(m)
}

func serializeToProto(m *gen.Metadata) ([]byte, error) {
	return proto.Marshal(m)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "Batman",
	Description: "The Batman movie.",
	Director:    "Christoper Nolan",
}

var genMetadata = &gen.Metadata{
	Id:         "123",
	Title:      "Batman",
	Desciption: "The Batman movie.",
	Director:   "Christopher Nolan",
}

func main() {
	jsonBytes, err := serializeToJSON(metadata)
	checkError(err)

	xmlBytes, err := serializeToXML(metadata)
	checkError(err)

	protoBytes, err := serializeToProto(genMetadata)
	checkError(err)

	fmt.Printf("JSON size:\t%dB\n", len(jsonBytes))
	fmt.Printf("XML size:\t%dB\n", len(xmlBytes))
	fmt.Printf("Proto size:\t%dB\n", len(protoBytes))
}
