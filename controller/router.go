package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/Asliddin3/energy-maximum/config"
	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/file"
	"github.com/Asliddin3/energy-maximum/pkg/hash"
	"github.com/Asliddin3/energy-maximum/pkg/humanizer"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/Asliddin3/energy-maximum/pkg/sms"
	"github.com/Asliddin3/energy-maximum/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db           *gorm.DB
	filesService *file.FilesService
	log          *logger.MyLogger
	cfg          *config.Config
	hash         *hash.Hash
	humanizer    *humanizer.ManagerHumanizer
}

func NewHandler(db *gorm.DB, log *logger.MyLogger, cfg config.Config, hash *hash.Hash, hum *humanizer.ManagerHumanizer) *Handler {
	return &Handler{
		db:           db,
		log:          log,
		filesService: file.NewFilesService(cfg),
		cfg:          &cfg,
		hash:         hash,
		humanizer:    hum,
	}
}

func (h *Handler) Init(server *gin.Engine) {
	sms := sms.NewSmsSender()
	api := server.Group("api")
	{
		h.NewAnalogController(api)
		h.NewPagesController(api)
		h.NewModuleController(api)
		h.NewRolesController(api)
		h.NewPublicOfferController(api)
		h.NewCategoryController(api)
		h.NewServiceController(api)
		h.NewCountryController(api)
		h.NewAboutController(api)
		h.NewCustomerController(api, sms)
		h.NewContactController(api)
		h.NewBannerController(api)
		h.NewProductController(api)
		h.NewParameterController(api)
		h.NewAdminController(api)
		h.NewBrandController(api)
		h.NewNewsController(api)
		h.NewVacancyController(api)
		h.NewOrderController(api)
	}
	server.StaticFS("/public/", http.Dir(h.cfg.StaticFilePath))
}
func (h *Handler) GetAdmin(c *gin.Context) *models.AdminMetadata {
	admin := c.MustGet("admin").(models.AdminMetadata)
	return &admin
}
func (h *Handler) GetCustomer(c *gin.Context) *models.CustomerMetadata {
	customer := c.MustGet("customer").(models.CustomerMetadata)
	return &customer
}

func (h *Handler) DeserializeCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		cookie, err := ctx.Cookie("access_token")
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)
		if len(fields) != 0 {
			if fields[0] == "Bearer" {
				access_token = fields[1]
			} else if fields[0] != "" {
				access_token = fields[0]
			}
		} else if err == nil {
			access_token = cookie
		}
		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}
		sub, err := utils.ValidateToken(access_token, h.cfg.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		user := models.CustomerMetadata{
			Id: int(sub.(float64)),
		}

		err = h.db.Model(&models.Customer{}).Where("id=?", user.Id).UpdateColumn("last_visit", time.Now()).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Set("customer", user)
		// initializers.RateLimit(ctx)
		ctx.Next()
	}
}

// DeserializeUser this method will getting user id and branch.Use it for separate data by branches.
func (h *Handler) DeserializeAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		cookie, err := ctx.Cookie("access_token")
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)
		if len(fields) != 0 {
			if fields[0] == "Bearer" {
				access_token = fields[1]
			} else if fields[0] != "" {
				access_token = fields[0]
			}
		} else if err == nil {
			access_token = cookie
		}
		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}
		sub, err := utils.ValidateToken(access_token, h.cfg.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		user := models.AdminMetadata{
			Id: int(sub.(float64)),
		}

		rows := h.db.Model(&models.Admins{}).Where("id=? AND deleted_at IS NULL AND is_active=true", sub).UpdateColumn("last_visit", time.Now())
		if rows.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, rows.Error.Error())
			return
		} else if rows.RowsAffected == 0 {
			newResponse(ctx, http.StatusUnauthorized, "not found admin")
			return
		}
		ctx.Set("admin", user)
		// initializers.RateLimit(ctx)
		ctx.Next()
	}
}
