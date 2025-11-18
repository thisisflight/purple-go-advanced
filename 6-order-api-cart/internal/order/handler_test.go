package order

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"purple/links/configs"
	"purple/links/internal/user"
	"purple/links/pkg/db"
)

const TEST_SECRET_KEY = "/2+XnmJGz1j3ehIVI/5P9kl+CghrE3DcS7rnT+qar5w="

func setupTestDB(t *testing.T) (*db.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	require.NoError(t, err)

	return &db.DB{DB: gormDB}, mock
}

func TestOrderHandler_Create_E2E(t *testing.T) {
	db, mock := setupTestDB(t)

	cfg := &configs.Config{
		Auth: configs.AuthConfig{
			Secret: TEST_SECRET_KEY,
		},
	}

	mock.ExpectQuery(`SELECT count\(\*\) FROM "product" WHERE id IN .*`).
		WithArgs(100).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT \* FROM "user" WHERE phone_number = .*`).
		WithArgs("89206016463", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "phone_number", "name", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "89206016463", "Test User", time.Now(), time.Now(), nil))

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "order"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 0, "Test order notes").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectQuery(`SELECT \* FROM "product" WHERE .*id.*`).
		WithArgs(100, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "created_at", "updated_at", "deleted_at"}).
			AddRow(100, "Product 1", 100.0, time.Now(), time.Now(), nil))
	mock.ExpectQuery(`INSERT INTO "order_item"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			1,
			100,
			2,
			"100",
			"0",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	userRepo := user.NewUserRepository(db)
	orderRepo := NewOrderRepository(db)
	orderValidator := NewOrderValidator(db)

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
			{ProductID: 100, Quantity: 2},
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

	if rr.Code != http.StatusCreated {
		t.Logf("Response body: %s", rr.Body.String())
	}
	assert.Equal(t, http.StatusCreated, rr.Code)

	if rr.Code == http.StatusCreated {
		var response OrderCreateResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, uint(1), response.UserID)
		assert.Equal(t, orderRequest.Notes, response.Notes)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func generateTestToken(phoneNumber string, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["phoneNumber"] = phoneNumber
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	return token.SignedString([]byte(secret))
}
