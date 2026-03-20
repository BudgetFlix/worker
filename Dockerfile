FROM debian:bookworm-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
ffmpeg \
openjdk-17-jre-headless && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY worker.jar worker.jar

ENV NEW_DIR=/budgetflix/media/inbox/new
ENV PROCESS_DIR=/budgetflix/media/inbox/process
ENV DONE_DIR=/budgetflix/media/inbox/done
ENV ERROR_DIR=/budgetflix/media/inbox/error
ENV ERROR_LOG=/budgetflix/media/inbox/error/errorLog.txt
ENV DATA_BASE=/budgetflix/database/budgetflix.db
ENV MOVIE_SOURCE=/budgetflix/media/library/movies
ENV SERIES_SOURCE=/budgetflix/media/library/series

ENTRYPOINT ["java", "-jar", "worker.jar"]