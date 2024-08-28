package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/yinnohs/simple-bank/db/mock"
	db "github.com/yinnohs/simple-bank/db/sqlc"
	"github.com/yinnohs/simple-bank/util"
)

func TestGetAccountAPI(t *testing.T) {

	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name:      "should return a status 200 when finding a correct account (OK)",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(
						gomock.Any(),
						gomock.Eq(account.ID)).
					Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		//TODO: add more cases
		{
			name:      "should return a status 404 when account is not found (NOT FOUND)",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(
						gomock.Any(),
						gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "should return an internal server error 500  (INTERNAL SERVER ERROR)",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(
						gomock.Any(),
						gomock.Eq(account.ID)).
					Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "should return an Bad request  error 400  (BAD REQUEST)",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountById(
						gomock.Any(),
						gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			//given
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/account/%d", tc.accountID)

			//when
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			//then
			tc.checkResponse(t, *recorder)
		})

	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance:  util.RandomBalance(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
