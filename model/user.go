package model

// User is data of user in Qiita.
type User struct {
	Description       string `json:"description"`         // 自己紹介文
	FacebookID        string `json:"facebook_id"`         // Facebook ID
	FolloweesCount    int    `json:"followees_count"`     // このユーザがフォローしているユーザの数
	FollowersCount    int    `json:"followers_count"`     // このユーザをフォローしているユーザの数
	GithubLoginName   string `json:"github_login_name"`   // GitHub ID
	ID                string `json:"id"`                  // ユーザID
	ItemsCount        int    `json:"items_count"`         // このユーザが qiita.com 上で公開している投稿の数 (Qiita:Teamでの投稿数は含まれません)
	LinkedinID        string `json:"linkedin_id"`         // LinkedIn ID
	Location          string `json:"location"`            // 居住地
	Name              string `json:"name"`                // 設定している名前
	Organization      string `json:"organization"`        // 所属している組織
	PermanentID       int    `json:"permanent_id"`        // ユーザごとに割り当てられる整数のID
	ProfileImageURL   string `json:"profile_image_url"`   // 設定しているプロフィール画像のURL
	TwitterScreenName string `json:"twitter_screen_name"` // Twitterのスクリーンネーム
	WebsiteURL        string `json:"website_url"`         // 設定しているWebサイトのURL
}
