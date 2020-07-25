package main

import(
        "fmt"
        "os/exec"
        "os"
        "log"
        "time"
)

var (
	start time.Time
	end time.Time
	minutes int = 120
	message string
)

func check_error(e error) {
    if e != nil {
        log.Fatal("Error- ", e)
    }
}

func verify_deployment() (int, string) {
	tries := 0
        status_verified := false
        for (status_verified==false && tries<=minutes){
                tries++
		time.Sleep(60*time.Second) //sleep for a minute before checking the juju status
                c2 := exec.Command("grep", "-e allocating", "-e blocked", "-e pending", "-e waiting", "-e maintenance", "-e executing", "-e error")
                c1 := exec.Command("juju", "status")
                pipe, _ := c1.StdoutPipe()
                defer pipe.Close()
                c2.Stdin = pipe
                c1.Start()
                stdout, _ := c2.Output()
		if string(stdout) == "" {
			status_verified = true
		}
		_ = c1.Wait()
	}
	if !status_verified{
                return -1,"Tries Expired, "
        } else {
                return 0,""
        }
}

func deploy_script() {
        cmd := exec.Command("/bin/sh","-x", "deploy-contrail.sh", os.Args[1], os.Args[2])
        err := cmd.Run()
        check_error(err)
}

func juju_deployment(){
	start = time.Now()
	deploy_script()
	return_code, return_message := verify_deployment()
	end = time.Now()

        if return_code == -1{
                message = return_message + "Failed deployment"
        } else {
                message = return_message + "Successfully deployed"
        }
}

func write_result() {
	result := "\nJuju deployment - " + os.Args[1] + " " + os.Args[2] + "\n" + message
	result += "\nStarted at " + start.String() + "\nEnded at " + end.String() + "\nTime taken = " + end.Sub(start).String()
	file := "result.txt"
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		err_file := os.Remove(file)
                check_error(err_file)
        }
        output, err_open := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
        check_error(err_open)
        defer output.Close()
	_, err_write := fmt.Fprintln(output, result)
        check_error(err_write)
}

func main() {
	juju_deployment()
	write_result()
}
