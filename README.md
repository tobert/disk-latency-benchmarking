# Disk Latency Benchmarking

I'm working on a couple talks and blog posts about disk latency. This repo contains a mix of
scripts, notes, configurations, and generated data used in the process.

## Scripts

There are a couple scripts in the scripts/ directory. More to come as I generate graphs, etc..

 latstat.pl - simple prototype to display data from /proc/diskstats
 log2json.pl - quick & dirty fio log to JSON converter

## Setup

My hardware setup is on my blog at http://albertptobey.blogspot.com/2014/03/benchmarking-disk-latency-setup.html

For posterity it is:


| Component    | Q | Description                         |
| ------------ | - | ----------------------------------- |
| Case         | 1 | Cooler Master HAF XB                |
| Power Supply | 1 | COOLMAX CU series CU-700B 700W      |
| Motherboard  | 1 | Intel S1200BTL                      |
| CPU          | 1 | Intel Xeon E31270 3.4Ghz            |
| Memory       | 4 | Kingston KVR1333D3E9SK2             |
| Root Drive   | 2 | Seagate ST9500530NS Enterprise SATA |
| Graphics     | 1 | XFX Radeon 6450 2GB                 |
| PCIe SSD     | 1 | FusionIO ioDrive II                 |
| PCIe SAS     | 1 | LSI Logic LSI00346 9300-4i SGL      |
| SAS chassis  | 1 | Thermaltake RC1400101A MAX-1542     |

## Drives

Samsung 840 Pro 128GB (2.5" SATA SSD)
FusionIO ioDrive II (PCIe 8x SSD)
Western Digital WD2500KS (3.5" SATA 7,200RPM)
Seagate ST9500430SS (2.5" SAS, will also test RAID10 on these)
Western Digital Velociraptor WD3000BLFS (2.5" SATA 10,000RPM)
Western Digital WD5002AALX (3.5" SATA 7,200RPM, also RAID1)

