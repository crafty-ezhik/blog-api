package auth

import (
	mock_auth "github.com/crafty-ezhik/blog-api/internal/auth/mock"
	"github.com/crafty-ezhik/blog-api/internal/models"
	"go.uber.org/mock/gomock"
	"testing"
)

type mockBehavior func(s *mock_auth.MockAuthService, user *models.User)

type TestTable struct {
	// name - имя теста
	name string

	// inputBody - тело запроса
	inputBody string

	//inputUser - структура пользователя
	inputUser *models.User

	// mockBehavior - функция mockBehavior
	mockBehavior mockBehavior

	// expectedStatus - Ожидаемый статус код
	expectedStatus int

	// expectedRequestBody - Ожидаемое тело ответа
	expectedRequestBody string
}

func TestAuthHandler_Register(t *testing.T) {

	// Определяем тестовую таблицу
	testTable := []TestTable{
		{
			name:      "OK",
			inputBody: `{"email": "email@email.com", "password": "password", "name": "Alex", "age": 32}`,
			inputUser: &models.User{Email: "email@email.com", Password: "password", Name: "Alex", Age: 32},
			mockBehavior: func(s *mock_auth.MockAuthService, user *models.User) {
				s.EXPECT().Register(user).Return(true, nil)
			},
			expectedStatus:      201,
			expectedRequestBody: `{"message": "You have successfully registered","success": true}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Для каждого теста необходимо создавать контроллер и финишировать при завершении теста
			c := gomock.NewController(t)
			defer c.Finish()

			// Создадим mock сервиса
			mockAuth := mock_auth.NewMockAuthService(c)

			// Вызываем нашу функцию и передаем структуру пользователя
			testCase.mockBehavior(mockAuth, testCase.inputUser)

			service := &AuthServiceimpl{}
		})
	}
}
