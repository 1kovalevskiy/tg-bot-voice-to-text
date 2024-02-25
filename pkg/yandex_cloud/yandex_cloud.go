package yandex_cloud

import (
	"context"
	"fmt"
	"time"

	"log/slog"
	"net"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/ai/stt/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type gRPCClient struct {
	connection *grpc.ClientConn
	client     stt.RecognizerClient
	context    context.Context
	notify     chan error
}

func NewgRPCClient(token string) *gRPCClient {
	conn, err := grpc.DialContext(context.Background(), net.JoinHostPort("stt.api.cloud.yandex.net", "443"),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	md := metadata.Pairs("authorization", fmt.Sprintf("Api-Key %s", token))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	client := stt.NewRecognizerClient(conn)

	return &gRPCClient{
		connection: conn,
		client:     client,
		context:    ctx,
		notify:     make(chan error, 1),
	}

}

func (g *gRPCClient) NewStreamer() (stt.Recognizer_RecognizeStreamingClient, context.CancelFunc) {
	context, cancel := context.WithTimeout(g.context, time.Minute*60)
	stream, err := g.client.RecognizeStreaming(context)
	if err != nil {
		slog.Error(err.Error())
		cancel()
		return nil, nil
	}
	stream.Send(GenSettings())
	return stream, cancel
}

// func (g *gRPCClient) Notify() <-chan error {
// 	return g.notify
// }

func (g *gRPCClient) Shutdown() {
	g.connection.Close()
}

func GenSettings() *stt.StreamingRequest {
	options := &stt.StreamingOptions{
		RecognitionModel: &stt.RecognitionModelOptions{
			AudioFormat: &stt.AudioFormatOptions{
				AudioFormat: &stt.AudioFormatOptions_ContainerAudio{
					ContainerAudio: &stt.ContainerAudio{
						ContainerAudioType: stt.ContainerAudio_OGG_OPUS,
					},
				},
			},
			TextNormalization: &stt.TextNormalizationOptions{
				TextNormalization: stt.TextNormalizationOptions_TEXT_NORMALIZATION_ENABLED,
				ProfanityFilter:   false,
				LiteratureText:    true,
			},
			LanguageRestriction: &stt.LanguageRestrictionOptions{
				RestrictionType: stt.LanguageRestrictionOptions_WHITELIST,
				LanguageCode:    []string{"ru-RU"},
			},
			AudioProcessingType: stt.RecognitionModelOptions_REAL_TIME,
		},
	}
	return &stt.StreamingRequest{Event: &stt.StreamingRequest_SessionOptions{SessionOptions: options}}
}
