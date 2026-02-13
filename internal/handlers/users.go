package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MrTomatePNG/projeto-m/internal/auth"
	"github.com/MrTomatePNG/projeto-m/internal/database"
	"github.com/MrTomatePNG/projeto-m/internal/middleware"
	"github.com/MrTomatePNG/projeto-m/internal/services"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	s        *services.UserService
	validate *validator.Validate
	jwtm     *auth.JWTManager
}

func NewUserHandler(db *database.Queries, jwtm *auth.JWTManager) *UserHandler {
	s := services.NewUserService(db)
	return &UserHandler{
		s:        s,
		validate: validator.New(),
		jwtm:     jwtm,
	}
}

// Helper para responder JSON
func (h *UserHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user struct {
			Username string `json:"username" validate:"required,min=3,max=50"`
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			h.respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid JSON format",
			})
			return
		}

		if err := h.validate.Struct(user); err != nil {
			errors := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				errors[err.Field()] = err.Tag()
			}
			h.respondJSON(w, http.StatusBadRequest, errors)
			return
		}

		userDTO, err := h.s.RegisterUser(user.Username, user.Email, user.Password)
		if err != nil {
			// Aqui agora funciona!
			if errors.Is(err, services.ErrUsernameTaken) {
				h.respondJSON(w, http.StatusConflict, map[string]string{
					"error": "Username já existe",
				})
				return
			}
			if errors.Is(err, services.ErrEmailTaken) {
				h.respondJSON(w, http.StatusConflict, map[string]string{
					"error": "Email já existe",
				})
				return
			}
			h.respondJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Erro ao criar usuário",
			})
			return
		}

		h.respondJSON(w, http.StatusCreated, userDTO)
	}
}

func (h *UserHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var login struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8"`
		}

		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			h.respondJSON(w, http.StatusBadRequest, map[string]string{
				"error": "Invalid JSON format",
			})
			return
		}

		if err := h.validate.Struct(&login); err != nil {
			errors := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				errors[err.Field()] = err.Tag()
			}
			h.respondJSON(w, http.StatusBadRequest, errors)
			return
		}

		// Buscar user no banco
		user, err := h.s.GetUserByEmail(login.Email)
		if err != nil {
			h.respondJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Email ou senha inválida",
			})
			return
		}

		// Comparar: hash do banco com a senha digitada (NÃO gerar novo hash!)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
			h.respondJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Email ou senha inválida",
			})
			return
		}

		token, err := h.jwtm.Generate(user.ID)
		if err != nil {
			h.respondJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Erro ao gerar token",
			})
			return
		}

		// Login bem-sucedido
		h.respondJSON(w, http.StatusOK, map[string]interface{}{
			"message":      "Login bem-sucedido",
			"access_token": token,
		})
	}
}

func (h *UserHandler) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			h.respondJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized",
			})
			return
		}

		user, err := h.s.GetUserByID(userID)
		if err != nil {
			h.respondJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Erro ao buscar usuário",
			})
			return
		}

		h.respondJSON(w, http.StatusOK, map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"bio":      user.Bio.String,
		})
	}
}
