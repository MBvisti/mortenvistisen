package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/knights-analytics/hugot"
)

func main() {
	// decoded := base36.Decode("1eawesp")
	// decodedMinus := decoded - 1
	// encoded := base36.Encode(decodedMinus)
	//
	// slog.Info(
	// 	"this is results",
	// 	"decoded",
	// 	decoded,
	// 	"encoded",
	// 	encoded,
	// 	"decoded_minus",
	// 	decodedMinus,
	// 	"formatted",
	// 	fmt.Sprintf("t3_%v", encoded),
	// )
	session, err := hugot.NewSession(hugot.WithOnnxLibraryPath("./onnxruntime-linux-x64.so"))
	if err != nil {
		panic(err)
	}
	// dlOpts := hugot.NewDownloadOptions()
	// mPath, err := session.DownloadModel(
	// 	"sentence-transformers/all-MiniLM-L6-v2",
	// 	// "distilbert/distilbert-base-uncased-finetuned-sst-2-english",
	// 	// "KnightsAnalytics/distilbert-base-uncased-finetuned-sst-2-english",
	// 	// "sentence-transformers/all-mpnet-base-v2",
	// 	// "sentence-transformers/all-MiniLM-L12-v2",
	// 	"./",
	// 	dlOpts,
	// )
	// if err != nil {
	// 	panic(err)
	// }

	config := hugot.FeatureExtractionConfig{
		ModelPath: "./sentence-transformers_all-MiniLM-L6-v2/",
		Name:      "testPipelineClassi",
	}

	fPipeline, err := hugot.NewPipeline(session, config)
	if err != nil {
		log.Print("############################# HERE ############################")
		panic(err)
	}

	batchRes, err := fPipeline.RunPipeline(
		[]string{"suck a biiiiiiiiiiiiig cock", "suck a cock two"},
	)
	if err != nil {
		panic(err)
	}

	s, err := json.Marshal(batchRes)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(s))
}
