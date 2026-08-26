package main

import (
	"asf/config"
	"asf/desktop"
	"asf/dodatki"
	"asf/hardware"
	"asf/osinfo"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	ColorLabelLight string
	ColorLabelDark  string
	ColorSepLight   string
	ColorSepDark    string
	ColorValueLight string
	ColorValueDark  string
	ColorReset      = "\033[0m"
)

func init() {
	if isTrueColor() {
		ColorLabelLight = "\033[38;2;129;161;193m"
		ColorLabelDark = "\033[38;2;136;192;208m"
		ColorSepLight = "\033[38;2;76;86;106m"
		ColorSepDark = "\033[38;2;67;76;94m"

		ColorValueLight = "\033[38;2;187;199;209m"
		ColorValueDark = "\033[38;2;165;177;188m"
	} else {
		ColorLabelLight = "\033[36m"
		ColorLabelDark = "\033[96m"
		ColorSepLight = "\033[90m"
		ColorSepDark = "\033[37m"
		ColorValueLight = "\033[37m"
		ColorValueDark = "\033[37m"
	}
}

func isTrueColor() bool {
	ct := strings.ToLower(os.Getenv("COLORTERM"))
	if strings.Contains(ct, "truecolor") || strings.Contains(ct, "24bit") {
		return true
	}
	return false
}

func main() {
	allFlag := flag.Bool("all", false, "Enable all options")
	logooffFlag := flag.Bool("nologo", false, "Disable logo")
	onlylogoFlag := flag.Bool("onlylogo", false, "Display only logo")
	noneFlag := flag.Bool("none", false, "Disable logo and colors")
	ttsFlag := flag.Bool("tts", false, "That's the Spirit umbrella logo")
	defaultConfigFlag := flag.Bool("default", false, "Use Default Config")
	infoConfigFlag := flag.Bool("info", false, "Info String")
	KFlag := flag.Bool("sharp", false, "Logo is a Croc")

	flag.Parse()
	cfg := config.LoadConfig()
	if *infoConfigFlag {
		cfg.EnableLogo = true
		cfg.EnableUserHost = false
		cfg.EnableOSInfo = false
		cfg.EnableKernel = false
		cfg.EnablePackages = false
		cfg.EnableDEWM = false
		cfg.EnableCPU = false
		cfg.EnableGPU = false
		cfg.EnableRAM = false
		cfg.EnableSwap = false
		cfg.EnableMusic = false
		cfg.EnableUptime = false
		cfg.EnableGTKTheme = false
		cfg.EnableIconTheme = false
		cfg.EnableFont = false
		cfg.EnableShell = false
		cfg.EnableBattery = false
		cfg.EnableDiskInfo = false
		cfg.EnableNetworkInfo = false
		cfg.EnableScreenResolution = false
		fmt.Println("asfetch (asf) - quenðin7 2026")
		os.Exit(0)
	}
	if *allFlag {
		cfg = config.GetAllEnabledConfig()
	}

	if *defaultConfigFlag {
		cfg = config.GetDefaultConfig()
	}

	useSpecialLogo := false
	var specialLogoLines []string

	if *ttsFlag {
		useSpecialLogo = true
		specialLogoLines = []string{
			"                       @+                       ",
			"                       @@                       ",
			"              :@@@@@@@@@@@@@@@@@@               ",
			"           @@@%   @@@      @@#   @@@#           ",
			"         @@=     @@          @@     @@@         ",
			"       @@       @@            @@      @@@       ",
			"     %@@        @             @@        @@      ",
			"    @@         @@              @@        %@=    ",
			"   @@          @@              @@         :@+   ",
			"  @@           @@              @@          *@   ",
			"  @.%@@@@@@@:  @@   @@@@@@@@   @@  @@@@@@@@ @@  ",
			" @@@#       @@@@@@@@   @@   @@ @@@@+       @@@ @ ",
			"*@@           @@@@     @@    :@@@@           @@ ",
			"@@             @@      @@     %@@             @@",
			"@* @@      @@      @.             @@",
			"        @              @@               @       ",
			"       @@              @@               @@      ",
			"        @              @@               @.      ",
			"                 @@    @@        @              ",
			"                 @@    @@        @@             ",
			"                       @@       %@* ",
		}
	} else if *KFlag {
		useSpecialLogo = true
		specialLogoLines = []string{
			"                     .--.  .--.",
			"                    /    \\/    \\",
			"                   | .-.  .-.   \\",
			"                   |/_  |/_  |   \\",
			"                   || `\\|| `\\|    `----.",
			"                   |\\0_/ \\0_/    --,    \\_",
			" .--\"\"\"\"\"-.       /              (` \\     `-.",
			"/          \\-----'-.              \\          \\",
			"\\  () ()                         /`\\          \\",
			"|                         .___.-'   |          \\",
			"\\                        /` \\|      /           ;",
			" `-.___             ___.' .-.`.---.|             \\",
			"    \\| ``-..___,.-'`\\| / /   /     |              `\\",
			"     `      \\|      ,`/ /   /   ,  /",
			"             `      |\\ /   /    |\\/",
			"              ,   .'`-;   '     \\/",
			"         ,    |\\-'  .'   ,   .-'`",
			"       .-|\\--;`` .-'     |\\.'",
			"      ( `\"'-.|\\ (___,.--'`'",
			"       `-.    `\"`          _.--'",
			"          `.          _.-'`-.",
			"            `''---''`` ",
		}
	}

	infoPairs := []struct {
		Label string
		Value string
	}{}
	if *noneFlag {
		cfg.EnableLogo = false
		cfg.EnableUserHost = true
		cfg.EnableOSInfo = true
		cfg.EnableKernel = true
		cfg.EnablePackages = true
		cfg.EnableDEWM = true
		cfg.EnableCPU = true
		cfg.EnableGPU = true
		cfg.EnableRAM = true
		cfg.EnableSwap = true
		cfg.EnableMusic = true
		cfg.EnableUptime = true
		cfg.EnableGTKTheme = false
		cfg.EnableIconTheme = false
		cfg.EnableFont = false
		cfg.EnableShell = false
		cfg.EnableBattery = true
		cfg.EnableDiskInfo = true
		cfg.EnableNetworkInfo = true
		cfg.EnableScreenResolution = false
		ColorLabelLight = ""
		ColorLabelDark = ""
		ColorSepLight = ""
		ColorSepDark = ""
		ColorValueLight = ""
		ColorValueDark = ""
		ColorReset = ""
	}
	if *logooffFlag {
		cfg.EnableLogo = false
	}
	if *onlylogoFlag {
		cfg.EnableLogo = true
		cfg.EnableUserHost = false
		cfg.EnableOSInfo = false
		cfg.EnableKernel = false
		cfg.EnablePackages = false
		cfg.EnableDEWM = false
		cfg.EnableCPU = false
		cfg.EnableGPU = false
		cfg.EnableRAM = false
		cfg.EnableSwap = false
		cfg.EnableMusic = false
		cfg.EnableUptime = false
		cfg.EnableGTKTheme = false
		cfg.EnableIconTheme = false
		cfg.EnableFont = false
		cfg.EnableShell = false
		cfg.EnableBattery = false
		cfg.EnableDiskInfo = false
		cfg.EnableNetworkInfo = false
		cfg.EnableScreenResolution = false
	}
	if *ttsFlag {
		cfg.EnableLogo = true
	}

	if cfg.EnableUserHost {
		username, hostname := dodatki.GetUserAndHost()
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"User ", username + "@" + hostname})
	}

	if cfg.EnableOSInfo {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"OS ", osinfo.GetOSInfo()})
	}
	if cfg.EnableKernel {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"Kernel ", osinfo.GetKernel()})
	}
	if cfg.EnablePackages {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"Packages ", osinfo.GetPackageCount()})
	}

	if cfg.EnableDEWM {
		de, wm := desktop.GetDEWM()
		if de != "" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"DE ", de})
		}
		if wm != "" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"WM ", wm})
		}
	}

	if cfg.EnableGTKTheme {
		gtkTheme := desktop.GetGTKTheme()
		if gtkTheme != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"GTK ", gtkTheme})
		}
	}
	if cfg.EnableIconTheme {
		iconTheme := desktop.GetIconTheme()
		if iconTheme != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Icons ", iconTheme})
		}
	}
	if cfg.EnableFont {
		font := desktop.GetFont()
		if font != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Font ", font})
		}
	}
	if cfg.EnableShell {
		shell := dodatki.GetShell()
		if shell != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Shell ", shell})
		}
	}
	if cfg.EnableUptime {
		uptime := dodatki.GetUptime()
		if uptime != "unknown" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Uptime ", uptime})
		}
	}
	if cfg.EnableBattery {
		batteryInfo := hardware.GetBatteryInfo()
		if batteryInfo != "N/A" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Battery ", batteryInfo})
		}
	}
	if cfg.EnableDiskInfo {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"/  ", hardware.GetDiskInfo()})
	}
	if cfg.EnableNetworkInfo {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"Network ", hardware.GetNetworkInfo()})
	}
	if cfg.EnableScreenResolution {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"Resolution ", desktop.GetScreenResolution()})
	}
	if cfg.EnableCPU {
		infoPairs = append(infoPairs, struct {
			Label string
			Value string
		}{"CPU ", hardware.GetCPUInfo()})
	}
	if cfg.EnableGPU {
		gpuInfo := hardware.GetGPUInfo()
		if gpuInfo != "Unknown GPU" && gpuInfo != "N/A" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"GPU ", gpuInfo})
		}
	}

	if cfg.EnableRAM || cfg.EnableSwap {
		ramInfo, swapInfo := hardware.GetMemoryAndSwapInfo()
		if cfg.EnableRAM {
			if ramInfo != "unknown" {
				infoPairs = append(infoPairs, struct {
					Label string
					Value string
				}{"RAM ", ramInfo})
			}
		}
		if cfg.EnableSwap {
			if swapInfo != "unknown" && swapInfo != "No Swap" {
				infoPairs = append(infoPairs, struct {
					Label string
					Value string
				}{"Swap ", swapInfo})
			}
		}
	}

	if cfg.EnableMusic {
		musicInfo := dodatki.GetCurrentMusic()
		if musicInfo != "Not playing" {
			infoPairs = append(infoPairs, struct {
				Label string
				Value string
			}{"Spotify ", musicInfo})
		}
	}

	maxLabelLen := 0
	for _, pair := range infoPairs {
		if len(pair.Label) > maxLabelLen {
			maxLabelLen = len(pair.Label)
		}
	}

	var infoLines []string
	for i, pair := range infoPairs {
		var labelColor, sepColor, valueColor string
		if i%2 == 0 {
			labelColor = ColorLabelLight
			sepColor = ColorSepDark
			valueColor = ColorValueLight
		} else {
			labelColor = ColorLabelDark
			sepColor = ColorSepLight
			valueColor = ColorValueDark
		}
		alignedLabel := fmt.Sprintf("%s%-*s%s", labelColor, maxLabelLen, pair.Label, ColorReset)
		separator := fmt.Sprintf("%s│%s", sepColor, ColorReset)
		value := fmt.Sprintf("%s%s%s", valueColor, pair.Value, ColorReset)
		infoLines = append(infoLines, fmt.Sprintf("%s%s %s", alignedLabel, separator, value))
	}

	if !*onlylogoFlag {
		infoLines = append(infoLines, "")

		row1 := ""
		for i := 0; i < 8; i++ {
			row1 += fmt.Sprintf("\033[4%dm   ", i)
		}
		infoLines = append(infoLines, row1+ColorReset)

		row2 := ""
		for i := 0; i < 8; i++ {
			row2 += fmt.Sprintf("\033[10%dm   ", i)
		}
		infoLines = append(infoLines, row2+ColorReset)
	}

	var logo []string
	if useSpecialLogo {
		logo = specialLogoLines
	} else if cfg.EnableLogo {
		logo, _ = config.LoadLogoFromFile(cfg.LogoPath)
	}

	maxLogoWidth := 0
	for _, line := range logo {
		w := 0
		for _, r := range line {
			if r >= 0x2800 && r <= 0x28FF {
				w += 1
			} else {
				w += 1
			} // Korekta Braille: 1 zamiast 2
		}
		if w > maxLogoWidth {
			maxLogoWidth = w
		}
	}

	totalH := len(logo)
	if len(infoLines) > totalH {
		totalH = len(infoLines)
	}

	fmt.Println()
	for i := 0; i < totalH; i++ {
		lLine, iLine := "", ""
		if i < len(logo) {
			lLine = logo[i]
		}
		if i < len(infoLines) {
			iLine = infoLines[i]
		}

		curW := 0
		for _, r := range lLine {
			if r >= 0x2800 && r <= 0x28FF {
				curW += 1
			} else {
				curW += 1
			}
		}

		if cfg.EnableLogo {
			fmt.Printf("%s%s%s%s%s\n", ColorValueLight, lLine, ColorReset, strings.Repeat(" ", maxLogoWidth-curW+4), iLine)
		} else {
			fmt.Printf("%s%s\n", strings.Repeat(" ", 4), iLine)
		}
	}
	fmt.Println()
}
