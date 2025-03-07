package repository

import (
	"database/sql"

	"github.com/josehdez0203/realstate/models"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id int) (*models.User, error)
	AddUser(user models.User) (*models.User, error)
	AllMovies(genre ...int) ([]*models.Movie, error)
	OneMovie(id int) (*models.Movie, error)
	OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error)
	ALlGenres() ([]*models.Genre, error)
	InsertMovie(movie models.Movie) (int, error)
	UpdateMovieGenres(id int, genresIDs []int) error
	UpdateMovie(movie models.Movie) error
	DeleteMovie(id int) error
}
