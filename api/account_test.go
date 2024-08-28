package api

import (
	"bytes"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().
		GetAccountById(
			gomock.Any(),
			gomock.Eq(account.ID)).
		Times(1).Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/api/v1/account/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	server.router.ServeHTTP(recorder, request)

	// check respones inside the recorder

	require.Equal(t, http.StatusOK, recorder.Code)

	// check body

	requireBodyMatchAccount(t, recorder.Body, account)
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
