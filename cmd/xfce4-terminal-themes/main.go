package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/pflag"
	"gopkg.in/ini.v1"
	"os"
	"path"
	"sort"
	"strings"
)

var (
	// Version of this tool
	Version        = "0.1.0"
	configFileName = "terminalrc"
	themesFileName = "themes"
)

func filePathFor(name string) (filePath string) {
	configDir := os.Getenv("XDG_HOME")
	if configDir == "" {
		configDir = path.Join(os.Getenv("HOME"), ".config")
	}

	filePath = path.Join(configDir, "xfce4", "terminal", name)

	return
}

func readConfig(name string) (file *ini.File, err error) {
	filePath := filePathFor(name)
	file, err = ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, filePath)

	return
}

func themeNames(themes *ini.File) (list []string) {
	for _, t := range themes.SectionStrings() {
		if t != "DEFAULT" {
			list = append(list, t)
		}
	}

	sort.Strings(list)

	return
}

func currentTheme(config *ini.Section) (theme string) {
	var b bytes.Buffer

	themeName := config.Key("ThemeName").String()
	fontName := config.Key("FontName").String()

	b.WriteString(fmt.Sprintf("Theme name: %s\n", themeName))
	b.WriteString(fmt.Sprintf("Font name: %s\n", fontName))

	theme = b.String()

	return
}

func setTheme(config *ini.File, themes *ini.File, themeName string) {
	c := config.Section("Configuration")

	for _, k := range themes.Section(themeName).Keys() {
		c.Key(k.Name()).SetValue(k.String())
	}

	config.SaveTo(filePathFor(configFileName))
}

func main() {
	config, err := readConfig(configFileName)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	themes, err := readConfig(themesFileName)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	listThemes := pflag.BoolP("themes", "l", false, "List theme names")
	current := pflag.BoolP("current", "c", false, "Display current theme")
	version := pflag.BoolP("version", "V", false, "Show version")
	pflag.BoolP("help", "h", false, "Show help")

	pflag.Parse()
	args := pflag.Args()

	switch {
	case *listThemes:
		fmt.Println(strings.Join(themeNames(themes), "\n"))
		os.Exit(0)

	case *current:
		section := config.Section("Configuration")
		fmt.Println(currentTheme(section))
		os.Exit(0)

	case *version:
		fmt.Printf("%s %s\n", os.Args[0], Version)
		os.Exit(0)

	case len(args) > 0:
		themeName := strings.Join(args, " ")
		setTheme(config, themes, themeName)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS|THEME NAME]\n", os.Args[0])
		pflag.PrintDefaults()

		os.Exit(0)
	}
}
