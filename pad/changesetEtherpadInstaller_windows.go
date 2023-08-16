//go:build windows
// +build windows

package pad

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

type SHELLEXECUTEINFO struct {
	CbSize       uint32
	FMask        uint32
	Hwnd         uintptr
	LpVerb       *uint16
	LpFile       *uint16
	LpParameters *uint16
	LpDirectory  *uint16
	NShow        int32
	HInstApp     uintptr
	// Optional fields
	LpIDList     uintptr
	LpClass      *uint16
	HkeyClass    uintptr
	DwHotKey     uint32
	HIcon        uintptr
	HProcess     uintptr
}

const (
	SEE_MASK_NOCLOSEPROCESS = 0x00000040
	SW_NORMAL               = 1
)

var (
	modshell32           = syscall.NewLazyDLL("shell32.dll")
	procShellExecuteExW  = modshell32.NewProc("ShellExecuteExW")
)

func ShellExecuteEx(lpExecInfo *SHELLEXECUTEINFO) bool {
	ret, _, _ := procShellExecuteExW.Call(uintptr(unsafe.Pointer(lpExecInfo)))
	return ret != 0
}


func installEtherpad(folderPath string) error {

		// Check if the file or directory exists
		path := strings.ReplaceAll(folderPath, "/", `\`) + `\etherpad-lite\node_modules\ep_etherpad-lite`
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			// If it exists, delete it
			err := os.RemoveAll(path)
			if err != nil {
				return fmt.Errorf("error deleting the file ep_etherpad-lite: %v\n", err)
			}
		}



		// Add exit 0 at the end of the install file
		path = strings.ReplaceAll(folderPath, "/", `\`) + `\etherpad-lite\src\bin\installOnWindows.bat`

		// open file using READ & WRITE permission
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("error opening the file installOnWindows.bat: %v\n", err)
		}
		defer file.Close()
	
		// add "exit /B 0" at the end of the file
		_, err = file.WriteString("exit\n")
		if err != nil {
			return fmt.Errorf("error writing to the file installOnWindows.bat: %v\n", err)
		}



		// Run the installOnWindows.bat file as admin
		verb := "runas"
		args := "&& pause"
		dir := strings.ReplaceAll(folderPath, "/", `\`) + `\etherpad-lite`
		
		// Erstellen Sie den absoluten Pfad zur BAT-Datei
		absScriptPath, err := filepath.Abs(dir + `\src\bin\installOnWindows.bat`)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %v", err)
		}
		
		verbPtr, _ := syscall.UTF16PtrFromString(verb)
		scriptPathPtr, _ := syscall.UTF16PtrFromString(absScriptPath)
		argsPtr, _ := syscall.UTF16PtrFromString(args)
		dirPtr, _ := syscall.UTF16PtrFromString(dir)
		
		var showCmd int32 = 1 // SW_NORMAL
		
		var execInfo SHELLEXECUTEINFO
		execInfo.CbSize = uint32(unsafe.Sizeof(execInfo))
		execInfo.FMask = SEE_MASK_NOCLOSEPROCESS
		execInfo.Hwnd = 0
		execInfo.LpVerb = verbPtr
		execInfo.LpFile = scriptPathPtr
		execInfo.LpParameters = argsPtr
		execInfo.LpDirectory = dirPtr
		execInfo.NShow = showCmd
		
		if !ShellExecuteEx(&execInfo) {
			return fmt.Errorf("failed to run script as admin")
		}
		
		// Warten Sie auf den Abschluss des Prozesses
		syscall.WaitForSingleObject(syscall.Handle(execInfo.HProcess), syscall.INFINITE)

		return nil
}