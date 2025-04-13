# Performance before and after SQLite migration

Measured on page `hinst-website/#/personal-goals/247488?activePostDate=2025-04-07+19%3A37%3A53`, on local network

## Before
* goalPosts: 33 ms
* goalPost: 16 ms
* goal: 3 ms
* TOTAL: 52 ms

## After
* goalPosts: 16 ms
* goalPost: 8 ms
* goal: 8 ms
* TOTAL: 32 ms
