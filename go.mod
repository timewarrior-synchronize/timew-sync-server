module git.rwth-aachen.de/computer-aided-synthetic-biology/bachelorpraktika/2020-67-timewarrior-sync/timew-sync-server

go 1.15

replace timewsync/sync => ./sync

replace timewsync => ./

require (
	timewsync v0.0.0-00010101000000-000000000000 // indirect
	timewsync/sync v0.0.0-00010101000000-000000000000
)
