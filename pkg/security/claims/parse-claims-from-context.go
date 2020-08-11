package claims

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BRBussy/goback/pkg/exception"
)

func ParseFromContext(ctx context.Context) (Claims, error) {
	// look for claims in context
	marshalledClaimsInterface := ctx.Value("Claims")
	if marshalledClaimsInterface == nil {
		return nil, NewErrClaimsNotInContext()
	}

	// try an cast claims to string
	marshalledClaims, ok := marshalledClaimsInterface.(json.RawMessage)
	if !ok {
		return nil, exception.NewErrUnexpected(
			errors.New("unable to cast context to json.RawMessage"),
		)
	}

	// parse the claims from json
	var serializedClaims SerializedClaims
	if err := json.Unmarshal(marshalledClaims, &serializedClaims); err != nil {
		return nil, NewErrJSONUnmarshallError(err)
	}

	return serializedClaims.Claims, nil
}
