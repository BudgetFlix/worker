package hu.budgetflix.worker.model.database;

import com.fasterxml.jackson.databind.ObjectMapper;
import hu.budgetflix.worker.model.MediaType;
import hu.budgetflix.worker.model.Stat;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

public  class JsonReader {
    private static final ObjectMapper mapper = new ObjectMapper();

    public static MediaType extractTheType (Path jsonPath) {
        return MediaType.MOVIE;
    }

    public static Stat jsonToObject(Path directory) {
        try{

            return mapper.readValue(findJsonFile(directory).toFile(),Stat.class);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    private static Path findJsonFile(Path folderPath) throws IOException {
        try (var stream = Files.list(folderPath)) {
            return stream
                    .filter(Files::isRegularFile)
                    .filter(p -> p.toString().endsWith(".json"))
                    .findFirst()
                    .orElseThrow(() ->
                            new IllegalStateException("Nincs JSON fájl a mappában: " + folderPath)
                    );
        }
    }
}
