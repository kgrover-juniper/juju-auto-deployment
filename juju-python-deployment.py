import os
import sys
import time
import subprocess
from datetime import datetime
minutes = 120

def verify_deployment():
    tries = 0
    status_verified = False
    while (status_verified is False and tries<=minutes):
        tries += 1
        time.sleep(60)
        c1 = subprocess.Popen(["juju", "status"], stdout=subprocess.PIPE)
        c2 = subprocess.Popen(["grep", "-e allocating", "-e blocked", "-e pending", "-e waiting", "-e maintenance", "-e executing", "-e error"], 
                stdin=c1.stdout,stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout = c2.communicate()[0]
        if not stdout:
            status_verified = True
        c1.wait()
    if status_verified is False:
        return -1,"Tries Expired, "
    else:
        return 0,""

def juju_deployment():
    global start,end,message
    start = datetime.now()
    deploy_script = subprocess.call(["./deploy-contrail.sh", sys.argv[1], sys.argv[2]], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    return_code, return_message = verify_deployment()
    end = datetime.now()
    if return_code == -1:
        message = return_message + "Failed deployment"
    else:
        message = return_message + "Successfully deployed"


def write_result():
    result = "\nJuju deployment " + str(sys.argv[1]) + " " + str(sys.argv[2]) + "\n" + message
    result += "\nStarted at " + str(start) + "\nEnded at " + str(end) + "\nTime taken = " + str(end-start) + "\n"
    fout = "result.txt"
    if os.path.exists(fout):
        os.remove(fout)
    with open(fout, 'w') as outfile:
        outfile.write(result)

def main():
    juju_deployment()
    write_result()

if __name__ == "__main__":
    main()
