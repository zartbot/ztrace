# ztrace

Traceroute in Go

![](https://github.com/zartbot/ztrace/blob/main/ztrace.png?raw=true)

## Motivation
Most network engineers need fast traceroute result and correlate it with service provider information and Geo locations.

ztrace is a go based solution to leverage multithread capability which could get the result in just one second.

The ztrace is a common lib for traceroute service, you can use it as a probe and export the stats result to 3rd party database.

The example is just output them to the console.

## How to build

```bash
git clone https://github.com/zartbot/ztrace
cd example 
make
cd build/ztrace
sudo ./ztrace_linux_x86 www.cisco.com
```
> The build target only support MAC and linux with both arm/x86 platform. you could load and run it on a Cisco Router.

## Usage


```
Usage:
  ./ztrace [-src source] [-proto protocol] [-ttl ttl] [-rate packetRate] [-wide Widescreen mode] [-path NumOfECMPPath] host
Example:
 ./ztrace www.cisco.com
 ./ztrace -ttl 30 -rate 1 -path 8 -wide www.cisco.com
Option:
  -lat float
    	Latitude (default 31.02)
  -long float
    	Longitude (default 121.1)
  -p int
    	Max ECMP Number (default 16)
  -path int
    	Max ECMP Number (default 16)
  -proto string
    	Protocol[icmp|tcp|udp] (default "udp")
  -r float
    	Packet Rate per second (default 1)
  -rate float
    	Packet Rate per second (default 1)
  -src string
    	Source 
  -ttl int
    	Max TTL (default 64)
  -w	Widescreen mode
  -wide
    	Widescreen mode

```
> Linux platform support all the traceroute protocol, The MAC raw socket can only listen ICMP/UDP, TCP is not supported.

```bash
[www.cisco.com]Traceroute Report

+-------+-----------------------+--------------------------------+------------------+------------------+------------+--------------------------------+-----------------+----------------+------------+------------+-------+
| TTL   |        Server         |              Name              |       City       |     Country      |    ASN     |               SP               | Distance[tRTT]  |      p95       |  Latency   |   Jitter   | Loss  |
+-------+-----------------------+--------------------------------+------------------+------------------+------------+--------------------------------+-----------------+----------------+------------+------------+-------+
|     1 | 192.168.101.1         | zartRT                         |                  |                  | 0          |                                |                 |         3.25ms |     1.48ms |     0.98ms |  0.0% |
|     2 | 192.168.1.1           |                                |                  |                  | 0          |                                |                 |         6.48ms |     3.69ms |     1.70ms |  0.0% |
|     3 | 100.65.0.1            |                                |                  |                  | 0          |                                |                 |        34.69ms |    25.52ms |     6.92ms |  0.0% |
|     4 | 61.152.54.125         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 12ms] |        34.03ms |    15.81ms |     9.30ms |  0.0% |
|       | 61.152.53.149         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 12ms] |        31.60ms |    17.86ms |     8.41ms |  0.0% |
|     5 | 61.152.25.66          |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        31.71ms |    20.10ms |     6.22ms |  0.0% |
|       | 61.152.25.190         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        30.96ms |    22.67ms |     7.75ms | 25.0% |
|       | 61.152.24.110         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        32.14ms |    21.01ms |     9.02ms |  0.0% |
|       | 61.152.24.118         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        34.29ms |    20.10ms |    10.47ms |  0.0% |
|       | 61.152.25.18          |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        30.95ms |    21.97ms |     8.66ms | 25.0% |
|       | 61.152.25.174         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        31.87ms |    22.24ms |     9.08ms | 25.0% |
|       | 61.152.25.98          |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        31.73ms |    22.17ms |     9.03ms | 25.0% |
|       | 101.95.88.126         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        33.79ms |    19.35ms |    10.38ms | 25.0% |
|       | 61.152.24.142         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        24.99ms |    17.72ms |     5.29ms | 25.0% |
|       | 101.95.88.122         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        25.97ms |    17.69ms |     6.43ms | 25.0% |
|       | 61.152.24.242         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        16.80ms |    14.05ms |     2.21ms | 25.0% |
|       | 61.152.24.226         |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        31.93ms |    22.25ms |     8.77ms | 25.0% |
|       | 61.152.25.30          |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        25.21ms |    19.20ms |     5.34ms | 25.0% |
|       | 61.152.25.50          |                                | Shanghai         | China            | 4812       | China Telecom Group            |     29km[ 15ms] |        31.96ms |    22.95ms |     8.33ms | 25.0% |
|     6 | 202.97.24.138         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 29ms] |        32.28ms |    23.04ms |     8.50ms | 25.0% |
|       | 202.97.57.145         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 29ms] |        31.94ms |    18.64ms |     7.69ms |  0.0% |
|       | 202.97.50.154         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 29ms] |        31.73ms |    22.79ms |     8.15ms | 25.0% |
|       | 202.97.57.25          |                                |                  | China            | 4134       | Chinanet                       |    805km[ 29ms] |        32.11ms |    23.05ms |     8.22ms | 25.0% |
|       | 202.97.94.237         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 29ms] |        19.93ms |    13.78ms |     5.27ms | 25.0% |
|       | 202.97.61.6           |                                |                  | China            | 4134       | Chinanet                       |    805km[ 29ms] |        30.82ms |    22.95ms |     7.23ms | 25.0% |
|     7 | 202.97.83.22          |                                |                  | China            | 4134       | Chinanet                       |    805km[ 32ms] |        31.86ms |    20.36ms |     7.55ms |  0.0% |
|       | 202.97.83.62          |                                |                  | China            | 4134       | Chinanet                       |    805km[ 32ms] |        30.87ms |    22.38ms |     8.13ms | 25.0% |
|       | 202.97.51.161         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 32ms] |        31.96ms |    18.99ms |     9.75ms |  0.0% |
|       | 202.97.83.18          |                                |                  | China            | 4134       | Chinanet                       |    805km[ 32ms] |        32.25ms |    23.22ms |     7.91ms | 25.0% |
|       | 202.97.74.1           |                                |                  | China            | 0          |                                |    805km[ 32ms] |        32.11ms |    23.32ms |     7.48ms |  0.0% |
|       | 202.97.85.34          |                                |                  | China            | 4134       | Chinanet                       |    805km[ 32ms] |        33.88ms |    23.23ms |     8.58ms |  0.0% |
|       | 202.97.83.38          |                                |                  | China            | 4134       | Chinanet                       |    805km[ 32ms] |        34.08ms |    18.92ms |    10.19ms |  0.0% |
|       | 202.97.12.197         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 32ms] |        33.91ms |    21.14ms |     7.98ms |  0.0% |
|     8 | 202.97.27.210         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       167.09ms |   160.08ms |     3.19ms |  0.0% |
|       | 202.97.41.202         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       163.98ms |   157.39ms |     7.49ms |  0.0% |
|       | 202.97.59.170         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       167.17ms |   163.51ms |     3.47ms |  0.0% |
|       | 202.97.90.222         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       180.62ms |   176.79ms |     2.92ms | 25.0% |
|       | 202.97.99.154         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       140.76ms |   136.22ms |     3.26ms | 25.0% |
|       | 202.97.52.178         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       143.87ms |   142.53ms |     0.95ms | 25.0% |
|       | 202.97.71.194         |                                |                  | China            | 0          |                                |    805km[ 35ms] |       166.01ms |   162.73ms |     2.67ms | 25.0% |
|       | 202.97.27.218         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       167.76ms |   164.88ms |     2.17ms | 25.0% |
|       | 202.97.89.142         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       173.84ms |   167.85ms |     4.40ms | 25.0% |
|       | 202.97.59.153         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       147.65ms |   141.30ms |     4.00ms |  0.0% |
|       | 202.97.58.170         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 35ms] |       150.41ms |   144.69ms |     5.47ms | 25.0% |
|     9 | 202.97.83.242         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 38ms] |       172.10ms |   145.78ms |    11.39ms |  0.0% |
|       | 202.97.92.37          |                                |                  | China            | 4134       | Chinanet                       |    805km[ 38ms] |       180.61ms |   162.57ms |    10.29ms |  0.0% |
|       | 202.97.50.78          |                                |                  | China            | 4134       | Chinanet                       |    805km[ 38ms] |       198.09ms |   184.91ms |     7.71ms |  0.0% |
|    10 | 218.30.53.109         |                                |                  | China            | 4134       | Chinanet                       |    805km[ 41ms] |       183.84ms |   160.34ms |    10.52ms |  0.0% |
|       | 129.250.9.73          | ae-63.a01.snjsca04.us.bb.gin.n |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[182ms] |       183.42ms |   179.63ms |     3.52ms |  0.0% |
|    11 | 129.250.2.105         | ae-7.r24.lsanca07.us.bb.gin.nt |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[185ms] |       186.21ms |   147.38ms |    14.07ms |  0.0% |
|       | 129.250.3.162         | ae-8.r25.snjsca04.us.bb.gin.nt |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[185ms] |       178.40ms |   161.49ms |     8.80ms |  0.0% |
|    12 | 129.250.7.68          | ae-7.r20.dllstx14.us.bb.gin.nt |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[188ms] |       195.09ms |   173.26ms |    10.37ms |  0.0% |
|       | 129.250.4.155         | ae-8.r21.dllstx14.us.bb.gin.nt |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[188ms] |       210.32ms |   202.74ms |     3.85ms |  0.0% |
|    13 | 129.250.4.117         | ae-11.r22.atlnga05.us.bb.gin.n |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[191ms] |       226.28ms |   219.00ms |     5.45ms | 25.0% |
|       | 129.250.3.55          | ae-1.r21.dllstx14.us.bb.gin.nt |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[191ms] |       205.43ms |   183.17ms |    10.01ms |  0.0% |
|    14 | 129.250.4.117         | ae-11.r22.atlnga05.us.bb.gin.n |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[194ms] |       214.90ms |   195.51ms |     7.42ms |  0.0% |
|       | 129.250.5.201         | ae-9.r04.atlnga05.us.bb.gin.nt |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[194ms] |       224.31ms |   219.58ms |     4.14ms | 25.0% |
|    15 | 129.250.5.201         | ae-9.r04.atlnga05.us.bb.gin.nt |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[197ms] |       217.47ms |   197.44ms |     8.71ms |  0.0% |
|       | 128.241.8.127         | ae-3.akamai.atlnga05.us.bb.gin |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[197ms] |       272.28ms |   259.82ms |     6.66ms |  0.0% |
|    16 | 128.241.8.127         | ae-3.akamai.atlnga05.us.bb.gin |                  | United States    | 2914       | NTT-COMMUNICATIONS-2914        |  11366km[200ms] |       270.39ms |   245.43ms |    13.80ms |  0.0% |
|       | 23.203.144.195        | ae2.datasite-atl2.netarch.akam |                  | United States    | 20940      | Akamai International B.V.      |  11366km[200ms] |      2371.35ms |  2371.35ms |     0.00ms | 75.0% |
|    17 | 23.203.144.195        | ae2.datasite-atl2.netarch.akam |                  | United States    | 20940      | Akamai International B.V.      |  11366km[203ms] |      2049.29ms |  2045.02ms |     2.36ms |  0.0% |
+-------+-----------------------+--------------------------------+------------------+------------------+------------+--------------------------------+-----------------+----------------+------------+------------+-------+
+-------+-----------------------+--------------------------------+------------------+------------------+------------+--------------------------------+-----------------+----------------+------------+------------+-------+
| Probe |        Server         |              Name              |       City       |     Country      |    ASN     |               SP               | Distance[tRTT]  |      p95       |  Latency   |   Jitter   | Loss  |
+-------+-----------------------+--------------------------------+------------------+------------------+------------+--------------------------------+-----------------+----------------+------------+------------+-------+
| TCP   | 104.102.227.49:80     |                                | Atlanta          | United States    | 16625      | AKAMAI-AS                      |  12333km[194ms] |       261.88ms |   243.37ms |    12.55ms | 25.0% |
|       | 104.102.227.49:443    |                                | Atlanta          | United States    | 16625      | AKAMAI-AS                      |  12333km[194ms] |      1271.75ms |   370.10ms |   271.74ms | 12.5% |
+-------+-----------------------+--------------------------------+------------------+------------------+------------+--------------------------------+-----------------+----------------+------------+------------+-------+
```