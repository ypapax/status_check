This is a service for monitoring availability of a particular set of services.
This is an example of a service list:
- google.com
- youtube.com
- facebook.com
- baidu.com
- wikipedia.org
- yahoo.com
- tmall.com
- amazon.com
- twitter.com
- live.com
- instagram.com

Every minute it checks availability and response time of a service. It stores results and gives 
reports through JSON API.

Service is unavailable if it gives 502, 503 or 504 status code.

Reports are the following:
- Amount of available and not available services for a given period of time.
- Amount of available services with a response time greater than 1 second for a given time. 
- Amount of available services with a response time less than 1 second for a given time.

# Running in docker compose
`./commands.sh runc`
# Running tests
`./commands.sh test`

```
--- PASS: TestApi (298.97s)
    --- PASS: TestApi/simple/services-count/available/1565073066/1565080266 (0.01s)
        main_test.go:106: requesting  http://localhost:3001/services-count/available/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 1}
    --- PASS: TestApi/simple/services-count/not-available/1565073066/1565080266 (0.00s)
        main_test.go:106: requesting  http://localhost:3001/services-count/not-available/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 1}
    --- PASS: TestApi/simple/services-count/faster/1000/1565073066/1565080266 (0.00s)
        main_test.go:106: requesting  http://localhost:3001/services-count/faster/1000/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 2}
    --- PASS: TestApi/simple/services-count/slower/1000/1565073066/1565080266 (0.00s)
        main_test.go:106: requesting  http://localhost:3001/services-count/slower/1000/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 0}
    --- PASS: TestApi/diff_status/services-count/available/1565073066/1565080266 (0.01s)
        main_test.go:106: requesting  http://localhost:3001/services-count/available/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 1}
    --- PASS: TestApi/diff_status/services-count/not-available/1565073066/1565080266 (0.00s)
        main_test.go:106: requesting  http://localhost:3001/services-count/not-available/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 1}
    --- PASS: TestApi/diff_status/services-count/faster/1000/1565073066/1565080266 (0.00s)
        main_test.go:106: requesting  http://localhost:3001/services-count/faster/1000/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 2}
    --- PASS: TestApi/diff_status/services-count/slower/1000/1565073066/1565080266 (0.00s)
        main_test.go:106: requesting  http://localhost:3001/services-count/slower/1000/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 0}
    --- PASS: TestApi/big/services-count/available/1565073066/1565080266 (0.01s)
        main_test.go:106: requesting  http://localhost:3001/services-count/available/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 1021}
    --- PASS: TestApi/big/services-count/not-available/1565073066/1565080266 (0.00s)
        main_test.go:106: requesting  http://localhost:3001/services-count/not-available/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 13}
    --- PASS: TestApi/big/services-count/faster/1000/1565073066/1565080266 (0.00s)
        main_test.go:106: requesting  http://localhost:3001/services-count/faster/1000/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 1014}
    --- PASS: TestApi/big/services-count/slower/1000/1565073066/1565080266 (0.00s)
        main_test.go:106: requesting  http://localhost:3001/services-count/slower/1000/1565073066/1565080266
        api_test.go:147: resp:  {"Count": 20}
PASS
ok  	github.com/ypapax/status_check/test	298.979s
```