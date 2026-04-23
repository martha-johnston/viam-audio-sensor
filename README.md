# audio-sensor

A [Viam](https://viam.com/) sensor module that wraps an `audio_in` component and reports whether audio is being detected. Each `Readings` call samples a short window of PCM16 audio from the configured `audio_in`, computes the signal's RMS level, and returns `audio_detected` = true when the RMS exceeds the configured threshold.

Useful for: simple sound-activated automations, debugging audio pipelines, smoke tests. Composes with any `audio_in` component (e.g. [audio-replay](https://github.com/martha-johnston/viam-audio-replay), RTSP-audio modules, live microphones).

## Model: `devin-hilly:audio-sensor:audio-sensor`

### Configuration

```json
{
  "audio_in": "my-audio",
  "threshold": 500.0
}
```

### Attributes

| Name             | Type    | Inclusion | Default | Description                                                                           |
|------------------|---------|-----------|---------|---------------------------------------------------------------------------------------|
| `audio_in`       | string  | Required  | —       | Name of the `audio_in` component to sample.                                           |
| `threshold`      | float   | Optional  | `500.0` | RMS level above which `audio_detected` is reported true.                              |
| `sample_seconds` | float   | Optional  | `1.0`   | How much audio to sample on each `Readings` call.                                     |

### Readings

```json
{
  "audio_detected": true,
  "rms": 1523.7,
  "samples": 48000
}
```

### Example configuration

```json
{
  "modules": [
    {
      "type": "registry",
      "name": "devin-hilly_audio-sensor",
      "module_id": "devin-hilly:audio-sensor",
      "version": "0.0.1"
    }
  ],
  "components": [
    {
      "name": "my-audio",
      "api": "rdk:component:audio_in",
      "model": "devin-hilly:audio-replay:audio",
      "attributes": { "video_path": "/abs/path/video.mp4" }
    },
    {
      "name": "my-audio-sensor",
      "api": "rdk:component:sensor",
      "model": "devin-hilly:audio-sensor:audio-sensor",
      "attributes": {
        "audio_in": "my-audio",
        "threshold": 500.0
      }
    }
  ]
}
```

## Building from source

```bash
make build
make module.tar.gz
```
