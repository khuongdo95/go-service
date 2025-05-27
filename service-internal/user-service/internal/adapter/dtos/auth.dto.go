package dtos

type (
	SignUpReq struct {
		Username string
		Password string
		Email    string
		Name     string
	}

	SignUpRes struct {
		Token *TokenInfo
	}

	TokenInfo struct {
		UserId      string
		IdToken     string
		AccessToken string
	}

	SignInReq struct {
		Username string
		Password string
	}

	SignInRes struct {
		Token *TokenInfo
	}

	IPAccess struct {
		Id              int64
		IpAddress       string
		TcpStage        bool
		TcpLive         bool
		RabbitmqStage   bool
		RabbitmqLive    bool
		TranslationsRmq bool
		TranslationsTcp bool
		Active          bool
		Metadata        *ModifiedEntity
	}

	CreateIPReq struct {
		IpAddress       string
		TcpStage        bool
		TcpLive         bool
		RabbitmqStage   bool
		RabbitmqLive    bool
		TranslationsRmq bool
		TranslationsTcp bool
		Active          bool
	}

	UpdateIPReq struct {
		Id              int64
		IpAddress       *string
		TcpStage        *bool
		TcpLive         *bool
		RabbitmqStage   *bool
		RabbitmqLive    *bool
		TranslationsRmq *bool
		TranslationsTcp *bool
		Active          *bool
	}

	UpdateIPRes struct {
		Success bool
	}

	DeleteIPResponse struct {
		Success bool
	}

	ListIPReq struct {
		Pagination *PaginationReq
		OnlyActive bool
	}

	ListIPRes struct {
		Pagination *PaginationRes
		Data       []*IPAccess
	}

	ListUserReq struct {
		Pagination *PaginationReq `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
	}
	ListUserRes struct {
		Users      []*UserRes
		Pagination *PaginationRes
	}

	UserRes struct {
		Id        string
		Name      string
		Email     string
		Username  string
		IpAddress string
		Metadata  *ModifiedEntity
	}

	GetIPReq struct {
		Id int64
	}

	DeleteIPReq struct {
		Id int64
	}

	CreateUserReq struct {
		Name      string
		Email     string
		Username  string
		Password  string
		IpAddress *string
	}

	UpdateUserReq struct {
		Id        string
		Name      *string
		Email     *string
		IpAddress *string
	}
)
