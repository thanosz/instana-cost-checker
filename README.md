Check the amount of data ingested by Instana server and produce a warning if exceeds the specified threshold; returns 1 if threshold exceeded, 0 otherwise.

Example: check current month usage and produce a warning if the total ingested data is at 70% of the allowed ingested data.

	instana-cost-checker -token TOKEN -endpoint unit-tenant.instana.io -maxallowed 7TB -threshold 0.7

 
Options:

  -endpoint string\
    	The endpoint to connect to (e.g. unit-tenant.instana.io, required)\
     
  -maxallowed string\
    	The maximum entitled data usage in MB, GB or TB (e.g. 7TB, required)\
     
  -month int\
    	The month of the year to request data for (optional, skip for current month) (default 12)\
     
  -threshold float\
    	The percentage to multiply with to generate a warning (optional) (default 0.8)\

  -token string\
    	The authentication token to use (required)\
     
  -year int\
    	The year (optional, skip for current year) (default 2024)
