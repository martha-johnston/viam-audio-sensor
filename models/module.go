package models

import (
	"context"
	"encoding/binary"
	"fmt"

	"go.viam.com/rdk/components/audioin"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	rutils "go.viam.com/rdk/utils"
)

var AudioSensor = resource.NewModel("devin-hilly", "audio-sensor", "audio-sensor")

const (
	defaultThreshold     = 500.0
	defaultSampleSeconds = 1.0
)

func init() {
	resource.RegisterComponent(
		sensor.API,
		AudioSensor,
		resource.Registration[sensor.Sensor, *Config]{
			Constructor: newAudioSensor,
		},
	)
}

// Config holds the JSON attributes.
type Config struct {
	AudioIn       string  `json:"audio_in"`
	Threshold     float64 `json:"threshold,omitempty"`
	SampleSeconds float32 `json:"sample_seconds,omitempty"`
}

// Validate ensures audio_in is set and declares it as a required dependency.
func (c *Config) Validate(path string) ([]string, []string, error) {
	if c.AudioIn == "" {
		return nil, nil, fmt.Errorf("audio_in is required")
	}
	if c.Threshold < 0 {
		return nil, nil, fmt.Errorf("threshold must be non-negative")
	}
	if c.SampleSeconds < 0 {
		return nil, nil, fmt.Errorf("sample_seconds must be non-negative")
	}
	return []string{c.AudioIn}, nil, nil
}

type audioSensor struct {
	resource.AlwaysRebuild

	name          resource.Name
	logger        logging.Logger
	audioIn       audioin.AudioIn
	threshold     float64
	sampleSeconds float32
}

func newAudioSensor(
	ctx context.Context,
	deps resource.Dependencies,
	rawConf resource.Config,
	logger logging.Logger,
) (sensor.Sensor, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	aIn, err := audioin.FromProvider(deps, conf.AudioIn)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve audio_in dependency %q: %w", conf.AudioIn, err)
	}

	threshold := conf.Threshold
	if threshold == 0 {
		threshold = defaultThreshold
	}
	sampleSeconds := conf.SampleSeconds
	if sampleSeconds == 0 {
		sampleSeconds = defaultSampleSeconds
	}

	return &audioSensor{
		name:          rawConf.ResourceName(),
		logger:        logger,
		audioIn:       aIn,
		threshold:     threshold,
		sampleSeconds: sampleSeconds,
	}, nil
}

func (s *audioSensor) Name() resource.Name {
	return s.name
}

// Readings requests `sample_seconds` of PCM16 audio from the configured
// audio_in, tracks the peak absolute sample value over that window, and
// reports it alongside the stream playback position and the configured
// threshold.
func (s *audioSensor) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	chunkChan, err := s.audioIn.GetAudio(ctx, rutils.CodecPCM16, s.sampleSeconds, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("GetAudio failed: %w", err)
	}

	var peak int16
	var latestTimestampNs int64
	for chunk := range chunkChan {
		if chunk == nil {
			continue
		}
		if chunk.EndTimestampNanoseconds > latestTimestampNs {
			latestTimestampNs = chunk.EndTimestampNanoseconds
		}
		data := chunk.AudioData
		for i := 0; i+1 < len(data); i += 2 {
			sample := int16(binary.LittleEndian.Uint16(data[i : i+2]))
			if sample > peak {
				peak = sample
			} else if -sample > peak {
				peak = -sample
			}
		}
	}

	streamSeconds := float64(latestTimestampNs) / 1e9
	return map[string]interface{}{
		"peak":           int64(peak),
		"stream_seconds": streamSeconds,
		"threshold":      s.threshold,
		"sound_detected": float64(peak) > s.threshold,
	}, nil
}

func (s *audioSensor) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *audioSensor) Status(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (s *audioSensor) Close(ctx context.Context) error {
	return nil
}
