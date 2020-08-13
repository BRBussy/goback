package authorisation

import (
	"github.com/BRBussy/goback/pkg/jsonrpc"
	"github.com/BRBussy/goback/pkg/mongo/filter"
	"github.com/BRBussy/goback/pkg/role"
	"github.com/BRBussy/goback/pkg/security/claims"
	"github.com/BRBussy/goback/pkg/user"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Middleware struct {
	authoriser Authoriser
	userStore  user.Store
	roleStore  role.Store
}

func NewMiddleware(
	authoriser Authoriser,
	userStore user.Store,
	roleStore role.Store,
) *Middleware {
	return &Middleware{
		authoriser: authoriser,
		userStore:  userStore,
		roleStore:  roleStore,
	}
}

func (a *Middleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse claims from context
		c, err := claims.ParseFromContext(r.Context())
		if err != nil {
			log.Error().Err(err).Msg("unable to parse claims from context")
			http.Error(w, "Unauthorized", http.StatusInternalServerError)
			return
		}

		// get the rpc method being requested
		rpcMethod, err := jsonrpc.GetRPCMethod(r)
		if err != nil {
			log.Error().Err(err).Msg("unable to get rpc method")
			http.Error(w, "Unauthorized", http.StatusInternalServerError)
			return
		}

		// retrieve the requesting user
		retrieveUserResponse, err := a.userStore.Retrieve(
			user.RetrieveRequest{
				Filter: filter.NewTextExactFilter(
					"id",
					c.UserID,
				),
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve user")
			http.Error(w, "Unauthorized", http.StatusInternalServerError)
			return
		}

		// retrieve all of the roles assigned to the existing user
		listRolesResponse, err := a.roleStore.List(
			role.ListRequest{
				Filter: filter.NewTextExactListFilter(
					"id",
					retrieveUserResponse.User.RoleIDs,
				),
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to retrieve user's roles")
			http.Error(w, "Unauthorized", http.StatusInternalServerError)
			return
		}

		// check through all of the permissions of all of the user's roles
		// to see if they have a permission matching the requested service
	nextRole:
		for _, r := range listRolesResponse.Records {
			for _, perm := range r.Permissions {
				if perm == rpcMethod {
					break nextRole
				}
			}
			// if execution reaches here the user does not
			// have permission to access this service
			log.Warn().Msg("user does not have permission to access service: " + rpcMethod)
			http.Error(w, "Unauthorized", http.StatusInternalServerError)
			return
		}

		// if execution reaches here they have the required permission
		// and thus pass the authorisation middleware
		next.ServeHTTP(w, r)
	})
}
