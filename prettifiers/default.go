package prettifiers

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/alecthomas/repr"
	"github.com/fatih/color"
	"github.com/globocom/prettylog/config"
	"github.com/globocom/prettylog/parsers"
)

const (
	SEPARATOR       = " "
	FIELD_SEPARATOR = "="
	FIELD_QUOTES    = "\""
)

type DefaultPrettifier struct{}

func (p *DefaultPrettifier) Prettify(line *parsers.ParsedLine) string {
	settings := config.GetSettings()
	buffer := &bytes.Buffer{}

	if settings.Timestamp.Visible {
		ts, err := time.Parse(time.RFC3339, line.Timestamp)
		if err != nil || settings.Timestamp.Format == "" {
			writeTo(buffer, line.Timestamp, 0, settings.Timestamp.Color)
		} else {
			writeTo(buffer, ts.Format(settings.Timestamp.Format), 0, settings.Timestamp.Color)
		}
	}

	if settings.Logger.Visible && line.Logger != "" {
		writeTo(buffer, line.Logger, settings.Logger.Padding, settings.Logger.Color)
	}

	if settings.Caller.Visible {
		writeTo(buffer, line.Caller, settings.Caller.Padding, settings.Caller.Color)
	}

	if settings.Level.Visible {
		writeTo(buffer, strings.ToUpper(line.Level), settings.Level.Padding, settings.Level.GetColorAttr(line.Level))
	}

	writeTo(buffer, line.Message, settings.Message.Padding, settings.Message.Color)
	writeFieldsTo(buffer, line.Fields, settings.Level.GetColorAttr(line.Level))

	return buffer.String()
}

func writeTo(buffer *bytes.Buffer, value string, padding int, colorAttrs []color.Attribute) {
	color := parseColor(colorAttrs)
	value = padRight(value, padding)

	if color != nil {
		value = color.Sprint(value)
	}

	buffer.WriteString(value)
	buffer.WriteString(SEPARATOR)
}

func writeFieldsTo(buffer *bytes.Buffer, fields [][]string, colorsAttrs []color.Attribute) {
	colorAttr := parseColor(colorsAttrs)
	settings := config.GetSettings()
	defaultFieldSettings := config.Field{
		Visible: true,
		Padding: 0,
		Color:   []color.Attribute{color.FgBlack, color.Faint},
	}

	for _, field := range fields {
		fieldSettings, fieldSettingsExist := settings.Fields[field[0]]
		if !fieldSettingsExist {
			fieldSettings = defaultFieldSettings
		}

		colorAttr = parseColor(fieldSettings.Color)
		if !fieldSettings.Visible {
			continue
		}

		if colorAttr != nil {
			buffer.WriteString(colorAttr.Sprint(field[0]))
		} else {
			buffer.WriteString(field[0])
		}
		buffer.WriteString(FIELD_SEPARATOR)

		fieldValue := field[1]
		if strings.Contains(fieldValue, " ") {
			fieldValue = fmt.Sprintf("%s%s%s", FIELD_QUOTES, fieldValue, FIELD_QUOTES)
		}

		repr.Println(fieldValue)

		fieldColor := parseColor(fieldSettings.Color)
		buffer.WriteString(fieldColor.Sprint(fieldValue))

		buffer.WriteString(fieldValue)
		buffer.WriteString(SEPARATOR)
	}
}

func padRight(str string, size int) string {
	size = size - len(str)
	if size < 0 {
		size = 0
	}
	return str + strings.Repeat(" ", size)
}

func parseColor(attributes []color.Attribute) *color.Color {
	var c *color.Color

	if len(attributes) > 0 {
		c = color.New(attributes[0])
	}

	if len(attributes) > 1 {
		c.Add(attributes[1])
	}

	if len(attributes) > 2 {
		c.Add(attributes[2])
	}

	return c
}
