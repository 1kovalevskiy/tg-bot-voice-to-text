package yandexcloudstt

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/ai/stt/v3"
)

type yandexClient interface {
	NewStreamer() (stt.Recognizer_RecognizeStreamingClient, context.CancelFunc)
}

type yandexSSTClient struct {
	client yandexClient
}

func NewSSTClient(client yandexClient) *yandexSSTClient {
	return &yandexSSTClient{
		client: client,
	}
}

func (sst *yandexSSTClient) RecognizeAudio(fileReader io.Reader) (string, error) {
	stream, cancel := sst.client.NewStreamer()
	defer cancel()
	defer stream.CloseSend()

	waitResponse := make(chan error)
	resp := make(chan string)

	go func() {
		var result string = ""
		counter := 0
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				slog.Info("Recieve from yandex-cloud", "message_count", counter)
				waitResponse <- nil
				resp <- result
				return
			}
			counter += 1
			if err != nil {
				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
				resp <- result
				return
			}
			final_refinement := res.GetFinalRefinement()
			if final_refinement == nil {
				continue
			}
			normalized_text := final_refinement.GetNormalizedText()
			if final_refinement == nil {
				continue
			}
			alternatives := normalized_text.GetAlternatives()
			if alternatives == nil {
				continue
			}
			if len(alternatives) == 0 {
				continue
			}
			result = alternatives[0].GetText()
		}
	}()

	reader := bufio.NewReader(fileReader)
	buf := make([]byte, 4096)
	for {
		num, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error(err.Error())
			return "", err
		}

		data := &stt.StreamingRequest{Event: &stt.StreamingRequest_Chunk{Chunk: &stt.AudioChunk{Data: buf[:num]}}}
		if err := stream.Send(data); err != nil {
			slog.Error(err.Error())
			return "", err
		}
	}
	if err := stream.CloseSend(); err != nil {
		slog.Error(err.Error())
		return "", err
	}
	err := <-waitResponse
	result := <-resp
	return result, err
}
