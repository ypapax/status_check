# Running in docker compose
`./commands.sh runc`
# Running tests
`./commands.sh test`

```
--- PASS: TestApi (372.16s)
    --- PASS: TestApi/simple/services-count/available/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 1}
    --- PASS: TestApi/simple/services-count/not-available/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 1}
    --- PASS: TestApi/simple/services-count/faster/1000/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 2}
    --- PASS: TestApi/simple/services-count/slower/1000/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 0}
    --- PASS: TestApi/diff_status/services-count/available/1565031305/1565038505 (0.01s)
        api_test.go:147: resp:  {"Count": 1}
    --- PASS: TestApi/diff_status/services-count/not-available/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 1}
    --- PASS: TestApi/diff_status/services-count/faster/1000/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 2}
    --- PASS: TestApi/diff_status/services-count/slower/1000/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 0}
    --- PASS: TestApi/big/services-count/available/1565031305/1565038505 (0.01s)
        api_test.go:147: resp:  {"Count": 1021}
    --- PASS: TestApi/big/services-count/not-available/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 13}
    --- PASS: TestApi/big/services-count/faster/1000/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 1014}
    --- PASS: TestApi/big/services-count/slower/1000/1565031305/1565038505 (0.00s)
        api_test.go:147: resp:  {"Count": 20}
PASS
ok  	github.com/ypapax/status_check/test	372.170s
```