package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"golang.org/x/sys/windows/registry"
)

var (
	blu        = color.New(color.FgBlue)
	boldBlue   = blu.Add(color.Bold)
	rd         = color.New(color.FgRed)
	boldRed    = rd.Add(color.Bold)
	grn        = color.New(color.FgGreen)
	boldGreen  = grn.Add(color.Bold)
	yel        = color.New(color.FgYellow)
	boldYellow = yel.Add(color.Bold)
	cyn        = color.New(color.FgCyan)
	boldCyan   = cyn.Add(color.Bold)
)

// ------------------------- Helper Functions --------------------------
func print(args ...interface{}) {
	fmt.Println(args...)
}

func check(e error, errMsg string) bool {
	if e != nil {
		print(errMsg)
	}
	return true
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func getNumberOfSubKeysAndValues(k registry.Key) (uint32, uint32) {
	keyInfo, err := k.Stat()
	check(err, "Unable to fetch Stat info from registry object...")
	return keyInfo.SubKeyCount, keyInfo.ValueCount
}

func openKey(hive registry.Key, subkey string, access uint32) registry.Key {
	key, err := registry.OpenKey(hive, subkey, access)
	check(err, "Unable to open registry key...")
	return key
}

// toTime converts an 8-byte Windows Filetime to time.Time.
func toTime(t []byte) time.Time {
	ft := &syscall.Filetime{
		LowDateTime:  binary.LittleEndian.Uint32(t[:4]),
		HighDateTime: binary.LittleEndian.Uint32(t[4:]),
	}
	return time.Unix(0, ft.Nanoseconds())
}

// ------------------- Program Body ---------------------
func getComputerInfo() {
	key := openKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.ALL_ACCESS,
	)
	defer key.Close()

	boldBlue.Println("◎ ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ Computer Build Info ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ ◎")
	print("")

	productName, _, err := key.GetStringValue("ProductName")
	check(err, "ProductName value not found in registry...")
	print("Product Name : " + productName)
	currentVersion, _, err := key.GetStringValue("CurrentVersion")
	check(err, "CurrentVersion value not found in registry...")
	print("Current Version : " + currentVersion)
	currentBuildNumber, _, err := key.GetStringValue("CurrentBuildNumber")
	check(err, "CurrentBuildNumber Value not found in registry...")
	print("Build Number : " + currentBuildNumber)
	registeredOwner, _, err := key.GetStringValue("RegisteredOwner")
	check(err, "RegisteredOwner value not found in registry...")
	print("Registered Owner : " + registeredOwner)
	print("")
}

func getInstalledApps() {
	key := openKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`,
		registry.ALL_ACCESS,
	)
	defer key.Close()

	numOfSubKeys, numOfValues := getNumberOfSubKeysAndValues(key)
	subkeys, err := key.ReadSubKeyNames(int(numOfSubKeys))
	check(err, "Unable to read subkeys...")

	boldRed.Println("◎ ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ Installed Applications ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ ◎")
	print("")
	for _, skey := range subkeys {
		k := openKey(
			registry.LOCAL_MACHINE,
			`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`+"\\"+skey,
			registry.ALL_ACCESS,
		)
		values, err := k.ReadValueNames(int(numOfValues))
		check(err, "Unable to read values from registry key...")
		if exist := find(values, "DisplayName"); exist {
			val, _, err := k.GetStringValue("DisplayName")
			check(err, "Unable to retrieve data from value DisplayName...")
			print("\u2022 " + val)
		} else {
			print("\u2022 " + skey)
		}
	}
}

func getEnVars() {
	key := openKey(
		registry.LOCAL_MACHINE,
		`SYSTEM\CurrentControlSet\Control\Session Manager\Environment`,
		registry.ALL_ACCESS,
	)
	defer key.Close()

	_, numOfValues := getNumberOfSubKeysAndValues(key)
	environmentVariables, err := key.ReadValueNames(int(numOfValues))
	check(err, "Unable to read values from registry key...")

	boldGreen.Println("\n◎ ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ Environment Variables ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ ◎")
	print("")

	for _, envar := range environmentVariables {
		envarValue, _, err := key.GetStringValue(envar)
		check(err, "Unable to retrieve data from value in registry key...")
		print(envar + " ☰☰ " + envarValue)
	}
	print("")
}

func getStartUpApps() {

}

func getJumpLists() {
	currentUser, err := user.Current()
	check(err, "Unable to fetch username...")
	username := strings.Split(currentUser.Username, `\`)
	jumpListPath := fmt.Sprintf(
		`C:\Users\%s\AppData\Roaming\Microsoft\Windows\Recent\AutomaticDestinations`,
		username[1],
	)
	jumpListFiles, err := ioutil.ReadDir(jumpListPath)
	check(err, "Unable to read files in jump list directory...")

	boldYellow.Println("◎ ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ Jump List Files ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ ◎")
	print("")
	for _, file := range jumpListFiles {
		print(file.Name())
	}
	print("")
}

func getLNKFiles() {

}

func getShellBags() {

}

func getPrefetchFiles() {

}

func getRecycleBinFiles() {
	recycleBinPath := `C:\$Recycle.Bin`
	recycleBinFiles, err := ioutil.ReadDir(recycleBinPath)
	check(err, "Unable to open recycle bin folder...")
	currentUser, err := user.Current()
	check(err, "Unable to get user info...")
	userSID := currentUser.Uid

	boldCyan.Println("◎ ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ Recycle Bin Files ☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶☶ ◎")

	for _, recycleFolder := range recycleBinFiles {
		folderName := recycleFolder.Name()
		if folderName == userSID {
			userRecycleBinContents, err := ioutil.ReadDir(recycleBinPath + `\` + folderName)
			check(err, "Unable able to open user's recycle bin folder...")
			for _, recycledFile := range userRecycleBinContents {
				if strings.HasPrefix(recycledFile.Name(), "$I") {
					fi, err := os.Open(recycleBinPath + `\` + folderName + `\` + recycledFile.Name())
					header := make([]byte, 8)
					_, err = fi.Read(header)
					check(err, "Unable to read data from dollar I recycle file...")

					operatingSystem := "Prior to Windows 10"
					if binary.LittleEndian.Uint64(header) == 2 {
						operatingSystem = "Windows 10"
					}

					fileSize := make([]byte, 8)
					_, err = fi.Read(fileSize)
					check(err, "Unable to read data from dollar I recycle file...")

					deletedTimeStamp := make([]byte, 8)
					_, err = fi.Read(deletedTimeStamp)
					check(err, "Unable to read data from dollar I recycle file...")

					dateDeleted := toTime(deletedTimeStamp)
					//convertToDate := fmt.Sprintf(`[DateTime]::FromFileTimeutc("%d")`, timeStamp)
					//date, err := exec.Command("powershell.exe", `-c`, convertToDate).CombinedOutput()
					check(err, "Unable to retrieve time stamp from recycled file...")

					fileNameLength := make([]byte, 4)
					_, _ = fi.Read(fileNameLength)

					dollarIFileSize, _ := fi.Stat()
					fileName := make([]byte, (dollarIFileSize.Size() - 8 - 8 - 8 - 4))
					_, err = fi.Read(fileName)

					print("")
					fmt.Print("File Name: ")
					boldCyan.Println(string(fileName))
					print("OS: " + operatingSystem)
					fmt.Print("File Deleted On: ")
					boldRed.Println(dateDeleted)
					print("File size: " + strconv.Itoa(int(binary.LittleEndian.Uint64(fileSize))))

				}
			}
		}
	}
}

func main() {
	getComputerInfo()
	getInstalledApps()
	getEnVars()
	getStartUpApps()
	getJumpLists()
	getLNKFiles()
	getShellBags()
	getPrefetchFiles()
	getRecycleBinFiles()
}
