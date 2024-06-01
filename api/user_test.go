package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	mockdb "simple_bank/db/mock"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {

	v, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.ComparePassword(e.password, v.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = v.HashedPassword
	return reflect.DeepEqual(e.arg, v)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("create user params %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(x db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{x, password}
}

func Test_CreateUser(t *testing.T) {
	user, password := createRandomUser(t)

	tests := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockIStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"user_name": user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockIStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					Email:    user.Email,
					FullName: user.FullName,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requestBodyMatchingUser(t, recorder.Body, user)
			},
		},
		// {
		// 	name:      "Not found",
		// 	accountID: account.ID,
		// 	buildStubs: func(store *mockdb.MockIStore) {
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusNotFound, recorder.Code)
		// 	},
		// },
		// {
		// 	name:      "InternalServerError",
		// 	accountID: account.ID,
		// 	buildStubs: func(store *mockdb.MockIStore) {
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		// 	},
		// },
		// {
		// 	name:      "Invalid ID",
		// 	accountID: 0,
		// 	buildStubs: func(store *mockdb.MockIStore) {
		// 		store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusBadRequest, recorder.Code)
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := mockdb.NewMockIStore(ctrl)
			tt.buildStubs(db)

			server := NewServer(db)
			recorder := httptest.NewRecorder()

			url := "/users"
			data, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)

			tt.checkResponse(t, recorder)
		})
	}
}

func createRandomUser(t *testing.T) (user db.User, password string) {

	password = util.RandomPassword()
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		FullName:       util.RandomName(),
	}
	return
}

func requestBodyMatchingUser(t *testing.T, body *bytes.Buffer, user db.User) {

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)

}
