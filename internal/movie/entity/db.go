package entity

type MovieRepo struct {
	ID       int64  `db:"id" json:"id"`
	Name     string `db:"name" json:"name" valid:"required"`
	Duration int    `db:"duration" json:"duration" valid:"required,range(1|1000)"`
	Genre    string `db:"genre" json:"genre" valid:"required"`
}
