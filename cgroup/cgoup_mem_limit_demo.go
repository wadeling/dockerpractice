package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	fmt.Printf("os args :%v",os.Args[0])
	fmt.Println()
	if os.Args[0] == "/proc/self/exe" {
		fmt.Printf("current id %d",syscall.Getpid())
		fmt.Println()

		cmd := exec.Command("sh","-c",`stress --vm-bytes 99m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{

		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Printf("stress run ok")
		fmt.Println()
	}

	fmt.Printf("ready exec /proc/self/exe,process id %v",syscall.Getpid())
	fmt.Println()
	cmd := exec .Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start();err != nil {
		fmt.Println("cmd start err",err)
		os.Exit(1)
	} else {
		fmt.Printf("child process id : %v",cmd.Process.Pid)
		fmt.Println()
		memPath := "testmemorylimit"

		err = os.Mkdir(path.Join(cgroupMemoryHierarchyMount,memPath),0755)
		if err != nil {
			fmt.Printf("mkdir err:%v",err)
			fmt.Println()
		}

		err = ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount,memPath,"tasks"),[]byte(strconv.Itoa(cmd.Process.Pid)),0644)
		if err != nil {
			fmt.Printf("write task err:%v",err)
			fmt.Println()
		}

		err = ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount,memPath,"memory.limit_in_bytes"),[]byte("100m"),0644)
		if err != nil {
			fmt.Printf("write memory limit err:%v",err)
			fmt.Println()
		}
	}
	p,err := cmd.Process.Wait()
	fmt.Printf("process wait:%+v,err:%v",*p,err)
	fmt.Println()
}