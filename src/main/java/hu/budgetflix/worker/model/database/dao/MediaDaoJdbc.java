package hu.budgetflix.worker.model.database.dao;

import hu.budgetflix.worker.model.Video;

import javax.sql.DataSource;
import java.nio.file.Path;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.time.LocalDateTime;

public class MediaDaoJdbc implements MediaDao {
    private DataSource dataSource;

    public MediaDaoJdbc (DataSource dataSource) {
        this.dataSource = dataSource;
    }


    @Override
    public void addNewMedia(Video vide, Path out) {
            String sql = "INSERT INTO movie (id, title,original_filename,status,hls_path,created_at) VALUES (?,?,?,?,?,?)";
        try (Connection connection = dataSource.getConnection()) {
            PreparedStatement st = connection.prepareStatement(sql);

            st.setString(1,vide.getId().toString());
            st.setString(2, vide.getName());
            st.setString(3, vide.getName());
            st.setString(4, vide.getStatus().toString());
            st.setString(5,out.resolve("index.m3u8").toString());
            st.setString(6, LocalDateTime.now().toString());

            st.executeUpdate();

        } catch (SQLException e) {
            throw new RuntimeException(e);
        }

    }

    @Override
    public void updateStatus(Video vide) {
        String sql = "UPDATE movie SET status = ? WHERE id = ? ";
        try (Connection con = dataSource.getConnection()) {
            PreparedStatement st =  con.prepareStatement(sql);

            st.setString(1,vide.getStatus().toString());
            st.setString(2,vide.getId().toString());

            st.executeUpdate();

        } catch (SQLException e) {
            throw new RuntimeException(e);
        }

    }
}
