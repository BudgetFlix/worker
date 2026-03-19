FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
ffmpeg \
openjdk-21-jre-headless && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY worker.jar worker.jar

ENV NEW_DIR=/media/inbox/new
ENV PROCESS_DIR=/media/inbox/process
ENV DONE_DIR=/media/inbox/done
ENV ERROR_DIR=/media/inbox/error
ENV ERROR_LOG=/media/inbox/error/errorLog.txt
ENV DATA_BASE=/media/database/budgetflix.db
ENV MOVIE_SOURCE=/media/library/movies
ENV SERIES_SOURCE=/media/library/series

ENTRYPOINT ["java", "-jar", "worker.jar"]