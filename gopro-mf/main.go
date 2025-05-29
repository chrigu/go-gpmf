package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"gopro/mp4"
	"gopro/telemetry"
)

func main() {
	app := &cli.App{
		Name:      "gopro-mf",
		Usage:     "Extract telemetry data from GoPro MP4 files",
		ArgsUsage: "in-mp4-file",
		Commands: []*cli.Command{
			{
				Name:    "print",
				Aliases: []string{"p"},
				Usage:   "Print telemetry data to the console",
				Action: func(cCtx *cli.Context) error {
					file, err := os.Open(cCtx.Args().First())
					if err != nil {
						log.Fatal("Error opening file:", err)
						return nil
					}
					defer file.Close()

					printTelemetryData(file)
					return nil
				},
			},
			{
				Name:    "extract-binary",
				Aliases: []string{"e"},
				Usage:   "Extract binary data from the MP4 file",
				Action: func(cCtx *cli.Context) error {
					inputFilename := cCtx.Args().First()
					file, err := os.Open(inputFilename)
					if err != nil {
						log.Fatal("Error opening file:", err)
						return nil
					}
					defer file.Close()

					extractBinaryData(file, inputFilename)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func printTelemetryData(file io.ReadSeeker) {
	// Extract metadata track from the MP4 file
	gpsData, gyroData, faceData, lumaData, huesData, sceneData := telemetry.ExtractTelemetryData(file, false)

	fmt.Printf("GPS Data: %v\n", gpsData)
	fmt.Printf("Gyro Data: %v\n", gyroData)
	fmt.Printf("Face Data: %v\n", faceData)
	fmt.Printf("Luma Data: %v\n", lumaData)
	fmt.Printf("Hues Data: %v\n", huesData)
	fmt.Printf("Scene Data: %v\n", sceneData)
}

func extractBinaryData(file io.ReadSeeker, inputFilename string) {
	gpmfRaw, _ := mp4.ExtractTelemetryFromMp4(file)
	if gpmfRaw == nil {
		log.Fatal("No telemetry data found in the file")
	}

	outputFilename := inputFilename[:len(inputFilename)-4] + ".bin"

	// Create output file
	outFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal("Error creating output file:", err)
	}
	defer outFile.Close()

	// Write the raw GPMF data to file
	_, err = outFile.Write(gpmfRaw)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}

	fmt.Printf("Successfully saved telemetry data to %s (size: %d bytes)\n", outputFilename, len(gpmfRaw))
}
