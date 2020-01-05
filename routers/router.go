// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/wangdianwen/go-core/controllers"
	servicesCaptcha "github.com/wangdianwen/go-core/middlewares/captcha/controllers"
	servicesEmail "github.com/wangdianwen/go-core/middlewares/email/controllers"
)

func init() {
	// for framework core features
	core := beego.NewNamespace("/brunton",
		beego.NSNamespace("/user",
			beego.NSRouter("/login/:uname/:passwd", &controllers.UserController{}, "get:Login"),
			beego.NSRouter("/login/:barcode", &controllers.UserController{}, "get:Login"),
			beego.NSRouter("/logout", &controllers.UserController{}, "get:LogOut"),
			beego.NSRouter("/current_user", &controllers.UserController{}, "get:CurrentUser"),
			beego.NSRouter("/user_suggestion/:userName/:limit", &controllers.UserController{}, "get:UserSuggestion"),
			beego.NSRouter("/search", &controllers.UserController{}, "get:SearchUsers"),
			beego.NSRouter("/uploadAvatar", &controllers.UserController{}, "post:UploadAvatar"),
			beego.NSRouter("/saveUser", &controllers.UserController{}, "post:SaveUser"),
			beego.NSRouter("/editUser/:uid", &controllers.UserController{}, "patch:EditUser"),
			beego.NSRouter("/banUser/:uid", &controllers.UserController{}, "patch:BanUser"),
			beego.NSRouter("/releaseUser/:uid", &controllers.UserController{}, "patch:ReleaseUser"),
			beego.NSRouter("/register", &controllers.UserController{}, "post:RegisterUser"),
			beego.NSRouter("/register/verifyEmail", &controllers.UserController{}, "get:VerifyEmail"),
			beego.NSRouter("/register/resendVerifyEmail", &controllers.UserController{}, "post:ResendVerifyEmail"),
			beego.NSRouter("/profile/editUser", &controllers.UserController{}, "patch:EditProfile"),
		),
		beego.NSNamespace("/client",
			beego.NSRouter("/optionClients", &controllers.ClientController{}, "get:OptionClients"),
			beego.NSRouter("/client_suggestion/:clientName/:limit", &controllers.ClientController{}, "get:ClientSuggestion"),
			beego.NSRouter("/search", &controllers.ClientController{}, "get:ClientSearch"),
			beego.NSRouter("/saveClient", &controllers.ClientController{}, "post:SaveClient"),
			beego.NSRouter("/editClient/:cid", &controllers.ClientController{}, "patch:EditClient"),
			beego.NSRouter("/banClient/:cid", &controllers.ClientController{}, "patch:BanClient"),
			beego.NSRouter("/releaseClient/:cid", &controllers.ClientController{}, "patch:ReleaseClient"),
		),
		beego.NSNamespace("/role",
			beego.NSRouter("/optionRoles/:cid", &controllers.RoleController{}, "get:OptionRoles"),
			beego.NSRouter("/search", &controllers.RoleController{}, "get:RoleSearch"),
			beego.NSRouter("/saveRole", &controllers.RoleController{}, "post:SaveRole"),
			beego.NSRouter("/editRole/:rid", &controllers.RoleController{}, "patch:EditRole"),
			beego.NSRouter("/banRoles/:rid", &controllers.RoleController{}, "patch:BanRole"),
			beego.NSRouter("/releaseRoles/:rid", &controllers.RoleController{}, "patch:ReleaseRole"),
		),
		beego.NSNamespace("/privilege",
			beego.NSRouter("/search", &controllers.PrivilegeController{}, "get:PrivilegeSearch"),
			beego.NSRouter("/savePrivilege", &controllers.PrivilegeController{}, "post:SavePrivilege"),
			beego.NSRouter("/editPrivilege/:pid", &controllers.PrivilegeController{}, "patch:EditPrivilege"),
			beego.NSRouter("/banPrivilege/:pid", &controllers.PrivilegeController{}, "patch:BanPrivilege"),
			beego.NSRouter("/releasePrivilege/:pid", &controllers.PrivilegeController{}, "patch:ReleasePrivilege"),
		),
		beego.NSNamespace("/menu",
			beego.NSRouter("/optionMenus", &controllers.MenuController{}, "get:OptionMenus"),
			beego.NSRouter("/search", &controllers.MenuController{}, "get:MenuSearch"),
			beego.NSRouter("/saveMenu", &controllers.MenuController{}, "post:SaveMenu"),
			beego.NSRouter("/editMenu/:mid", &controllers.MenuController{}, "patch:EditMenu"),
			beego.NSRouter("/delMenu/:mid", &controllers.MenuController{}, "delete:DelMenu"),
			beego.NSRouter("/current", &controllers.MenuController{}, "get:CurrentMenu"),
			beego.NSRouter("/currentNames", &controllers.MenuController{}, "get:CurrentMenuNames"),
		),
		beego.NSNamespace("/app",
			beego.NSRouter("/optionApps", &controllers.AppController{}, "get:OptionApps"),
			beego.NSRouter("/search", &controllers.AppController{}, "get:AppsList"),
			beego.NSRouter("/saveApp", &controllers.AppController{}, "post:SaveApp"),
			beego.NSRouter("/editApp/:appid", &controllers.AppController{}, "patch:EditApp"),
			beego.NSRouter("/banApp/:appid", &controllers.AppController{}, "patch:BanApp"),
			beego.NSRouter("/releaseApp/:appid", &controllers.AppController{}, "patch:ReleaseApp"),
			beego.NSRouter("/validate_products", &controllers.AppController{}, "get:ValidateProducts"),
			beego.NSRouter("/start_free_trial_products", &controllers.AppController{}, "post:FreeTrialProducts"),
		),
	)
	beego.AddNamespace(core)

	// middleware
	servicesNamespace := beego.NewNamespace("/middlewares",
		beego.NSNamespace("/captcha",
			beego.NSRouter("/test", &servicesCaptcha.Captcha{}, "get:Test"),
		),
		beego.NSNamespace("/email",
			beego.NSRouter("/accountActivation", &servicesEmail.Email{}, "post:EmailVerification"),
			beego.NSRouter("/test", &servicesEmail.Email{}, "get:TestEmail"),
		),
	)
	beego.AddNamespace(servicesNamespace)
}
