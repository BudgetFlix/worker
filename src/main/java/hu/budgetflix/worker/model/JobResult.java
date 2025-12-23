package hu.budgetflix.worker.model;

public record JobResult(int exitCode, String errorTail) {
    public boolean success (){return exitCode == 0;}
}
