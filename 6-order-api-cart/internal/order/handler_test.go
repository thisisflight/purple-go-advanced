package order

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"purple/links/configs"
	"purple/links/internal/product"
	"purple/links/internal/user"
	"purple/links/pkg/db"
)

const TEST_SECRET_KEY = "/2+XnmJGz1j3ehIVI/5P9kl+CghrE3DcS7rnT+qar5w="

func getTestDB() *db.DB {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := &configs.Config{
		Db: configs.DBConfig{
			Dsn: os.Getenv("TEST_DSN"),
		},
	}

	return db.NewDB(cfg)
}

func prepareTestData(t *testing.T, db *db.DB) (userID uint, productID uint) {
	testUser := &user.User{
		PhoneNumber: "89206016463",
		Name:        "Test User",
	}
	err := db.Create(testUser).Error
	require.NoError(t, err)
	userID = testUser.ID

	testProduct := &product.Product{
		Name:  "Test Product",
		Price: decimal.NewFromFloat(100.0),
	}
	err = db.Create(testProduct).Error
	require.NoError(t, err)
	productID = testProduct.ID

	return userID, productID
}

func cleanupTestData(db *db.DB, userID uint, productID uint, orderID uint) {
	db.Unscoped().Where("order_id = ?", orderID).Delete(&OrderItem{})
	db.Unscoped().Where("id = ?", orderID).Delete(&Order{})
	db.Unscoped().Where("user_id = ?", userID).Delete(&user.User{})
	db.Unscoped().Where("id = ?", productID).Delete(&product.Product{})
}

func generateTestToken(phoneNumber string, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["phoneNumber"] = phoneNumber
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	return token.SignedString([]byte(secret))
}

func TestOrderHandler_Create_E2E(t *testing.T) {
	testDB := getTestDB()

	cfg := &configs.Config{
		Auth: configs.AuthConfig{
			Secret: TEST_SECRET_KEY,
		},
	}

	userID, productID := prepareTestData(t, testDB)

	var orderID uint
	defer cleanupTestData(testDB, userID, productID, orderID)

	userRepo := user.NewUserRepository(testDB)
	orderRepo := NewOrderRepository(testDB)
	orderValidator := NewOrderValidator(testDB)

	orderService := NewOrderService(OrderServiceDeps{
		OrderRepository: orderRepo,
		OrderValidator:  orderValidator,
		UserRepository:  userRepo,
	})

	router := http.NewServeMux()
	deps := OrderHandlerDeps{
		OrderService: orderService,
		Config:       cfg,
	}
	NewOrderHandler(router, deps)

	orderRequest := OrderCreateRequest{
		Notes: "Test order notes",
		Items: []OrderItemRequest{
			{ProductID: int64(productID), Quantity: 2},
		},
	}

	requestBody, err := json.Marshal(orderRequest)
	require.NoError(t, err)

	token, err := generateTestToken("89206016463", TEST_SECRET_KEY)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/order", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response OrderCreateResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, orderRequest.Notes, response.Notes)
}
