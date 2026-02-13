package utils

import (
	"bytes"
	"context"
	"fmt"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

const (
	VideoMaxWidth       = 720
	VideoCRF            = 23
	VideoPreset         = "veryfast"
	VideoAudioBitrate   = "128k"
	VideoThumbAtSeconds = 1
)

type VideoProbeOptions struct {
	Timeout time.Duration
	Kwargs   ffmpeg.KwArgs
}

func ProbeVideo(path string, opts VideoProbeOptions) (string, error) {
	kwargs := ffmpeg.KwArgs{}
	for k, v := range opts.Kwargs {
		kwargs[k] = v
	}
	return ffmpeg.ProbeWithTimeoutExec(path, opts.Timeout, kwargs)
}

func TranscodeToMP4_720p(ctx context.Context, inPath, outPath string) error {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	s := ffmpeg.Input(inPath).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("min(%d\\,iw):-2", VideoMaxWidth)}).
		Output(outPath, ffmpeg.KwArgs{
			"c:v":      "libx264",
			"preset":   VideoPreset,
			"crf":      fmt.Sprintf("%d", VideoCRF),
			"pix_fmt":  "yuv420p",
			"movflags": "+faststart",
			"c:a":      "aac",
			"b:a":      VideoAudioBitrate,
		})

	s.Context = ctx

	if err := s.OverWriteOutput().WithOutput(&outBuf, &errBuf).Run(); err != nil {
		if errBuf.Len() > 0 {
			return fmt.Errorf("ffmpeg transcode error: %s: %w", errBuf.String(), err)
		}
		return err
	}
	return nil
}

func ExtractFrameJPEG(ctx context.Context, inPath string, atSeconds int) ([]byte, error) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	s := ffmpeg.Input(inPath, ffmpeg.KwArgs{
		"ss": fmt.Sprintf("%d", atSeconds),
	}).
		Output("pipe:", ffmpeg.KwArgs{
			"vframes": 1,
			"format":  "image2",
			"vcodec":  "mjpeg",
		})

	s.Context = ctx

	if err := s.OverWriteOutput().WithOutput(&outBuf, &errBuf).Run(); err != nil {
		if errBuf.Len() > 0 {
			return nil, fmt.Errorf("ffmpeg extract frame error: %s: %w", errBuf.String(), err)
		}
		return nil, err
	}
	return outBuf.Bytes(), nil
}
