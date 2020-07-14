package main

import(
        "fmt"
        "os/exec"
        "os"
        "bufio"
        "log"
        "time"
	"regexp"
)

var (
	minutes int = 120
	total_count int
	deployed int
	result string = ""
)

func check_status() {
	cmd := exec.Command("juju", "status")
        stdout, err := cmd.Output()
        check_error(err)

        file := "status.txt"
        remove_previous(file)
        recent_status, err_open := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
        check_error(err_open)
        defer recent_status.Close()
	_, err_output := fmt.Fprintln(recent_status, string(stdout))
        check_error(err_output)

}

func verify_deployment() (int,string){
        tries := 0
        status_verified := false
        for (status_verified==false && tries<=minutes){
                tries++
		total_count = 0
		deployed = 0
		time.Sleep(60*time.Second) //sleep for a minute before checking the juju status
		check_status()//check status after a minute
		fin, err_open_file := os.Open("status.txt")
                check_error(err_open_file)
                defer fin.Close()
                scanner := bufio.NewScanner(fin)
                for scanner.Scan() {
			total_count += 1
                        line := scanner.Text()
			in_progress := regexp.MustCompile(`allocating|pending|waiting|blocked|executing|maintenance`)
			in_error := regexp.MustCompile(`error`)
                        if in_error.MatchString(line)==true {
				return -1, "Error, "
			}
			if in_progress.MatchString(line)==false {
				deployed += 1
                        }
                }
		if total_count == deployed {
			status_verified = true
		}
        }
        if !status_verified{
                return -1,"Tries Expired, "
        } else {
                return 0,""
        }
}

func check_error(e error) {
    if e != nil {
        log.Fatal("Error- ", e)
    }
}

func remove_previous(filename string) {
	if file_exists(filename) {
		err_file := os.Remove(filename)
                check_error(err_file)
        }
}

func file_exists(fileName string) (bool) {
    _, err := os.Stat(fileName)
    if os.IsNotExist(err) {
        return false
    } else {
            return true
    }
}

func deploy_script() {
	cmd := exec.Command("/bin/sh","-x", "deploy-contrail.sh", os.Args[1], os.Args[2])
        err := cmd.Run()
        check_error(err)
}

func juju_deployment(){
	start := time.Now()
	deploy_script()
        return_code,return_message := verify_deployment()
	end := time.Now()
	elapsed_time := end.Sub(start)

        if return_code == -1{
                return_message += "Failed deployment\n"
        } else {
                return_message += "Successfully deployed\n"
        }
	result += "\nJuju deployment - " + os.Args[1] + " " + os.Args[2] + "\n"
	result += return_message
	result += "\nStarted at " + start.String()
	result += "\nEnded at " + end.String()
	result += "\nTime taken = " + elapsed_time.String()
}

func write_result() {
	file := "result.txt"
        remove_previous(file)
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
