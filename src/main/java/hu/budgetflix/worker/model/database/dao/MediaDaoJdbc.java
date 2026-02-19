package hu.budgetflix.worker.model.database.dao;

import hu.budgetflix.worker.model.media.Movie;

import javax.sql.DataSource;
import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.time.LocalDateTime;

public class MediaDaoJdbc implements MediaDao {
    private DataSource dataSource;

    public MediaDaoJdbc (DataSource dataSource) {
        this.dataSource = dataSource;
    }


    @Override
    public Long addNewMedia(Movie movie) {
            String sql = "INSERT INTO movie ( title,original_filename,status,created_at) VALUES (?,?,?,?) RETURNING id";
        try (Connection connection = dataSource.getConnection()) {
            PreparedStatement st = connection.prepareStatement(sql);


            st.setString(1, movie.getName());
            st.setString(2, movie.getName());
            st.setString(3, movie.getStatus().toString());
            st.setString(4, LocalDateTime.now().toString());

            st.executeUpdate();

            ResultSet rs = st.getGeneratedKeys();
            if(rs.next()){
                return rs.getLong(1);
            }

        } catch (SQLException e) {
            throw new RuntimeException(e);
        }
        return (long) 0;
    }

    @Override
    public void updatePatch(Movie movie) {
        String sql = "UPDATE movie SET hls_path = ? WHERE id = ? ";
        try (Connection con = dataSource.getConnection()) {
            PreparedStatement st =  con.prepareStatement(sql);

            st.setString(1,movie.getVideo().getOutPath().toString());
            st.setLong(2,movie.getId());

            st.executeUpdate();

        } catch (SQLException e) {
            throw new RuntimeException(e);
        }

    }

    @Override
    public void updateStatus(Movie movie) {
        String sql = "UPDATE movie SET status = ? WHERE id = ? ";
        try (Connection con = dataSource.getConnection()) {
            PreparedStatement st =  con.prepareStatement(sql);

            st.setString(1,movie.getStatus().toString());
            st.setLong(2,movie.getId());

            st.executeUpdate();

        } catch (SQLException e) {
            throw new RuntimeException(e);
        }

    }
}
