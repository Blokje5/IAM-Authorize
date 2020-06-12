package server

import "testing"

func TestErrorResponse_Error(t *testing.T) {

	tests := []struct {
		name          string
		errorResponse *ErrorResponse
		want          string
	}{
		{
			"Calling Error() an ErrorResponse returns a json string",
			NewErrorResponse(ErrConflict, "Conflict", "Uniqueness constraint violation", 409),
			`{"type":"http://localhost:8080/errors/conflict","title":"Conflict","detail":"Uniqueness constraint violation","status":409}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.errorResponse.Error(); got != tc.want {
				t.Errorf("ErrorResponse.Error() = %v, want %v", got, tc.want)
			}
		})
	}
}
