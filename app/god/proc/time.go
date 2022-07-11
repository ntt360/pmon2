package proc

import "syscall"

//we're using this file to handle the ktime timestamps that proc connector gives us

//get the time of boot, used to turn the ktime timestamp into a real date
func getBoottime() (int64, error) {

	//this seems to be how dmesg -T does it
	//This is dmesg's 'faillback' method as getting access to clock_gettime() from go is more trouble than it's worth.
	tv := &syscall.Timeval{}
	err := syscall.Gettimeofday(tv)
	if err != nil {
		return 0, err
	}

	info := &syscall.Sysinfo_t{}
	err = syscall.Sysinfo(info)
	if err != nil {
		return 0, err
	}

	return (tv.Sec - info.Uptime), nil
}
