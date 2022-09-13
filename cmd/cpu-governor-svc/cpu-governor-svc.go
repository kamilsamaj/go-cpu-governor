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
		<method name="SetMode">
			<arg direction="in" type="s"/>
			<arg direction="out" type="s"/>
		</method>
	</interface>` + introspect.IntrospectDataString + `
</node>`

func (f CpuGovernor) SetMode(modeName string) (string, *dbus.Error) {
	log.Println("Selected mode:", modeName)
	var modeSetterFunc func() error
	if modeName == "powersave" {
		modeSetterFunc = getPerfSetterFunc(intelcpu.CPUGovernorPowersave, intelcpu.CPUPreferencePower)
	} else if modeName == "balancepower" {
		modeSetterFunc = getPerfSetterFunc(intelcpu.CPUGovernorPowersave, intelcpu.CPUPreferenceBalancePower)
	} else if modeName == "balanceperformance" {
		modeSetterFunc = getPerfSetterFunc(intelcpu.CPUGovernorPowersave, intelcpu.CPUPreferenceBalancePerformance)
	} else if modeName == "performance" {
		modeSetterFunc = getPerfSetterFunc(intelcpu.CPUGovernorPerformance, intelcpu.CPUPreferencePerformance)
	} else {
		log.Printf("Unknown governor '%s'\n", modeName)
		return "", dbus.NewError("org.freedesktop.DBus.Properties.Error",
			[]interface{}{fmt.Errorf("unknown governor '%s'", modeName).Error()})
	}

	err := modeSetterFunc()
	if err != nil {
		return "", dbus.NewError("org.freedesktop.DBus.Properties.Error", []interface{}{err.Error()})
	}
	f.CpuGovernorName = modeName
	return modeName, nil
}

func getPerfSetterFunc(governor intelcpu.CPUCoreGovernor, perfPreference intelcpu.CPUPreference) func() error {
	f := func() error {
		cpu := intelcpu.New()
		cores, _ := cpu.GetCores()

		for _, core := range cores {
			err := core.SetGovernor(governor)
			if err != nil {
				return fmt.Errorf("cannot set cpu %s governor: %s", governor, err)
			}
			err = core.SetPreference(perfPreference)
			if err != nil {
				return fmt.Errorf("cannot set cpu %s preference: %s", perfPreference, err)
			}
		}
		return nil
	}
	return f
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

	// default mode when the service starts
	cpuGov := CpuGovernor{
		CpuGovernorName: "balancepower",
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
