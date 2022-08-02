## Two-Choice

### 60 Second Siege

```
# HELP pblb_processed_total The total number of processed requests
# TYPE pblb_processed_total counter
pblb_processed_total{node="nginx_a",status_class="2xx"} 37468
pblb_processed_total{node="nginx_b",status_class="2xx"} 36659
pblb_processed_total{node="nginx_c",status_class="2xx"} 37049
                                                total: 111176
```

```
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

```
# HELP pblb_processed_total The total number of processed requests
# TYPE pblb_processed_total counter
pblb_processed_total{node="nginx_a",status_class="2xx"} 633222
pblb_processed_total{node="nginx_b",status_class="2xx"} 632754
pblb_processed_total{node="nginx_c",status_class="2xx"} 630716

Lifting the server siege...
Transactions:		     1896517 hits
Availability:		      100.00 %
Elapsed time:		      900.79 secs
Data transferred:	     1112.33 MB
Response time:		        0.12 secs
Transaction rate:	     2105.39 trans/sec
Throughput:		        1.23 MB/sec
Concurrency:		      253.57
Successful transactions:     1896517
Failed transactions:	           2
Longest transaction:	        0.93
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

```
# HELP pblb_processed_total The total number of processed requests
# TYPE pblb_processed_total counter
pblb_processed_total{node="nginx_a",status_class="2xx"} 641796
pblb_processed_total{node="nginx_b",status_class="2xx"} 641796
pblb_processed_total{node="nginx_c",status_class="2xx"} 641795

Lifting the server siege...
Transactions:		     1925132 hits
Availability:		      100.00 %
Elapsed time:		      900.35 secs
Data transferred:	     1129.11 MB
Response time:		        0.12 secs
Transaction rate:	     2138.20 trans/sec
Throughput:		        1.25 MB/sec
Concurrency:		      254.34
Successful transactions:     1925132
Failed transactions:	           4
Longest transaction:	        0.99
Shortest transaction:	        0.00
```