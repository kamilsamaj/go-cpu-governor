package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"

	"github.com/gotk3/gotk3/gtk"
	"github.com/kamilsamaj/go-appindicator"
)

const (
	dbusIface = "com.github.kamilsamaj.CpuGovernor"
	dbusPath  = "/com/github/kamilsamaj/CpuGovernor"
)

var (
	powersaveItem          *gtk.MenuItem
	balancePowerItem       *gtk.MenuItem
	balancePerformanceItem *gtk.MenuItem
	performanceItem        *gtk.MenuItem
)

func init() {
	// omit timestamp from the log output
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// TODO: this could be probably less DRY
func setIndicatorItem(indicator *appindicator.Indicator, selectedOption string) error {
	switch selectedOption {
	case "powersave":
		indicator.SetIcon("indicator-cpufreq-25")
		powersaveItem.SetLabel("•  powersave")
		performanceItem.SetLabel("   balance power")
		performanceItem.SetLabel("   balance performance")
		performanceItem.SetLabel("   performance")
		return nil
	case "balancepower":
		indicator.SetIcon("indicator-cpufreq-50")
		powersaveItem.SetLabel("   powersave")
		performanceItem.SetLabel("•  balance power")
		performanceItem.SetLabel("   balance performance")
		performanceItem.SetLabel("   performance")
		return nil
	case "balanceperformance":
		indicator.SetIcon("indicator-cpufreq-75")
		powersaveItem.SetLabel("   powersave")
		performanceItem.SetLabel("   balance power")
		performanceItem.SetLabel("•  balance performance")
		performanceItem.SetLabel("   performance")
		return nil
	case "power":
		indicator.SetIcon("indicator-cpufreq-100")
		powersaveItem.SetLabel("   powersave")
		performanceItem.SetLabel("   balance power")
		performanceItem.SetLabel("   balance performance")
		performanceItem.SetLabel("\u2022  performance")
		return nil
	default:
		return fmt.Errorf("unrecognized %s option", selectedOption)
	}
}

func getIndicatorItemFunc(indicator *appindicator.Indicator, dBusConn *dbus.Conn, modeName, iconName string) func() {
	f := func() {
		sanitizedModeName := strings.ReplaceAll(modeName, " ", "")
		log.Printf("requesting %s mode\n", sanitizedModeName)
		resp, err := makeDbusCall(dBusConn, fmt.Sprintf("%s.SetMode", dbusIface), sanitizedModeName)
		if err != nil {
			log.Fatalf("cannot call %s.SetMode: %s", dbusIface, err.Error())
		}
		if resp == sanitizedModeName {
			setIndicatorItem(indicator, sanitizedModeName)
			log.Printf("%s mode set successfully\n", sanitizedModeName)
		} else {
			log.Fatalf("cannot set powersave mode, D-Bus call response: %s\n", resp)
		}
	}
	return f
}

func main() {
	// connect to the system D-Bus
	dBusConn, err := dbus.ConnectSystemBus()
	if err != nil {
		log.Fatalln("failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer dBusConn.Close()

	// initialize the Gtk3 AppIndicator
	gtk.Init(nil)

	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal(err)
	}

	indicator := appindicator.New("cpu-governor", "indicator-cpufreq-25", appindicator.CategoryApplicationStatus)
	indicator.SetTitle("CPU Governor")
	indicator.SetLabel("", "")
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	// add "powersave" item
	powersaveItem, err = gtk.MenuItemNewWithLabel("•  powersave")
	if err != nil {
		log.Fatal(err)
	}
	powersaveItem.Connect("activate",
		getIndicatorItemFunc(indicator, dBusConn, "powersave", "indicator-cpufreq-25"))
	menu.Add(powersaveItem)

	// add "balance power" item
	balancePowerItem, err = gtk.MenuItemNewWithLabel("   balance power")
	if err != nil {
		log.Fatal(err)
	}
	balancePowerItem.Connect("activate",
		getIndicatorItemFunc(indicator, dBusConn, "balance power", "indicator-cpufreq-50"))
	menu.Add(balancePowerItem)

	// add "balance performance" item
	balancePerformanceItem, err = gtk.MenuItemNewWithLabel("   balance performance")
	if err != nil {
		log.Fatal(err)
	}
	balancePerformanceItem.Connect("activate",
		getIndicatorItemFunc(indicator, dBusConn, "balance performance", "indicator-cpufreq-75"))
	menu.Add(balancePerformanceItem)

	performanceItem, err = gtk.MenuItemNewWithLabel("   performance")
	if err != nil {
		log.Fatal(err)
	}
	performanceItem.Connect("activate",
		getIndicatorItemFunc(indicator, dBusConn, "performance", "indicator-cpufreq-100"))
	menu.Add(performanceItem)

	// display the menu
	menu.ShowAll()
	gtk.Main()
}

func makeDbusCall(conn *dbus.Conn, method string, args ...interface{}) (string, error) {
	obj := conn.Object(dbusIface, dbusPath)
	var dBusResp string

	err := obj.Call(method, 0, args...).Store(&dBusResp)
	if err != nil {
		return "", fmt.Errorf("D-Bus call failed: method: %s, args: %v; %w", method, args, err)
	}
	return dBusResp, nil
}
