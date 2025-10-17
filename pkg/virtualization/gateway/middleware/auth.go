package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	authv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Authenticator exposes Gin middleware hooking into Kubernetes authentication.
type Authenticator interface {
	GinMiddleware() gin.HandlerFunc
}

type jwtAuthenticator struct {
	client *kubernetes.Clientset
}

// NewJWTImpersonationAuthenticator returns an Authenticator using TokenReview and impersonation headers.
func NewJWTImpersonationAuthenticator(cfg *rest.Config) Authenticator {
	return &jwtAuthenticator{client: kubernetes.NewForConfigOrDie(cfg)}
}

func (a *jwtAuthenticator) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c.GetHeader("Authorization"))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		review, err := a.client.AuthenticationV1().TokenReviews().Create(
			c,
			&authv1.TokenReview{Spec: authv1.TokenReviewSpec{Token: token}},
			metav1.CreateOptions{},
		)
		if err != nil || !review.Status.Authenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token rejected"})
			return
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), userKey{}, review.Status.User))
		c.Request.Header.Set("Impersonate-User", review.Status.User.Username)
		if review.Status.User.UID != "" {
			c.Request.Header.Set("Impersonate-Uid", review.Status.User.UID)
		}
		if len(review.Status.User.Groups) > 0 {
			c.Request.Header.Del("Impersonate-Group")
			for _, group := range review.Status.User.Groups {
				c.Request.Header.Add("Impersonate-Group", group)
			}
		}
		c.Next()
	}
}

func extractBearer(header string) string {
	if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return ""
	}
	return strings.TrimSpace(header[7:])
}

type userKey struct{}

// UserFrom retrieves the authentication user info.
func UserFrom(ctx context.Context) *authv1.UserInfo {
	if v, ok := ctx.Value(userKey{}).(authv1.UserInfo); ok {
		return &v
	}
	return nil
}
