package telemetry

import (
	"gopro/mp4"
	"gopro/parser"
)

type TimedGyro struct {
	parser.Gyroscope
	TimeSample
}

// todo: refactor
func AddTimestampsToGyroDataWithDownsample(
	gyroData [][]parser.Gyroscope,
	telemetryMetadata *mp4.TelemetryMetadata,
	downsampleIntervalMs uint32,
) []TimedGyro {
	var TimedGyros []TimedGyro
	var sampleIndex uint32 = 0
	var sampleScaleTime uint32 = 0

	var accumulatedGyro parser.Gyroscope
	var count uint32 = 0
	var lastSampleScaleTime int64 = 0

	downsampleScaleThreshold := int64(telemetryMetadata.TimeScale * downsampleIntervalMs / 1000)

	for _, timeToSample := range telemetryMetadata.TimeToSamples {
		for range int(timeToSample.SampleCount) {
			if sampleIndex >= uint32(len(gyroData)) {
				break
			}

			currentTimedGyros := gyroData[sampleIndex]
			sampleCount := uint32(len(currentTimedGyros))

			for _, gyro := range currentTimedGyros {
				// Accumulate gyro values
				accumulatedGyro.X += gyro.X
				accumulatedGyro.Y += gyro.Y
				accumulatedGyro.Z += gyro.Z
				count++

				// Calculate individual timestamp for this sample
				individualTime := calculateIndividualTime(telemetryMetadata.CreationTime, int64(sampleScaleTime), telemetryMetadata.TimeScale)

				// Check if enough time has passed to downsample
				if int64(sampleScaleTime)-lastSampleScaleTime >= downsampleScaleThreshold {
					avgGyro := averageGyro(accumulatedGyro, count)

					TimedGyros = append(TimedGyros, TimedGyro{
						Gyroscope: avgGyro,
						TimeSample: TimeSample{
							TimeStamp: individualTime,
						},
					})

					// Reset accumulators
					accumulatedGyro = parser.Gyroscope{}
					lastSampleScaleTime = int64(sampleScaleTime)
					count = 0
				}

				// Increment scaled time based on sample delta
				sampleScaleTime += timeToSample.SampleDelta / sampleCount
			}
			sampleIndex++
		}
	}

	return TimedGyros
}

// Helper: Compute average Gyroscope reading
func averageGyro(accumulated parser.Gyroscope, count uint32) parser.Gyroscope {
	return parser.Gyroscope{
		X: accumulated.X / float32(count),
		Y: accumulated.Y / float32(count),
		Z: accumulated.Z / float32(count),
	}
}

// Helper: Compute individual timestamp for a sample
func calculateIndividualTime(creationTime int64, sampleScaleTime int64, timeScale uint32) int64 {
	return int64(float64(sampleScaleTime)/float64(timeScale)*1000) + creationTime
}
