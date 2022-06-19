## Entain BE Technical Test

## Task 1
Added additional filter `race_visibility`. This filter will allow user/operator to call on ListRaces to filter out the output to race's visible parameter.

The changes necessary were to modify the `ListRacesRequestFilter` to accept an additional paremter to parse in another filter named `race_visibility` from both api and racing servers.

The modification to `applyFilter()` was necessary to query the database to check if there were any additional filter inputs for RaceVisibilty.

### Task 1 - Reflection
I haven't used the protobuff before and I failed to understand immediately why the functions couldn't parse single single input (e.g. {filter:{"meeting_ids":1}}). I had to re-read the .pb.go files to understand the functions generated by the proto was accepting an array of input of a certain type. 

## Task 2
Added `raceReportSortByAdvertisedTime()` to the `ListRaces()` so any newly added races can be ordered whenever a message is sent from the user/operator. 

### Task 2 - Reflection
I haven't implemented the bonus task as I usually tend to the MVP first and apply any additional feature/upgrades after discussing with the team and client. I'm working this task as how I would usually work in an environment with a team and assuming I would have to deliver the MVP first then bonus tasks at the end when MVP has been completed.

## Task 3
Added a new message field to `Race` to api and racing server. This allows user to visually see if a certain types of races are CLOSED or OPEN depending on the current time.

The database table is modified on runtime with when the `ListRaces` is called to update the status of all races. Cases that was considered was when a new race is added or just simply, time passing, the listed races needs to be up to date/time. I tried to make the code to account for future implementations or modifications needed to the database and made the functions more generic to add additional columns on runtime and modify the status at the same time.

### Task 3 - Reflection
It's been a while since I handled any database query (sqlite) so having to set this up was quite fun. I faced some difficulty when trying to make sure I was using the correct syntax to modify the selected row of the database but after a few tries (and bunch of online searching) I found out the best way I could implement how to add a new column. Unfortunately, I made a poor implementation to decern if the column of a specific name already exists. I couldn't figure out how to use the `INSER OR IGNORE` or if I could find an alternative to `IF NOT EXIST` command for sqlite3, so I made a dodgy error handling to ignore the 'duplicate column exists' error. In my future attempts, hopefully I'll understand the correct keywords necessary. 

## Task 4
I've implemented to the extent where I wanted to start testing my implementation. However, I noticed I wasn't able to generate the .pb.gw.go file which I believe should be containing the "GET" `RegisterRacingHandlerClient()` from the protoc generate command. I tried to read all the resources I could find, and tried to resolve it, but kept getting the error command below regarding 'protoc-gen-grpc-gateway' binary missing. Much of the search tells me the imports or the binary is a pre-compiled file which isn't well documented and isn't tracked by go (https://github.com/grpc-ecosystem/grpc-gateway/issues/1065#issuecomment-544241612). I'm not too sure how true this is today, so I tried to compile the binaries needed myself but in the end, it didn't work.

The left over task which I couldn't complete would be to generate the protoc .pb.gw.go file and finish off `applyFindRaceById()` and implement it for my interface `GetRace()`.

```
[Info  - 10:05:10 pm] 2022/06/19 22:05:10 background imports cache refresh starting

[Info  - 10:05:10 pm] 2022/06/19 22:05:10 background refresh finished after 3.2328ms

[Info  - 10:05:15 pm] 2022/06/19 22:05:15 protoc -I . --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative --grpc-gateway_out . --grpc-gateway_opt paths=source_relative racing/racing.proto

	operation="generate"

[Info  - 10:05:15 pm] 2022/06/19 22:05:15 'protoc-gen-grpc-gateway' is not recognized as an internal or external command,
operable program or batch file.

	operation="generate"

[Info  - 10:05:15 pm] 2022/06/19 22:05:15 --grpc-gateway_out
	operation="generate"

[Info  - 10:05:15 pm] 2022/06/19 22:05:15 : protoc-gen-grpc-gateway: Plugin failed with status code 1.

	operation="generate"

[Info  - 10:05:15 pm] 2022/06/19 22:05:15 api.go:3: running "protoc": exit status 1

	operation="generate"

[Error - 10:05:15 pm] 2022/06/19 22:05:15 command error: exit status 1

[Error - 10:05:15 pm] Request workspace/executeCommand failed.
  Message: exit status 1
  Code: 0 

```

### Task 4 - Reflection
Unfortuantely I don't have much to say other than my lack of knowledge in this area to complete this task. It was fun and I felt pretty comfortable using the grpc/protoc tool with golang as it handled much of the difficult and tedious tasks already, but I realize setting up the environment properly from the get go is also important to utilize this tool. I noticed the failure of the go:generate in api.go from Task 1 but didn't investigate any further as I was able to build the pb.go and grpc.pb.go from another command and didn't realize the importance of the 'protoc-gen-grpc-gateway' imports.