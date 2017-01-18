#Suppression

To build execute `build` script.

To run use `./suppress [-v] -input=<input file> -filters="<filter file1> [filter file2] [filter file3] ..."`

###Suppression is implemented in pipeline-like. The process executes in three steps:
- filters preprocessing. Read filters provided and parse it into internal structure.
- process input. Read input provided, classify into bad, stream to saver using filter from step #1.
- save the results and report current progress (# of matches/non-matches for logging).

####1) Process and prepare filter structure combined from "-filters" provided
* for each filter file:
	* stream lines to channel (so-called "fan-out")
* run workers to classify (using md5-hash regex) lines into email/md5 and stream result to channels
* collect data from [emails, md5s] channels to fill filter structure with ("fan-in")	

Step #1 pipeline view:
```
					 routine [email/md5 split worker]
								   w#1
								 /     \
stream filter to channel [lines] - w#2 - channel [md5, email] -> filter
								   ...
								 \     /				
								   w#n
```					  								   
####2) Process input using filter created 
* stream each line from input to `lines` channel
* run workers to `lines`
	* take `line` from channel
		* if `line` is not valid mail (using email regex)
			* send to `bad` channel
		* if `line` exists in `filter[emails]`
			* send to `matches` channel
		* if `line.md5` in `filter[md5]`
			* send to `matches` channel
		* else 
			* send to `clean` channel
		
####3) Run saver routine that starts workers and collect progress report data
* start 3 saver workers to write data from `matches`, `clean` and `bad` channels to corresponding files
* collect status from workers and report feedback on timer (set to 1 minute) tick to `main` routine

Steps #2-3 pipeline view:
```
		 routine [matches/clean/bad classifier]
						  c#1											 channel[matches] -> write to .matches
						/     \										   /						
stream input to [lines] - c#2 - channel [matches, clean, bad] -> saver - channel[clean]	  -> write to .clean
					      ...										   \						 
						\     /										     channel[bad]	  -> write to .bad
				     	  c#n	
```								  
All channels are buffered.
								
                
####Possible improvements:
- if few filter files provided, process them concurrently (depends on common use-case).
