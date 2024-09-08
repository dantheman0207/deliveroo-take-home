package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

// CronField represents a single field in a cron expression.
type CronField struct {
	label string
	min   int
	max   int
	value string
}

// CronFields represents all fields in a cron expression, including the command.
type CronFields struct {
	minute     CronField
	hour       CronField
	dayOfMonth CronField
	month      CronField
	dayOfWeek  CronField
	command    string
}

var (
	exitOnError = true
)

// main is the entry point of the program. It parses command-line arguments,
// expands the cron expression, and prints the results.
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered in main: %v\n", r)
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	if len(os.Args) < 7 {
		fmt.Println("Error: not enough arguments provided. Example usage: go run main.go */15 0 1,15 * 1-5 /usr/bin/find")
		if exitOnError {
			os.Exit(2)
		}
		return
	}

	// Parse and expand each field
	cronFields := CronFields{
		minute:     CronField{label: "minute", min: 0, max: 59},
		hour:       CronField{label: "hour", min: 0, max: 23},
		dayOfMonth: CronField{label: "day of month", min: 1, max: 31},
		month:      CronField{label: "month", min: 1, max: 12},
		dayOfWeek:  CronField{label: "day of week", min: 0, max: 6},
		command:    os.Args[6],
	}
	err := cronFields.expandCronFields(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		if exitOnError {
			os.Exit(4)
		}
		return
	}
	cronFields.printCronFields()
}

func (c *CronFields) expandCronFields(args []string) error {
	err := c.minute.expandCronField(args[0])
	if err != nil {
		return err
	}
	err = c.hour.expandCronField(args[1])
	if err != nil {
		return err
	}
	err = c.dayOfMonth.expandCronField(args[2])
	if err != nil {
		return err
	}
	err = c.month.expandCronField(args[3])
	if err != nil {
		return err
	}
	err = c.dayOfWeek.expandCronField(args[4])
	if err != nil {
		return err
	}
	return nil
}

func (c *CronFields) printCronFields() {
	// Output the results
	fmt.Printf("%-14s%s\n", "minute", c.minute.value)
	fmt.Printf("%-14s%s\n", "hour", c.hour.value)
	fmt.Printf("%-14s%s\n", "day of month", c.dayOfMonth.value)
	fmt.Printf("%-14s%s\n", "month", c.month.value)
	fmt.Printf("%-14s%s\n", "day of week", c.dayOfWeek.value)
	fmt.Printf("%-14s%s\n", "command", c.command)
}

// expandCronField expands a single cron field into a list of valid values.
// It supports asterisks, commas, ranges, and steps.
// It sets the value of the cron field to the expanded values.
// It handles any error encountered.
func (cf *CronField) expandCronField(field string) error {
	switch {
	case field == "*":
		return cf.expandCronFieldAsterisk()
	case strings.Contains(field, ","):
		return cf.expandCronFieldComma(field)
	case strings.Contains(field, "-"):
		return cf.expandCronFieldRange(field)
	case strings.Contains(field, "/"):
		return cf.expandCronFieldStep(field)
	default:
		num, err := strconv.Atoi(field)
		if err != nil {
			return fmt.Errorf("invalid input: %s", field)
		}
		if num < cf.min || num > cf.max {
			return fmt.Errorf("value out of range: %d", num)
		}
		cf.value = strconv.Itoa(num)
		return nil
	}
}

func (cf *CronField) expandCronFieldAsterisk() error {
	result := []int{}
	for i := cf.min; i <= cf.max; i++ {
		result = append(result, i)
	}
	cf.value = strings.Trim(fmt.Sprint(result), "[]")
	return nil
}

func (cf *CronField) expandCronFieldComma(field string) error {
	result := []int{}
	for _, v := range strings.Split(field, ",") {
		num, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid input: %s", v)
		}
		if num < cf.min || num > cf.max {
			return fmt.Errorf("value out of range: %d", num)
		}
		result = append(result, num)
	}
	cf.value = strings.Trim(fmt.Sprint(result), "[]")
	return nil
}

func (cf *CronField) expandCronFieldRange(field string) error {
	parts := strings.Split(field, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range: %s", field)
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid start of range: %s", parts[0])
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid end of range: %s", parts[1])
	}
	if start < cf.min || start > cf.max || end < cf.min || end > cf.max || start > end {
		return fmt.Errorf("invalid range: %s", field)
	}
	result := []int{}
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	cf.value = strings.Trim(fmt.Sprint(result), "[]")
	return nil
}

func (cf *CronField) expandCronFieldStep(field string) error {
	parts := strings.Split(field, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid step: %s", field)
	}
	step, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid step value: %s", parts[1])
	}
	if step <= 0 {
		return fmt.Errorf("step must be positive: %d", step)
	}
	if step > cf.max-cf.min {
		return fmt.Errorf("step is larger than the range")
	}
	result := []int{}
	for i := cf.min; i <= cf.max; i += step {
		result = append(result, i)
	}
	cf.value = strings.Trim(fmt.Sprint(result), "[]")
	return nil
}
