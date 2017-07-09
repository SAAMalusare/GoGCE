
# GoGCE
Google GCE GoLang

Script sverves following purposes:
1. Creating new VM instance on Google compute engine (Displays status, Internal IP & External IP of the created VM)
2. Creating New healthcheck instance on GCE
3. Checking Healthcheck of Given instance

-----------------------------------------------------------------------------------------
Script usage:
-----------------------------------------------------------------------------------------
GCEOps.go -pn=crygce-test -zn=us-central1-a -in=mynewinstance001 -mt=f1-micro -si=centos-cloud/global/images/centos-6-v20170620 -ops=I
------------------------------------------------------------------------------------------
-pn (Type string): PROJECT NAME
-zn (Type string): ZONE, e.g. us-central1-a
-in (Type string): INSTANCE NAME (VM NAME), e.g. mynewvm
-mt (Type string): MACHINE TYPE, e.g. f1-micro
-si (Type string): PATH TO SOURCE IMAGE, e.g. centos-cloud/global/images/centos-6-v20170620
-ops (Type string): TYPE OF OPERATIONS e.g. I for create new instance, 
                    H for create new HealthCheck, 
                    B for creating new VM & HealthCheck instance and
                    default is HealthCheck response.
-hin (Type string): INSTANCE NAME (VM NAME), e.g. mynewhealthcheck
-----------------------------------------------------------------------------------------
Parameter requiredness:
-----------------------------------------------------------------------------------------
if OPS = "I"
Mandatory Parameters: Project name, Zone, Source Image, Instance name, MachineType, Operation
e.g.
GCEOps.go -pn=crygce-test -zn=us-central1-a -in=mynewinstance001 -mt=f1-micro -si=centos-cloud/global/images/centos-6-v20170620 -ops=I

if OPS = "B"
Mandatory Parameters: Project name, Zone, Source Image, Instance name, Healthcheck Instance name, Operation
e.g.
GCEOps.go -pn=crygce-test -zn=us-central1-a -in=mynewinstance001 -mt=f1-micro -si=centos-cloud/global/images/centos-6-v20170620 -hin=healthcheckinstanc -ops=I

if OPS = "H"
Mandatory Parameters: Project name, Zone, Healthcheck Instance name, Operation
e.g.
GCEOps.go -pn=crygce-test -zn=us-central1-a -hin=healthcheckinstanc -ops=I

if OPS = ""
Mandatory Parameters: Healthcheck Instance name
e.g.
GCEOps.go -hin=healthcheckinstance
-----------------------------------------------------------------------------------------
