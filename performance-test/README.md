# Performance Test

Performance test for hinst-website.
Running on Orange Pi Zero 2W single board computer with 4x cores, 4 GB of RAM, microSD storage.
* Using Golang standard library for HTTP receivers
	* Throughput: 130.0 blog posts per second (avg: 93.2), 424.0 requests per second (avg: 303.9)
* Using Huma Rest framework with OpenAPI reflection support
	* Throughput: 123.5 blog posts per second (avg: 80.0), 402.8 requests per second (avg: 260.9)
