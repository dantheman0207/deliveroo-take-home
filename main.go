package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
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
	year       CronField
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
	currentYear := time.Now().Year()
	cronFields := CronFields{
		minute:     CronField{label: "minute", min: 0, max: 59},
		hour:       CronField{label: "hour", min: 0, max: 23},
		dayOfMonth: CronField{label: "day of month", min: 1, max: 31},
		month:      CronField{label: "month", min: 1, max: 12},
		dayOfWeek:  CronField{label: "day of week", min: 0, max: 6},
		year:       CronField{label: "year", min: currentYear - 2, max: currentYear + 25},
	}
	err := cronFields.expandCronFields(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		if exitOnError {
			os.Exit(3)
		}
		return
	}
	commandIndex := 6
	if cronFields.year.value != "" {
		commandIndex = 7
	}
	cronFields.command = strings.Join(os.Args[commandIndex:], " ")
	cronFields.printCronFields()
}

func (c *CronFields) expandCronFields(args []string) error {
	result, err := c.minute.expandCronField(args[0])
	if err != nil {
		return err
	}
	c.minute.value = result
	result, err = c.hour.expandCronField(args[1])
	if err != nil {
		return err
	}
	c.hour.value = result
	result, err = c.dayOfMonth.expandCronField(args[2])
	if err != nil {
		return err
	}
	c.dayOfMonth.value = result
	result, err = c.month.expandCronField(args[3])
	if err != nil {
		return err
	}
	c.month.value = result
	result, err = c.dayOfWeek.expandCronField(args[4])
	if err != nil {
		return err
	}
	c.dayOfWeek.value = result
	result, err = c.year.expandCronField(args[5])
	if err == nil && result != "" {
		c.year.value = result
	}

	return nil
}

func (c *CronFields) printCronFields() {
	// Output the results
	fmt.Printf("%-14s%s\n", c.minute.label, c.minute.value)
	fmt.Printf("%-14s%s\n", c.hour.label, c.hour.value)
	fmt.Printf("%-14s%s\n", c.dayOfMonth.label, c.dayOfMonth.value)
	fmt.Printf("%-14s%s\n", c.month.label, c.month.value)
	fmt.Printf("%-14s%s\n", c.dayOfWeek.label, c.dayOfWeek.value)
	if c.year.value != "" {
		fmt.Printf("%-14s%s\n", c.year.label, c.year.value)
	}
	fmt.Printf("%-14s%s\n", "command", c.command)
}

// expandCronField expands a single cron field into a list of valid values.
// It supports asterisks, commas, ranges, and steps.
// It sets the value of the cron field to the expanded values.
// It handles any error encountered.
func (cf *CronField) expandCronField(field string) (string, error) {
	var result string
	var err error
	switch {
	case field == "*":
		result, err = cf.expandCronFieldAsterisk()
		if err != nil {
			return "", err
		}
	case strings.Contains(field, ","):
		result, err = cf.expandCronFieldComma(field)
		if err != nil {
			return "", err
		}
	case strings.Contains(field, "/"):
		result, err = cf.expandCronFieldStep(field)
		if err != nil {
			return "", err
		}
	case strings.Contains(field, "-"):
		result, err = cf.expandCronFieldRange(field)
		if err != nil {
			return "", err
		}
	default:
		num, err := strconv.Atoi(field)
		if err != nil {
			return "", fmt.Errorf("invalid input: %s", field)
		}
		if num < cf.min || num > cf.max {
			return "", fmt.Errorf("value out of range: %d", num)
		}
		result = strconv.Itoa(num)
	}
	return result, nil
}

func (cf *CronField) expandCronFieldAsterisk() (string, error) {
	result := []int{}
	for i := cf.min; i <= cf.max; i++ {
		result = append(result, i)
	}
	resultString := strings.Trim(fmt.Sprint(result), "[]")
	return resultString, nil
}

func (cf *CronField) expandCronFieldComma(field string) (string, error) {
	result := []int{}
	for _, v := range strings.Split(field, ",") {
		nums, err := cf.expandCronField(v)
		if err != nil {
			return "", err
		}
		for _, v := range strings.Split(nums, " ") {
			num, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("invalid input: %s", v)
			}
			if num < cf.min || num > cf.max {
				return "", fmt.Errorf("value out of range: %d", num)
			}
			result = append(result, num)
		}
	}
	sort.Ints(result)
	resultString := strings.Trim(fmt.Sprint(result), "[]")
	return resultString, nil
}

func (cf *CronField) expandCronFieldRange(field string) (string, error) {
	parts := strings.Split(field, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid range: %s", field)
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid start of range: %s", parts[0])
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid end of range: %s", parts[1])
	}
	if start < cf.min || start > cf.max || end < cf.min || end > cf.max || start > end {
		return "", fmt.Errorf("invalid range: %s", field)
	}
	result := []int{}
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	resultString := strings.Trim(fmt.Sprint(result), "[]")
	return resultString, nil
}

func (cf *CronField) expandCronFieldStep(field string) (string, error) {
	parts := strings.Split(field, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid step: %s", field)
	}
	step, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid step value: %s", parts[1])
	}
	if step <= 0 {
		return "", fmt.Errorf("step must be positive: %d", step)
	}
	if parts[0] != "*" { // handle step with range
		parts = strings.Split(parts[0], "-")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid range for step base: %s", parts[0])
		}
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", fmt.Errorf("invalid start of step range: %s", parts[0])
		}
		cf.min = start
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("invalid end of step range: %s", parts[1])
		}
		cf.max = end
	}
	if step > cf.max-cf.min {
		return "", fmt.Errorf("step is larger than the range")
	}
	result := []int{}
	for i := cf.min; i <= cf.max; i += step {
		result = append(result, i)
	}
	resultString := strings.Trim(fmt.Sprint(result), "[]")
	return resultString, nil
}
