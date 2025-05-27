package authz

// const (
// 	_bearer        = "Bearer"
// 	_authorization = "authorization"
// )

// type authZServer struct {
// 	token token.Token
// 	cache cache.UserCache
// 	ip    ip.IPAccess
// 	ext   extractor.Extractor
// }

// NewServer creates a new authorization server.
// func NewServer(cache cache.UserCache, token token.Token, ip ip.IPAccess) authv3.AuthorizationServer {
// 	return &authZServer{
// 		token: token,
// 		cache: cache,
// 		ip:    ip,
// 		ext:   extractor.New(),
// 	}
// }

// Check implements authorization's Check interface which performs authorization check based on the
// attributes associated with the incoming request.
// func (s *authZServer) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
// 	forwardedIP := req.Attributes.Request.Http.Headers["x-forwarded-for"]
// 	err := s.checkIPAccess(ctx, forwardedIP)
// 	if err != nil {
// 		return buildDeniedResponse(int32(rpc.PERMISSION_DENIED), typev3.StatusCode_Forbidden), nil
// 	}
// 	authorization := req.Attributes.Request.Http.Headers[_authorization]
// 	tenantID := req.Attributes.Request.Http.Headers[string(constants.TenantID)]
// 	extracted := strings.Fields(authorization)
// 	if len(extracted) == 2 && extracted[0] == _bearer {
// 		tk, u, err := s.verifyAccessToken(ctx, extracted[1], tenantID)
// 		if err == nil {
// 			return &authv3.CheckResponse{
// 				HttpResponse: &authv3.CheckResponse_OkResponse{
// 					OkResponse: &authv3.OkHttpResponse{
// 						Headers: s.createHeaders(tk, u),
// 					},
// 				},
// 				Status: &status.Status{
// 					Code: int32(rpc.OK),
// 				},
// 			}, nil
// 		} else {
// 			if err == errTokenInactive {
// 				// Using StatusCode_Gone for the reason that kick by another device
// 				return buildDeniedResponse(int32(rpc.UNAUTHENTICATED), typev3.StatusCode_Conflict), nil
// 			}
// 		}
// 	} else {
// 		global.Log.Error("invalid header", errCacheNotFound)
// 	}

// 	return buildDeniedResponse(int32(rpc.UNAUTHENTICATED), typev3.StatusCode_Unauthorized), nil
// }

// func buildDeniedResponse(outerCode int32, innerCode typev3.StatusCode) *authv3.CheckResponse {
// 	return &authv3.CheckResponse{
// 		Status: &status.Status{
// 			Code: outerCode,
// 		},
// 		HttpResponse: &authv3.CheckResponse_DeniedResponse{
// 			DeniedResponse: &authv3.DeniedHttpResponse{
// 				Status: &typev3.HttpStatus{
// 					Code: innerCode,
// 				},
// 			},
// 		},
// 	}
// }

// var (
// 	errCouldNotParseToken      = errors.New("could not parse token")
// 	errTokenInvalid            = errors.New("invalid token")
// 	errCacheNotFound           = errors.New("cache token not found")
// 	errTokenInactive           = errors.New("inactive token")
// 	errUserNotFound            = errors.New("user not found")
// 	errUserCredentialsNotFound = errors.New("user credentials not found")
// 	errIPAccessDenied          = errors.New("ip access denied")
// )

// func (s *authZServer) verifyAccessToken(ctx context.Context, accessToken, tenantID string) (*signer.Token, *dtos.UserRes, error) {
// 	tk, err := s.token.AccessToken().Parse(accessToken)
// 	if err != nil {
// 		global.Log.Error("could not parse token", err)
// 		return nil, nil, errCouldNotParseToken
// 	}
// 	if tk.Valid() != nil {
// 		global.Log.Error("invalid token", errTokenInvalid)
// 		return nil, nil, errTokenInvalid
// 	}
// 	if tk.TenantID != tenantID {
// 		global.Log.Error("invalid TenantID", errTokenInvalid)
// 		return nil, nil, errTokenInvalid
// 	}

// 	u, err := s.cache.Get(ctx, tk.Subject)
// 	if err != nil {
// 		global.Log.Error("could not find user in cache", err)
// 		return nil, nil, errUserNotFound
// 	}

// 	return tk, &dtos.UserRes{Id: u.GetId()}, nil
// }

// func (s *authZServer) createHeaders(tk *signer.Token, user *dtos.UserRes) []*corev3.HeaderValueOption {
// 	headers := []*corev3.HeaderValueOption{
// 		{
// 			Append: &wrapperspb.BoolValue{Value: false},
// 			Header: &corev3.HeaderValue{
// 				Key:   string(constants.TokenID),
// 				Value: tk.Id,
// 			},
// 		},
// 		{
// 			Append: &wrapperspb.BoolValue{Value: false},
// 			Header: &corev3.HeaderValue{
// 				Key:   .UserID,
// 				Value: user.Id,
// 			},
// 		},
// 	}

// 	return headers
// }

// func (s *authZServer) checkIPAccess(ctx context.Context, xff string) error {
// 	if len(xff) == 0 {
// 		return nil
// 	}
// 	ip := strings.TrimSpace(strings.Split(xff, ",")[0])
// 	if !s.ip.IsAllowed(ctx, ip) {
// 		global.Log.Error("ip access denied", errIPAccessDenied)
// 		return errIPAccessDenied
// 	}
// 	return nil
// }
