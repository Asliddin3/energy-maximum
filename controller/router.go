package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Asliddin3/energy-maximum/config"
	"github.com/Asliddin3/energy-maximum/models"
	"github.com/Asliddin3/energy-maximum/pkg/file"
	"github.com/Asliddin3/energy-maximum/pkg/hash"
	"github.com/Asliddin3/energy-maximum/pkg/logger"
	"github.com/Asliddin3/energy-maximum/pkg/sms"
	"github.com/Asliddin3/energy-maximum/pkg/utils"

	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db           *gorm.DB
	BigCache     *bigcache.BigCache
	filesService *file.FilesService
	log          *logger.MyLogger
	cfg          *config.Config
	hash         *hash.Hash
}

func NewHandler(db *gorm.DB, log *logger.MyLogger, cache *bigcache.BigCache, cfg config.Config, hash *hash.Hash) *Handler {
	return &Handler{
		db:           db,
		log:          log,
		filesService: file.NewFilesService(cfg),
		BigCache:     cache,
		cfg:          &cfg,
		hash:         hash,
	}
}

func (h *Handler) Init(server *gin.Engine) {
	sms := sms.NewSmsSender()
	api := server.Group("api")
	{
		h.NewAnalogController(api)
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
		tokenKey, err := utils.ValidateToken(access_token, h.cfg.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		sub := strings.Split(tokenKey.(string), ":")[1]
		key := fmt.Sprintf("customer:%v", sub)
		var user models.CustomerMetadata
		userByte, err := h.BigCache.Get(key)
		if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		if err == bigcache.ErrEntryNotFound {
			result := h.db.Debug().Model(&models.Customer{}).Where("id=?", fmt.Sprint(sub)).
				Select("id").First(&user)
			if result.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no longer exists"})
				return
			}
			fmt.Println("gotten user from db", user)
			res, err := json.Marshal(&user)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
				return
			}

			err = h.BigCache.Set(key, res)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
				return
			}
		} else {
			err = json.Unmarshal(userByte, &user)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				return
			}
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
		key := fmt.Sprintf("admin:%v", sub)
		var user models.AdminMetadata
		userByte, err := h.BigCache.Get(key)
		if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		if err == bigcache.ErrEntryNotFound {
			result := h.db.Debug().Model(&models.Admins{}).Where("id=? AND deleted_at IS NULL AND is_active=true", fmt.Sprint(sub)).
				Select("id,is_superuser").First(&user)
			if result.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no longer exists"})
				return
			}
			fmt.Println("gotten user from db", user)
			res, err := json.Marshal(&user)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
				return
			}

			err = h.BigCache.Set(key, res)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
				return
			}
		} else {
			err = json.Unmarshal(userByte, &user)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
		err = h.db.Model(&models.Admins{}).Where("id=?", user.Id).UpdateColumn("last_visit", time.Now()).Error
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.Set("admin", user)
		// initializers.RateLimit(ctx)
		ctx.Next()
	}
}
