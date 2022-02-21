package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gopowersupply/intelcpu"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

type CpuGovernor struct {
	CpuGovernorName string
}

const (
	dbusIface = "com.github.kamilsamaj.CpuGovernor"
	dbusPath  = "/com/github/kamilsamaj/CpuGovernor"
)

const introSpec = `<node>
	<interface name="com.github.kamilsamaj.CpuGovernor">
		<method name="SetGovernor">
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `
</node>`

func (f CpuGovernor) SetGovernor(govName string) (string, *dbus.Error) {
	log.Println("Selected governor:", govName)
	if govName == "performance" {
		err := setPerformanceGovernor()
		if err != nil {
			return "", dbus.NewError("org.freedesktop.DBus.Properties.Error", []interface{}{err.Error()})
		}
	} else if govName == "powersave" {
		err := setPowersaveGovernor()
		if err != nil {
			return "", dbus.NewError("org.freedesktop.DBus.Properties.Error", []interface{}{err.Error()})
		}
	} else {
		log.Printf("Unknown governor '%s'\n", govName)
	}
	f.CpuGovernorName = govName
	return govName, nil
}

func setPowersaveGovernor() error {
	cpu := intelcpu.New()
	cores, _ := cpu.GetCores()

	for _, core := range cores {
		err := core.SetGovernor(intelcpu.CPUGovernorPowersave)
		if err != nil {
			return fmt.Errorf("cannot set cpu powersave governor: %w", err)
		}
		err = core.SetPreference(intelcpu.CPUPreferencePower)
		if err != nil {
			return fmt.Errorf("cannot set cpu power preference: %w", err)
		}
	}
	return nil
}

func setPerformanceGovernor() error {
	cpu := intelcpu.New()
	cores, _ := cpu.GetCores()

	for _, core := range cores {
		err := core.SetGovernor(intelcpu.CPUGovernorPerformance)
		if err != nil {
			return fmt.Errorf("cannot set cpu performance governor: %w", err)
		}
		err = core.SetPreference(intelcpu.CPUPreferencePerformance)
		if err != nil {
			return fmt.Errorf("cannot set cpu performance preference: %w", err)
		}
	}
	return nil
}

func init() {
	// omit timestamp from the log output
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

func main() {
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cpuGov := CpuGovernor{
		CpuGovernorName: "powersave",
	}

	conn.Export(cpuGov, dbusPath, dbusIface)
	conn.Export(introspect.Introspectable(introSpec), dbusPath,
		"org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName(dbusIface,
		dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Fatalln(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, dbusIface, "name already taken")
		os.Exit(1)
	}
	log.Println("Setting up CPU governor:", cpuGov.CpuGovernorName)
	log.Printf("Listening on D-Bus interface %s, D-Bus Path: %s ...\n", dbusIface, dbusPath)
	select {}
}
