package main

import (
	"fmt"
	"net/http"
	"time"
)

// 本来はDBや環境変数で管理
const (
	AdminUser    = "ikafly"
	AdminPass    = "already_set_true_password"
	CookieName   = "inside_session"
	SessionValue = "authenticated_user_shrimp" // 簡易的な固定トークン
)

func main() {
	// ルーティング
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", authMiddleware(indexHandler)) // 認証が必要なページ
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("Server starting at :8080...")
	http.ListenAndServe(":8080", nil)
}

// 認証チェック用ミドルウェア
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil || cookie.Value != SessionValue {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// ログイン処理
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// パスを "sites/login.html" に合わせる
		http.ServeFile(w, r, "sites/login.html")
		return
	}

	// POST処理
	user := r.FormValue("username")
	pass := r.FormValue("password")

	if user == AdminUser && pass == AdminPass {
		// Cookieをセット
		http.SetCookie(w, &http.Cookie{
			Name:     CookieName,
			Value:    SessionValue,
			Path:     "/",
			HttpOnly: true, // セキュリティ向上（JSからアクセス不可）
			Expires:  time.Now().Add(24 * time.Hour),
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		fmt.Fprintf(w, "認証失敗: ユーザー名かパスワードが違います")
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "sites/index.html")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   CookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
