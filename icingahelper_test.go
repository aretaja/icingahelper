package icingahelper

import "testing"

func TestNewCheck(t *testing.T) {
	check := NewCheck("CPU")
	if check.name != "CPU" {
		t.Error("NewCheck name - Expected CPU, got ", check.name)
	}
	if check.retVal != 3 {
		t.Error("NewCheck retVal - Expected 3, got ", check.retVal)
	}
	if check.perf != nil {
		t.Error("NewCheck perf - Expected nil, got ", check.perf)
	}
	if check.unkn != nil {
		t.Error("NewCheck unkn - Expected nil, got ", check.unkn)
	}
	if check.crit != nil {
		t.Error("NewCheck crit - Expected nil, got ", check.crit)
	}
	if check.warn != nil {
		t.Error("NewCheck warn - Expected nil, got ", check.warn)
	}
	if check.ok != nil {
		t.Error("NewCheck ok - Expected nil, got ", check.ok)
	}
}

func TestRetVal(t *testing.T) {
	check := NewCheck("CPU")
	if check.RetVal() != 3 {
		t.Error("RetVal - Expected 3, got ", check.RetVal())
	}
}

func TestAlarmLevel(t *testing.T) {
	check := NewCheck("CPU")
	if r, err := check.AlarmLevel(5, "6", "7"); r != 0 || err != nil {
		t.Error("AlarmLevel OK - Expected 0 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(7, "6", "7"); r != 1 || err != nil {
		t.Error("AlarmLevel WARN - Expected 1 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(8, "6", "7"); r != 2 || err != nil {
		t.Error("AlarmLevel CRIT - Expected 2 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(3, "3:5", "2:6"); r != 0 || err != nil {
		t.Error("AlarmLevel Outside OK - Expected 0 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(2, "3:5", "2:6"); r != 1 || err != nil {
		t.Error("AlarmLevel Outside WARN lower - Expected 1 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(6, "3:5", "2:6"); r != 1 || err != nil {
		t.Error("AlarmLevel Outside WARN upper - Expected 1 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(1, "3:5", "2:6"); r != 2 || err != nil {
		t.Error("AlarmLevel Outside CRIT lower - Expected 1 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(7, "3:5", "2:6"); r != 2 || err != nil {
		t.Error("AlarmLevel Outside CRIT upper - Expected 1 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(2, "@2:6", "@3:5"); r != 0 || err != nil {
		t.Error("AlarmLevel Inside OK lower - Expected 0 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(7, "@2:6", "@3:5"); r != 0 || err != nil {
		t.Error("AlarmLevel Inside OK upper - Expected 1 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(3, "@2:6", "@3:5"); r != 1 || err != nil {
		t.Error("AlarmLevel Inside WARN - Expected 1 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(4, "@2:6", "@3:5"); r != 2 || err != nil {
		t.Error("AlarmLevel Inside CRIT - Expected 1 and nil, got ", r, err)
	}
	if r, err := check.AlarmLevel(4, "@2:6:", "@3:5-"); r != 3 || err == nil {
		t.Error("AlarmLevel Bad range - Expected 3 and error, got ", r, err)
	}
}

func TestAddPerfData(t *testing.T) {
	check := NewCheck("CPU")
	check.AddPerfData("cpu usage", "20", "%", "80", "90", "0", "100")
	check.AddPerfData("cpu load", "20", "%", "80", "90", "0", "100")
	if check.perf == nil || check.perf[0] != "cpu usage=20%;80;90;0;100" || check.perf[1] != "cpu load=20%;80;90;0;100" {
		t.Errorf("AddPerfData - Expected []string{\"cpu usage=20%%;80;90;0;100\", \"cpu load=20%%;80;90;0;100\"}, got %#v", check.perf)
	}
}

func TestAddMsg(t *testing.T) {
	check := NewCheck("CPU")
	check.AddMsg(2, "shortCmsg", "longCmsg")
	check.AddMsg(2, "shortCmsg", "longCmsg")
	if check.crit == nil || check.crit[0].short != "shortCmsg" || check.crit[0].long != "longCmsg" || check.crit[1].short != "shortCmsg" || check.crit[1].long != "longCmsg" {
		t.Errorf("AddMsg - Expected []icingahelper.msg{icingahelper.msg{short:\"shortCmsg\", long:\"longCmsg\"}, icingahelper.msg{short:\"shortCmsg\", long:\"longCmsg\"}}, got %#v", check.crit)
	}

	check.AddMsg(1, "shortWmsg", "longWmsg")
	check.AddMsg(1, "shortWmsg", "longWmsg")
	if check.warn == nil || check.warn[0].short != "shortWmsg" || check.warn[0].long != "longWmsg" || check.warn[1].short != "shortWmsg" || check.warn[1].long != "longWmsg" {
		t.Errorf("AddMsg - Expected []icingahelper.msg{icingahelper.msg{short:\"shortWmsg\", long:\"longWmsg\"}, icingahelper.msg{short:\"shortWmsg\", long:\"longWmsg\"}}, got %#v", check.warn)
	}

	check.AddMsg(0, "shortOmsg", "longOmsg")
	check.AddMsg(0, "shortOmsg", "longOmsg")
	if check.ok == nil || check.ok[0].short != "shortOmsg" || check.ok[0].long != "longOmsg" || check.ok[1].short != "shortOmsg" || check.ok[1].long != "longOmsg" {
		t.Errorf("AddMsg - Expected []icingahelper.msg{icingahelper.msg{short:\"shortOmsg\", long:\"longOmsg\"}, icingahelper.msg{short:\"shortOmsg\", long:\"longOmsg\"}}, got %#v", check.ok)
	}

	check.AddMsg(3, "shortUmsg", "longUmsg")
	check.AddMsg(3, "shortUmsg", "longUmsg")
	if check.unkn == nil || check.unkn[0].short != "shortUmsg" || check.unkn[0].long != "longUmsg" || check.unkn[1].short != "shortUmsg" || check.unkn[1].long != "longUmsg" {
		t.Errorf("AddMsg - Expected []icingahelper.msg{icingahelper.msg{short:\"shortUmsg\", long:\"longUmsg\"}, icingahelper.msg{short:\"shortUmsg\", long:\"longUmsg\"}}, got %#v", check.unkn)
	}
}

func TestFinalMsg(t *testing.T) {
	check := NewCheck("CPU")
	check.AddPerfData("cpu usage", "20", "%", "80", "90", "0", "100")

	for i := 0; i < 4; i++ {
		check.AddMsg(i, "shortmsg", "longmsg")
		check.AddMsg(i, "shortmsg", "longmsg")
	}

	if r := check.FinalMsg(); r != "CPU: UNKNOWN - shortmsg(c); shortmsg(c); shortmsg(w); shortmsg(w); shortmsg(u); shortmsg(u); shortmsg; shortmsg |cpu usage=20%;80;90;0;100\n\nlongmsg(c)\nlongmsg(c)\nlongmsg(w)\nlongmsg(w)\nlongmsg(u)\nlongmsg(u)\nlongmsg(ok)\nlongmsg(ok)" {
		t.Errorf("FinalMsg UNKN - Expected \"CPU: UNKNOWN - shortmsg(c); shortmsg(c); shortmsg(w); shortmsg(w); shortmsg(u); shortmsg(u); shortmsg; shortmsg |cpu usage=20%%;80;90;0;100\\n\\nlongmsg(c)\\nlongmsg(c)\\nlongmsg(w)\\nlongmsg(w)\\nlongmsg(u)\\nlongmsg(u)\\nlongmsg(ok)\\nlongmsg(ok)\", got %#v", r)
	}

	check.retVal = 2
	if r := check.FinalMsg(); r != "CPU: CRITICAL - shortmsg(c); shortmsg(c); shortmsg(w); shortmsg(w); shortmsg(u); shortmsg(u); shortmsg; shortmsg |cpu usage=20%;80;90;0;100\n\nlongmsg(c)\nlongmsg(c)\nlongmsg(w)\nlongmsg(w)\nlongmsg(u)\nlongmsg(u)\nlongmsg(ok)\nlongmsg(ok)" {
		t.Errorf("FinalMsg CRIT - Expected \"CPU: CRITICAL - shortmsg(c); shortmsg(c); shortmsg(w); shortmsg(w); shortmsg(u); shortmsg(u); shortmsg; shortmsg |cpu usage=20%%;80;90;0;100\\n\\nlongmsg(c)\\nlongmsg(c)\\nlongmsg(w)\\nlongmsg(w)\\nlongmsg(u)\\nlongmsg(u)\\nlongmsg(ok)\\nlongmsg(ok)\", got %#v", r)
	}

	check.retVal = 1
	if r := check.FinalMsg(); r != "CPU: WARNING - shortmsg(c); shortmsg(c); shortmsg(w); shortmsg(w); shortmsg(u); shortmsg(u); shortmsg; shortmsg |cpu usage=20%;80;90;0;100\n\nlongmsg(c)\nlongmsg(c)\nlongmsg(w)\nlongmsg(w)\nlongmsg(u)\nlongmsg(u)\nlongmsg(ok)\nlongmsg(ok)" {
		t.Errorf("FinalMsg WARN - Expected \"CPU: WARNING - shortmsg(c); shortmsg(c); shortmsg(w); shortmsg(w); shortmsg(u); shortmsg(u); shortmsg; shortmsg |cpu usage=20%%;80;90;0;100\\n\\nlongmsg(c)\\nlongmsg(c)\\nlongmsg(w)\\nlongmsg(w)\\nlongmsg(u)\\nlongmsg(u)\\nlongmsg(ok)\\nlongmsg(ok)\", got %#v", r)
	}

	check.retVal = 0
	if r := check.FinalMsg(); r != "CPU: OK - shortmsg(c); shortmsg(c); shortmsg(w); shortmsg(w); shortmsg(u); shortmsg(u); shortmsg; shortmsg |cpu usage=20%;80;90;0;100\n\nlongmsg(c)\nlongmsg(c)\nlongmsg(w)\nlongmsg(w)\nlongmsg(u)\nlongmsg(u)\nlongmsg(ok)\nlongmsg(ok)" {
		t.Errorf("FinalMsg OK - Expected \"CPU: OK - shortmsg(c); shortmsg(c); shortmsg(w); shortmsg(w); shortmsg(u); shortmsg(u); shortmsg; shortmsg |cpu usage=20%%;80;90;0;100\\n\\nlongmsg(c)\\nlongmsg(c)\\nlongmsg(w)\\nlongmsg(w)\\nlongmsg(u)\\nlongmsg(u)\\nlongmsg(ok)\\nlongmsg(ok)\", got %#v", r)
	}
}
