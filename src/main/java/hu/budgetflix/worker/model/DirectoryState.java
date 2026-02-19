package hu.budgetflix.worker.model;

public class DirectoryState {
    private long stableSince = 0;
    private  boolean submitted = false;

    public DirectoryState () {

    }

    public boolean isSubmitted() {
        return submitted;
    }

    public void setSubmitted(boolean submitted) {
        this.submitted = submitted;
    }

    public long getStableSince() {
        return stableSince;
    }

    public void setStableSince(long stableSince) {
        this.stableSince = stableSince;
    }
}
