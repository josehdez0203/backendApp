package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/josehdez0203/realstate/models"
)

func (m *PostgresDBRepo) AllMovies(genre ...int) ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}
	query := fmt.Sprintf(`
	select
	  id, title, release_date, runtime,
	  mpaa_rating, description, coalesce(image,''),
	  created_at, updated_at
	from
	  movies %s
	order by
	  title
	`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []*models.Movie

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		log.Println(movie.MPAARating)
		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *PostgresDBRepo) OneMovie(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	log.Println("Obt movies 🎥", id)

	query := `select id, title, release_date, runtime, mpaa_rating,
	description, coalesce(image,'') as image, created_at, updated_at
	from movies where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	var movie models.Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		log.Println("❌", err)
		return nil, err
	}
	// Obtener generos

	query = `select g.id, g.genre from movies_genres mg
	     left join genres g on (mg.genre_id = g.id)
	     where mg.movie_id = $1
	     order by g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(&g.ID, &g.Genre)
		if err != nil {
			return nil, err
		}

		genres = append(genres, &g)
	}
	movie.Genres = genres
	return &movie, nil
}

func (m *PostgresDBRepo) OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, title, release_date, runtime, mpaa_rating,
	description, coalesce(image,''), created_at, updated_at
	from movies where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	log.Println("movs", query, id)

	var movie models.Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		log.Println("Movies", err)
		return nil, nil, err
	}

	// Obtener generos

	query = `select g.id, g.genre from movies_genres mg
	     left join genres g on (mg.genre_id = g.id)
	     where mg.movie_id = $1
	     order by g.genre`

	rows, err := m.DB.QueryContext(ctx, query, id)
	log.Println("gen", query, id)
	if err != nil {
		log.Println("Genres", err)
		return nil, nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	var genresArray []int

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(&g.ID, &g.Genre)
		if err != nil {
			return nil, nil, err
		}
		genres = append(genres, &g)
		genresArray = append(genresArray, g.ID)
	}

	movie.Genres = genres
	movie.GenresArray = genresArray

	var allGenres []*models.Genre
	query = `select id, genre from genres order by genre`
	gRows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	defer gRows.Close()

	for gRows.Next() {
		var g models.Genre
		err := gRows.Scan(&g.ID, &g.Genre)
		if err != nil {
			return nil, nil, err
		}
		allGenres = append(allGenres, &g)
	}

	return &movie, allGenres, nil
}

func (m *PostgresDBRepo) ALlGenres() ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, genre, created_at, updated_at from genres order by genre`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var genres []*models.Genre

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)

	}
	return genres, nil
}

func (m *PostgresDBRepo) InsertMovie(movie models.Movie) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	statement := `insert into movies(title, description, release_date, runtime,
	mpaa_rating, created_at, updated_at, image) 
	values($1,$2,$3,$4,$5,$6,$7,$8) returning id`

	var newID int

	err := m.DB.QueryRowContext(ctx, statement,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.RunTime,
		movie.MPAARating,
		movie.CreatedAt,
		movie.UpdatedAt,
		movie.Image,
	).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (m *PostgresDBRepo) UpdateMovie(movie models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	statement := `update movies set title= $1, description = $2, release_date = $3,
	runtime = $4, mpaa_rating = $5, updated_at = $6, image = $7 where id = $8`

	_, err := m.DB.ExecContext(ctx, statement,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.RunTime,
		movie.MPAARating,
		movie.UpdatedAt,
		movie.Image,
		movie.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostgresDBRepo) UpdateMovieGenres(id int, genresIDs []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	statement := `delete from movies_genres where movie_id = $1`
	_, err := m.DB.ExecContext(ctx, statement, id)
	if err != nil {
		return err
	}

	for _, n := range genresIDs {
		statement := `insert into movies_genres(movie_id, genre_id) values($1,$2)`
		_, err := m.DB.ExecContext(ctx, statement, id, n)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *PostgresDBRepo) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	statement := `delete from movies where id = $1`
	_, err := m.DB.ExecContext(ctx, statement, id)
	if err != nil {
		return err
	}

	return nil
}
