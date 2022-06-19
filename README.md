## Entain BE Technical Test

# Task 1
Added additional filter `race_visibility`. This filter will allow user/operator to call on ListRaces to filter out the output to race's visible parameter.

The changes necessary were to modify the `ListRacesRequestFilter` to accept an additional parameter to parse in another filter named `race_visibility` from both api and racing servers.

The modification to `applyFilter()` was necessary to query the database to check if there were any additional filter inputs for RaceVisibility.

# Task 1 - Reflection
I haven't used the protobuf before and I failed to understand immediately why the server was not accepting a single input. I had to re-read the .pb.go files to understand the functions generated by the protoc to understand it was accepting an array of a certain type.