package getassignedproblems_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Pickausernaame/chat-service/internal/types"
	getassignedproblems "github.com/Pickausernaame/chat-service/internal/usecases/manager/get-assigned-problems"
)

func TestRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		request getassignedproblems.Request
		wantErr bool
	}{
		// Negative
		{
			name: "require manager id",
			request: getassignedproblems.Request{
				ManagerID: types.UserIDNil,
			},
			wantErr: true,
		},

		// Positive
		{
			name: "success",
			request: getassignedproblems.Request{
				ManagerID: types.NewUserID(),
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
