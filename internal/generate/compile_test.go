package generate

import (
	"os"
	"reflect"
	"testing"
	"time"

	"go.einride.tech/can"
	"go.einride.tech/can/pkg/descriptor"
	examplecan "go.einride.tech/can/testdata/gen/go/example"
	"gotest.tools/v3/assert"
)

func TestCompile_ExampleDBC(t *testing.T) {
	finish := runTestInDir(t, "../..")
	defer finish()
	const exampleDBCFile = "testdata/dbc/example/example.dbc"
	exampleDatabase := &descriptor.Database{
		SourceFile: exampleDBCFile,
		Version:    "",
		Nodes: []*descriptor.Node{
			{
				Name: "DBG",
			},
			{
				Name:        "DRIVER",
				Description: "The driver controller driving the car",
			},
			{
				Name: "IO",
			},
			{
				Name:        "MOTOR",
				Description: "The motor controller of the car",
			},
			{
				Name:        "SENSOR",
				Description: "The sensor controller of the car",
			},
		},

		Messages: []*descriptor.Message{
			{
				ID:         1,
				Name:       "EmptyMessage",
				SenderNode: "DBG",
			},

			{
				ID:          100,
				Name:        "DriverHeartbeat",
				Length:      1,
				SenderNode:  "DRIVER",
				Description: "Sync message used to synchronize the controllers",
				SendType:    descriptor.SendTypeCyclic,
				CycleTime:   time.Second,
				Signals: []*descriptor.Signal{
					{
						Name:          "Command",
						Start:         0,
						Length:        8,
						Scale:         1,
						ReceiverNodes: []string{"SENSOR", "MOTOR"},
						ValueDescriptions: []*descriptor.ValueDescription{
							{Value: 0, Description: "None"},
							{Value: 1, Description: "Sync"},
							{Value: 2, Description: "Reboot"},
							{Value: 3, Description: "Headlights On"},
						},
					},
				},
			},

			{
				ID:         101,
				Name:       "MotorCommand",
				Length:     1,
				SenderNode: "DRIVER",
				SendType:   descriptor.SendTypeCyclic,
				CycleTime:  100 * time.Millisecond,
				Signals: []*descriptor.Signal{
					{
						Name:          "Steer",
						Start:         0,
						Length:        4,
						IsSigned:      true,
						Scale:         1,
						Offset:        -5,
						Min:           -5,
						Max:           5,
						ReceiverNodes: []string{"MOTOR"},
					},
					{
						Name:          "Drive",
						Start:         4,
						Length:        4,
						Scale:         1,
						Max:           9,
						ReceiverNodes: []string{"MOTOR"},
					},
				},
			},

			{
				ID:         200,
				Name:       "SensorSonars",
				Length:     8,
				SenderNode: "SENSOR",
				SendType:   descriptor.SendTypeCyclic,
				CycleTime:  100 * time.Millisecond,
				Signals: []*descriptor.Signal{
					{
						Name:          "Mux",
						IsMultiplexer: true,
						Start:         0,
						Length:        4,
						Scale:         1,
						ReceiverNodes: []string{"DRIVER", "IO"},
					},
					{
						Name:          "ErrCount",
						Start:         4,
						Length:        12,
						Scale:         1,
						ReceiverNodes: []string{"DRIVER", "IO"},
					},
					{
						Name:             "Left",
						IsMultiplexed:    true,
						MultiplexerValue: 0,
						Start:            16,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DRIVER", "IO"},
					},
					{
						Name:             "NoFiltLeft",
						IsMultiplexed:    true,
						MultiplexerValue: 1,
						Start:            16,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DBG"},
					},
					{
						Name:             "Middle",
						IsMultiplexed:    true,
						MultiplexerValue: 0,
						Start:            28,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DRIVER", "IO"},
					},
					{
						Name:             "NoFiltMiddle",
						IsMultiplexed:    true,
						MultiplexerValue: 1,
						Start:            28,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DBG"},
					},
					{
						Name:             "Right",
						IsMultiplexed:    true,
						MultiplexerValue: 0,
						Start:            40,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DRIVER", "IO"},
					},
					{
						Name:             "NoFiltRight",
						IsMultiplexed:    true,
						MultiplexerValue: 1,
						Start:            40,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DBG"},
					},
					{
						Name:             "Rear",
						IsMultiplexed:    true,
						MultiplexerValue: 0,
						Start:            52,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DRIVER", "IO"},
					},
					{
						Name:             "NoFiltRear",
						IsMultiplexed:    true,
						MultiplexerValue: 1,
						Start:            52,
						Length:           12,
						Scale:            0.1,
						ReceiverNodes:    []string{"DBG"},
					},
				},
			},

			{
				ID:         400,
				Name:       "MotorStatus",
				Length:     3,
				SenderNode: "MOTOR",
				SendType:   descriptor.SendTypeCyclic,
				CycleTime:  100 * time.Millisecond,
				Signals: []*descriptor.Signal{
					{
						Name:          "WheelError",
						Start:         0,
						Length:        1,
						Scale:         1,
						ReceiverNodes: []string{"DRIVER", "IO"},
					},
					{
						Name:          "SpeedKph",
						Start:         8,
						Length:        16,
						Scale:         0.001,
						Unit:          "km/h",
						ReceiverNodes: []string{"DRIVER", "IO"},
					},
				},
			},

			{
				ID:         500,
				Name:       "IODebug",
				Length:     6,
				SenderNode: "IO",
				SendType:   descriptor.SendTypeEvent,
				Signals: []*descriptor.Signal{
					{
						Name:          "TestUnsigned",
						Start:         0,
						Length:        8,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
					},
					{
						Name:          "TestEnum",
						Start:         8,
						Length:        6,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
						DefaultValue:  int(examplecan.IODebug_TestEnum_Two),
						ValueDescriptions: []*descriptor.ValueDescription{
							{Value: 1, Description: "One"},
							{Value: 2, Description: "Two"},
						},
					},
					{
						Name:          "TestSigned",
						Start:         16,
						Length:        8,
						IsSigned:      true,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
					},
					{
						Name:          "TestFloat",
						Start:         24,
						Length:        8,
						Scale:         0.5,
						ReceiverNodes: []string{"DBG"},
					},
					{
						Name:          "TestBoolEnum",
						Start:         32,
						Length:        1,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
						ValueDescriptions: []*descriptor.ValueDescription{
							{Value: 0, Description: "Zero"},
							{Value: 1, Description: "One"},
						},
					},
					{
						Name:          "TestScaledEnum",
						Start:         40,
						Length:        2,
						Scale:         2,
						Min:           0,
						Max:           6,
						ReceiverNodes: []string{"DBG"},
						ValueDescriptions: []*descriptor.ValueDescription{
							{Value: 0, Description: "Zero"},
							{Value: 1, Description: "Two"},
							{Value: 2, Description: "Four"},
							{Value: 3, Description: "Six"},
						},
					},
				},
			},
			{
				ID:         600,
				Name:       "IOFloat32",
				Length:     8,
				SenderNode: "IO",
				SendType:   descriptor.SendTypeNone,
				Signals: []*descriptor.Signal{
					{
						Name:          "Float32ValueNoRange",
						Length:        32,
						IsSigned:      true,
						IsFloat:       true,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
					},
					{
						Name:          "Float32WithRange",
						Start:         32,
						Length:        32,
						IsSigned:      true,
						IsFloat:       true,
						Scale:         1,
						Min:           -100,
						Max:           100,
						ReceiverNodes: []string{"DBG"},
					},
				},
			},
			{
				ID:         700,
				Name:       "SignalNameFormatting",
				Length:     8,
				SenderNode: "IO",
				SendType:   descriptor.SendTypeNone,
				Signals: []*descriptor.Signal{
					{
						Name:          "non_capitalized_signal",
						Length:        8,
						IsSigned:      true,
						IsFloat:       false,
						Scale:         1,
						ReceiverNodes: []string{"DBG"},
					},
				},
			},
		},
	}
	input, err := os.ReadFile(exampleDBCFile)
	assert.NilError(t, err)
	result, err := Compile(exampleDBCFile, input)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Warnings) > 0 {
		t.Fatal(result.Warnings)
	}
	assert.DeepEqual(t, exampleDatabase, result.Database)
}

func TestCompile_Float64SignalWarningExpected(t *testing.T) {
	finish := runTestInDir(t, "../..")
	defer finish()
	const exampleDBCFile = "testdata/dbc-invalid/example/example_float64_signal.dbc"
	input, err := os.ReadFile(exampleDBCFile)
	assert.NilError(t, err)
	result, err := Compile(exampleDBCFile, input)
	if err != nil {
		t.Fatal(err)
	}
	// We expect one warning from unsupported float64 signal
	assert.Equal(t, len(result.Warnings), 1)
}

func TestCompile_Float32InvalidSignalNameWarningExpected(t *testing.T) {
	finish := runTestInDir(t, "../..")
	defer finish()
	const exampleDBCFile = "testdata/dbc-invalid/example/example_float32_invalid_signal_name.dbc"
	input, err := os.ReadFile(exampleDBCFile)
	assert.NilError(t, err)
	result, err := Compile(exampleDBCFile, input)
	if err != nil {
		t.Fatal(err)
	}
	// We expect one warning for incorrect signal name in SIGVAL_TYPE_ declaration
	assert.Equal(t, len(result.Warnings), 1)
}

func TestCompile_Float32InvalidSignalLengthWarningExpected(t *testing.T) {
	finish := runTestInDir(t, "../..")
	defer finish()
	const exampleDBCFile = "testdata/dbc-invalid/example/example_float32_invalid_signal_length.dbc"
	input, err := os.ReadFile(exampleDBCFile)
	assert.NilError(t, err)
	result, err := Compile(exampleDBCFile, input)
	if err != nil {
		t.Fatal(err)
	}
	// We expect one warning for incorrect signal length in declaration of float32 signal
	assert.Equal(t, len(result.Warnings), 1)
}

func Test_CopyFrom_PreservesOutOfRangeValues(t *testing.T) {
	descriptor := examplecan.Messages().MotorCommand
	frame := can.Frame{
		ID:         descriptor.ID,
		Length:     descriptor.Length,
		IsExtended: descriptor.IsExtended,
	}
	// 0xF is 15, but max is set to 9
	descriptor.Drive.MarshalUnsigned(&frame.Data, 0xF)
	// Unmarshal out of bounds value
	original := examplecan.NewMotorCommand()
	if err := original.UnmarshalFrame(frame); err != nil {
		t.Errorf("Failed to unmarshal frame: %v", err)
	}
	// When we CopyFrom original message to m2
	m2 := examplecan.NewMotorCommand().CopyFrom(original)
	// Then we expect the messages and the frames to be identical
	if !reflect.DeepEqual(m2, original) {
		t.Errorf("Expected new message (%v) and original (%v) to be identical", m2, original)
	}
	if m2.Frame() != original.Frame() {
		t.Errorf("Expected frames of messages to be identical (%v != %v)", m2.Frame(), original.Frame())
	}
}

func TestCompile_MetadataInvalidSignalReference(t *testing.T) {
	finish := runTestInDir(t, "../..")
	defer finish()
	const exampleDBCFile = "testdata/dbc-invalid/example/example_metadata_invalid_signal_reference.dbc"
	input, err := os.ReadFile(exampleDBCFile)
	assert.NilError(t, err)
	result, err := Compile(exampleDBCFile, input)
	if err != nil {
		t.Fatal(err)
	}
	// We expect one warning for incorrect signal length in declaration of float32 signal
	assert.Equal(t, len(result.Warnings), 1)
}
