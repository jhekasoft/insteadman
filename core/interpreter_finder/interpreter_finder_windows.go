// +build windows

package interpreterFinder

import (
	"syscall"
)

func exactFilePaths() []string {
	paths := []string{}

	for _, drive := range getDrives() {
		drivePaths := []string{
			drive + ":\\Program Files (x86)\\Games\\INSTEAD\\sdl-instead.exe",
			drive + ":\\Program Files\\Games\\INSTEAD\\sdl-instead.exe",
			drive + ":\\Program Files (x86)\\INSTEAD\\sdl-instead.exe",
			drive + ":\\Program Files\\INSTEAD\\sdl-instead.exe",
		}
		paths = append(paths, drivePaths...)
	}

	return paths
}

// https://stackoverflow.com/a/23135463
func getDrives() []string {
	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	getLogicalDrivesHandle, _ := syscall.GetProcAddress(kernel32, "GetLogicalDrives")

	ret, _, callErr := syscall.Syscall(uintptr(getLogicalDrivesHandle), 0, 0, 0, 0)
	if callErr != 0 {
		return []string{}
	}

	return bitsToDrives(uint32(ret))
}

func bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i])
		}
		bitMap >>= 1
	}

	return drives
}
