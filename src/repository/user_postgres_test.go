package repository

import (
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestPostgres_showUserSegments(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	type args struct {
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    string
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"segments_list"}).
					AddRow("AVITO_SALE,AVITO_CALL")

				mock.ExpectQuery(`SELECT segments_list FROM (.+) WHERE user_id = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			input: args{
				userId: 1,
			},
			want: "AVITO_SALE,AVITO_CALL",
		},
		{
			name: "No segments",
			mock: func() {
				rows := sqlmock.NewRows([]string{"segments_list"})

				mock.ExpectQuery(`SELECT segments_list FROM (.+) WHERE user_id = \$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
			input: args{
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.ShowUserSegments(tt.input.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
