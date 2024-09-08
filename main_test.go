package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name: "Basic test",
			args: []string{"*/15", "0", "1,15", "*", "1-5", "/usr/bin/find"},
			expected: `minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
`,
		},
		// Add more test cases here
		{
			name: "Test with all *",
			args: []string{"*", "*", "*", "*", "*", "/usr/bin/find"},
			expected: `minute        0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59
hour          0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
day of month  1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   0 1 2 3 4 5 6
command       /usr/bin/find
`,
		},
		{
			name: "Test with all min",
			args: []string{"0", "0", "1", "1", "0", "/usr/bin/find"},
			expected: `minute        0
hour          0
day of month  1
month         1
day of week   0
command       /usr/bin/find
`,
		},
		{
			name: "Test with all max",
			args: []string{"59", "23", "31", "12", "6", "/usr/bin/find"},
			expected: `minute        59
hour          23
day of month  31
month         12
day of week   6
command       /usr/bin/find
`,
		},
		{
			name:     "Test with invalid input",
			args:     []string{"invalid", "0", "1", "1", "1", "/usr/bin/find"},
			expected: "invalid input: invalid\n",
		},
		// No arguments provided
		{
			name:     "No arguments provided",
			args:     []string{},
			expected: `Error: not enough arguments provided. Example usage: go run main.go */15 0 1,15 * 1-5 /usr/bin/find`,
		},
		// No command provided
		{
			name:     "No command provided",
			args:     []string{"*/15", "0", "1,15", "*", "1-5"},
			expected: `Error: not enough arguments provided. Example usage: go run main.go */15 0 1,15 * 1-5 /usr/bin/find`,
		},
	}

	exitOnError = false
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Test panicked: %v", r)
				}
			}()

			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			_ = io.MultiWriter(w, old)
			os.Stdout = w

			// // Set up arguments
			os.Args = append([]string{"cmd"}, tt.args...)

			// // Run main
			main()

			// // Restore stdout
			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			io.Copy(&buf, r)
			actual := buf.String()

			if strings.TrimSpace(actual) != strings.TrimSpace(tt.expected) {
				t.Errorf("Expected output:\n%s\nGot:\n%s", tt.expected, actual)
			}
		})
	}
	exitOnError = true
}

func TestExpandCronField(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		min, max int
		expected string
	}{
		{"Asterisk", "*", 0, 23, "0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23"},
		{"Comma", "1,3,5", 0, 5, "1 3 5"},
		{"Range", "1-3", 0, 5, "1 2 3"},
		{"Step", "*/2", 0, 5, "0 2 4"},
		{"Single", "3", 0, 5, "3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := CronField{min: tt.min, max: tt.max}
			cf.expandCronField(tt.field)
			if cf.value != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, cf.value)
			}
		})
	}
}

func TestExpandCronFieldInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		min, max int
		expected string
	}{
		{"InvalidInput", "invalid", 0, 5, "invalid input: invalid"},
		{"InvalidGroupInput", "1,NAN,3", 0, 5, "invalid input: NAN"},
		{"InvalidGroupInputOutOfRange", "1,30,3", 0, 5, "value out of range: 30"},
		{"InvalidRangeInputDash", "1-NAN-NAN2", 0, 5, "invalid range: 1-NAN-NAN2"},
		{"InvalidRangeInputStart", "NAN-NAN2", 0, 5, "invalid start of range: NAN"},
		{"InvalidRangeInputEnd", "1-NAN", 0, 5, "invalid end of range: NAN"},
		{"InvalidRange", "1-30", 0, 5, "invalid range: 1-30"},
		{"InvalidStep", "*/30", 0, 5, "step is larger than the range"},
		{"InvalidStepValue", "*/0", 0, 5, "step must be positive: 0"},
		{"InvalidStepInput", "*/1/2", 0, 5, "invalid step: */1/2"},
		{"InvalidStepInput", "*/NAN", 0, 5, "invalid step value: NAN"},
		{"InvalidInputForRange", "8", 0, 5, "value out of range: 8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cf := CronField{min: tt.min, max: tt.max}
			err := cf.expandCronField(tt.field)
			if err == nil {
				t.Errorf("Expected error, got nil")
				return
			}
			if err.Error() != tt.expected {
				t.Errorf("Expected error message %q, got %q", tt.expected, err.Error())
			}
		})
	}
}

func TestExpandCronFieldsErrors(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{"Invalid step input for minute", []string{"*/60", "0", "1,15", "*", "1-5"}, "step is larger than the range"},
		{"Invalid range input for day of week", []string{"*/15", "0", "1,15", "*", "1-30"}, "invalid range: 1-30"},
		{"Invalid comma input for day of month", []string{"*/15", "0", "1,50", "*", "1-5"}, "value out of range: 50"},
		{"Invalid input for month", []string{"*/15", "0", "1,15", "50", "1-5", "/usr/bin/find"}, "value out of range: 50"},
		{"Invalid input for hour", []string{"*/15", "050", "1,15", "0", "1-5", "/usr/bin/find"}, "value out of range: 50"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cronFields := CronFields{
				minute:     CronField{label: "minute", min: 0, max: 59},
				hour:       CronField{label: "hour", min: 0, max: 23},
				dayOfMonth: CronField{label: "day of month", min: 1, max: 31},
				month:      CronField{label: "month", min: 1, max: 12},
				dayOfWeek:  CronField{label: "day of week", min: 0, max: 6},
				command:    "/usr/bin/find",
			}
			err := cronFields.expandCronFields(tt.args)
			if err == nil {
				t.Errorf("Expected error, got nil")
			}
			if err.Error() != tt.expected {
				t.Errorf("Expected error message %q, got %q", tt.expected, err.Error())
			}
		})
	}
}
