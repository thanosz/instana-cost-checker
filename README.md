Check the amount of data ingested by Instana server and produce a warning if exceeds the specified threshold; returns 1 if threshold exceeded, 0 otherwise.

 
## Options:

	  -endpoint string
	    	The endpoint to connect to (e.g. unit-tenant.instana.io, required)
	  -maxallowed string
	    	The maximum entitled data usage in MB, GB or TB (e.g. 7TB, required)
	  -month int
	    	The month of the year to request data for (optional, skip for current month) (default 12)
	  -threshold float
	    	The percentage to multiply with to generate a warning (optional) (default 0.8)
	  -token string
	    	The authentication token to use (required)
	  -verbose
	    	Verbose output for each day
	  -year int
	    	The year (optional, skip for current year) (default 2024)


## Example: 
Check current month usage and produce a warning if the total ingested data is at 70% of the allowed ingested data. 

	instana-cost-checker -token TOKEN -endpoint unit-tenant.instana.io -maxallowed 7TB -threshold 0.7 -verbose
		
	2024-12-23 02:00:00 +0200 EET
		(infra), 95 MB
		(trace), 337 GB
	2024-12-24 02:00:00 +0200 EET
		(infra), 258 MB
		(trace), 758 GB
	2024-12-25 02:00:00 +0200 EET
		(infra), 258 MB
		(trace), 739 GB
	2024-12-26 02:00:00 +0200 EET
		(infra), 258 MB
		(trace), 802 GB
	2024-12-27 02:00:00 +0200 EET
		(infra), 258 MB
		(trace), 795 GB
	2024-12-28 02:00:00 +0200 EET
		(infra), 258 MB
		(trace), 816 GB
	2024-12-29 02:00:00 +0200 EET
		(infra), 258 MB
		(trace), 826 GB
	2024-12-30 02:00:00 +0200 EET
		(infra), 258 MB
		(trace), 790 GB
	
	Totals:
	   infra: 1.9 GB
	  traces: 5.9 TB
	
	Total Usage for month December: 5.9 TB
	
	Threshold warning!
	exit status 1