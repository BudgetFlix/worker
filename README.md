#  BudgetFlix Media Worker

> Media ingestion and processing service for BudgetFlix.

---

## Overview

The **Media Worker** is responsible for processing newly uploaded video files and preparing them for streaming.

It runs as a **Java 17 daemon** and continuously watches a directory for incoming media.

---

##  Status

>  In active development

---
## Responsibilities

When a new file is detected, the worker:

*  Detects and validates media files
*  Extracts basic metadata
*  Creates the required folder structure
*  Processes video using FFmpeg
*  Generates HLS streams (`.m3u8`, `.ts`)
*  Creates or updates database records
*  Queue-based processing

---

## Processing Flow

```id="flow1"
File detected
   ↓
Validation
   ↓
Database record created
   ↓
HLS output generated
   ↓
Processing (FFmpeg)
   ↓
Database updated
```


---

## Tech Stack

* Java 17
* FFmpeg
* File system watcher (WatchService)
* SQLite database

---

##  Folder Structure (example)

```id="folder1"
/library
  /movies
    /id
      /hls
        index.m3u8
        segment_000.ts
```

---

## FFmpeg Processing

The worker uses FFmpeg to:

* convert video to HLS format
* segment video into `.ts` files
* generate `.m3u8` playlist

Example output:

```
index.m3u8
segment_000.ts
segment_001.ts
...
```

---

##  Notes

* The worker must have access to:

    * upload directory
    * output media directory
* Ensure FFmpeg is properly installed
* Large files may take time to process

---

##  Future Improvements

* Retry mechanism for failed jobs
* Parallel processing
* Metadata extraction (duration, resolution)
* Thumbnail generation
* Progress tracking
* Logging & monitoring
* Generates thumbnails 


---

##  Design Philosophy

* Keep processing isolated from API
* Make pipeline deterministic and repeatable
* Prepare media for efficient streaming (HLS)

---

##  Related

* `budgetflix` → frontend
* `budgetflix-backend` → API & streaming access

---

