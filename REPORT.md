## Two-Choice

```
# HELP pblb_processed_total The total number of processed requests
# TYPE pblb_processed_total counter
pblb_processed_total{node="nginx_a",status_class="2xx"} 37468
pblb_processed_total{node="nginx_b",status_class="2xx"} 36659
pblb_processed_total{node="nginx_c",status_class="2xx"} 37049

Lifting the server siege...
Transactions:		      110923 hits
Availability:		      100.00 %
Elapsed time:		       60.75 secs
Data transferred:	       65.06 MB
Response time:		        0.14 secs
Transaction rate:	     1825.89 trans/sec
Throughput:		        1.07 MB/sec
Concurrency:		      253.95
Successful transactions:      110923
Failed transactions:	           0
Longest transaction:	        0.50
Shortest transaction:	        0.00
```

## Round Robin

```
# HELP pblb_processed_total The total number of processed requests
# TYPE pblb_processed_total counter
pblb_processed_total{node="nginx_a",status_class="2xx"} 37260
pblb_processed_total{node="nginx_b",status_class="2xx"} 37260
pblb_processed_total{node="nginx_c",status_class="2xx"} 37260

Lifting the server siege...
Transactions:		      111525 hits
Availability:		       99.99 %
Elapsed time:		       60.25 secs
Data transferred:	       65.41 MB
Response time:		        0.14 secs
Transaction rate:	     1851.04 trans/sec
Throughput:		        1.09 MB/sec
Concurrency:		      254.01
Successful transactions:      111525
Failed transactions:	          10
Longest transaction:	        0.68
Shortest transaction:	        0.00
```
