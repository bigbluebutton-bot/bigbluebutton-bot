package bbb

// For users and current-user
type User struct {
	MeetingID               string            `json:"meetingId"`
	UserID                  string            `json:"userId"`
	ClientType              string            `json:"clientType"`
	Validated               bool              `json:"validated"`
	Left                    bool              `json:"left"`
	Approved                bool              `json:"approved"`
	AuthTokenValidatedTime  int64             `json:"authTokenValidatedTime"`
	InactivityCheck         bool              `json:"inactivityCheck"`
	LoginTime               int64             `json:"loginTime"`
	Authed                  bool              `json:"authed"`
	Avatar                  string            `json:"avatar"`
	BreakoutProps           UserBreakoutProps `json:"breakoutProps"`
	Color                   string            `json:"color"`
	EffectiveConnectionType interface{}       `json:"effectiveConnectionType"`
	Emoji                   string            `json:"emoji"`
	ExtID                   string            `json:"extId"`
	Guest                   bool              `json:"guest"`
	GuestStatus             string            `json:"guestStatus"`
	IntID                   string            `json:"intId"`
	Locked                  bool              `json:"locked"`
	LoggedOut               bool              `json:"loggedOut"`
	Mobile                  bool              `json:"mobile"`
	Name                    string            `json:"name"`
	Pin                     bool              `json:"pin"`
	Presenter               bool              `json:"presenter"`
	ResponseDelay           int               `json:"responseDelay"`
	Role                    string            `json:"role"`
	SortName                string            `json:"sortName"`
}
type UserBreakoutProps struct {
	IsBreakoutUser bool   `json:"isBreakoutUser"`
	ParentID       string `json:"parentId"`
}