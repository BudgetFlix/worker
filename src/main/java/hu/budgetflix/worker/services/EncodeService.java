package hu.budgetflix.worker.services;

import java.io.IOException;
import java.nio.file.Path;

public interface EncodeService {
    void buildUpStructure(Path directory) throws IOException;
    void startEncode ();

    void shutDown();
}
