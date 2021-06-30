// Package icingahelper implements icingaCheck object and some functions to
// ease Icinga plugin development.
package icingahelper

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Exported part

// Version of release
const Version = "1.0.0"

// icingaCheck object
type IcingaCheck struct {
	name                 string   // name - first word in final return message.
	perf                 []string // perf - performance data
	retVal               int      // retVal - plugin exit value
	unkn, crit, warn, ok []msg    // messages of results which have corrsponding level
}

// Message data
type msg struct {
	short string // short message
	long  string // long message
}

// Initialize new icingaCheck object
func NewCheck(name string) *IcingaCheck {
	return &IcingaCheck{
		name:   name,
		retVal: 3,
	}
}

// Get icingaCheck retVal
func (c *IcingaCheck) RetVal() int {
	return c.retVal
}

// Check threshold
//  Returns alarm level (int) and error if any
func (c *IcingaCheck) AlarmLevel(v int64, wa, cr string) (int, error) {
	level := 3
	re, _ := regexp.Compile(`^(@)?(?:(-?[0-9]*):)?(?:(-?[0-9]*))$`)

	// Parse correct types from submitted threshold strings
	warn := re.FindSubmatch([]byte(wa))
	crit := re.FindSubmatch([]byte(cr))

	if len(warn) < 1 || len(crit) < 1 {
		return level, fmt.Errorf("not valid threshold w - %s, c - %s", wa, cr)
	}

	var wInside, cInside bool
	var wMin, cMin, wMax, cMax int64 = math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64

	if string(warn[1]) == "@" {
		wInside = true
	}

	if string(crit[1]) == "@" {
		cInside = true
	}

	if str := string(warn[2]); str != "" {
		wMin, _ = strconv.ParseInt(string(warn[2]), 10, 64)
	}

	if str := string(warn[3]); str != "" {
		wMax, _ = strconv.ParseInt(string(warn[3]), 10, 64)
	}

	if str := string(crit[2]); str != "" {
		cMin, _ = strconv.ParseInt(string(crit[2]), 10, 64)
	}

	if str := string(crit[3]); str != "" {
		cMax, _ = strconv.ParseInt(string(crit[3]), 10, 64)
	}

	// Calculate alarm level based on threshold
	isAlarm := func(i bool, min, max int64) bool {
		if i {
			if v > min && v < max {
				return true
			} else {
				return false
			}
		} else {
			if v < min || v > max {
				return true
			} else {
				return false
			}
		}
	}

	if isAlarm(cInside, cMin, cMax) {
		level = 2
	} else if isAlarm(wInside, wMin, wMax) {
		level = 1
	} else {
		level = 0
	}

	// Change icingaCheck retVal if needed
	if (c.retVal == 3 && level != 3) || (c.retVal != 3 && level != 3 && level > c.retVal) {
		c.retVal = level
	}

	return level, nil
}

// Add perormance data
// unit - "us", "ms", "s", "%", "b", "kb", "mb", "gb", "tb", "c", or the empty string
// max, min - must be "" if not defined
// warn, crit - [[@]<int>:]<int>
//  fe. addPerfData("cpu usage", "20", "%", "0", "100", "80", "90")
func (c *IcingaCheck) AddPerfData(label, value, unit, warn, crit, min, max string) {

	out := fmt.Sprintf("%s=%s%s;%s;%s;%s;%s", label, value, unit, warn, crit, min, max)

	c.perf = append(c.perf, out)
}

// Add to check return message(s)
func (c *IcingaCheck) AddMsg(level int, short, long string) {
	m := msg{
		short: short,
		long:  long,
	}

	switch level {
	case 2:
		t := c.crit
		t = append(t, m)
		c.crit = t
	case 1:
		t := c.warn
		t = append(t, m)
		c.warn = t
	case 0:
		t := c.ok
		t = append(t, m)
		c.ok = t
	default:
		t := c.unkn
		t = append(t, m)
		c.unkn = t
	}
}

// Returns plugin output message
func (c *IcingaCheck) FinalMsg() string {
	level := "UNKNOWN"

	switch c.retVal {
	case 2:
		level = "CRITICAL"
	case 1:
		level = "WARNING"
	case 0:
		level = "OK"
	}

	var sm, lm []string

	if c.crit != nil {
		for _, v := range c.crit {
			sm = append(sm, fmt.Sprintf("%s(c)", v.short))

			if v.long != "" {
				lm = append(lm, fmt.Sprintf("%s(c)", v.long))
			}
		}
	}

	if c.warn != nil {
		for _, v := range c.warn {
			sm = append(sm, fmt.Sprintf("%s(w)", v.short))

			if v.long != "" {
				lm = append(lm, fmt.Sprintf("%s(w)", v.long))
			}
		}
	}

	if c.unkn != nil {
		for _, v := range c.unkn {
			sm = append(sm, fmt.Sprintf("%s(u)", v.short))

			if v.long != "" {
				lm = append(lm, fmt.Sprintf("%s(u)", v.long))
			}
		}
	}

	if c.ok != nil {
		for _, v := range c.ok {
			sm = append(sm, v.short)

			if v.long != "" {
				lm = append(lm, fmt.Sprintf("%s(ok)", v.long))
			}
		}
	}

	perf := ""
	if c.perf != nil {
		perf = fmt.Sprint("|", strings.Join(c.perf, " "))
	}

	msg := fmt.Sprintf("%s: %s - %s %s\n", c.name, level, strings.Join(sm, "; "), perf)
	if len(lm) > 0 {
		msg = fmt.Sprintf("%s\n%s", msg, strings.Join(lm, "\n"))
	}
	return msg
}
