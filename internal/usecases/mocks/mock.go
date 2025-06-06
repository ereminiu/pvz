// Code generated by MockGen. DO NOT EDIT.
// Source: usecases.go
//
// Generated by this command:
//
//	mockgen -source=usecases.go -destination=mocks/mock.go
//

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	context "context"
	reflect "reflect"

	entities "github.com/ereminiu/pvz/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockOrderRepository is a mock of OrderRepository interface.
type MockOrderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryMockRecorder
	isgomock struct{}
}

// MockOrderRepositoryMockRecorder is the mock recorder for MockOrderRepository.
type MockOrderRepositoryMockRecorder struct {
	mock *MockOrderRepository
}

// NewMockOrderRepository creates a new mock instance.
func NewMockOrderRepository(ctrl *gomock.Controller) *MockOrderRepository {
	mock := &MockOrderRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepository) EXPECT() *MockOrderRepositoryMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockOrderRepository) AddOrder(ctx context.Context, order *entities.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockOrderRepositoryMockRecorder) AddOrder(ctx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockOrderRepository)(nil).AddOrder), ctx, order)
}

// RemoveOrder mocks base method.
func (m *MockOrderRepository)RemoveOrder(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveOrder", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveOrder indicates an expected call of RemoveOrder.
func (mr *MockOrderRepositoryMockRecorder) RemoveOrder(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOrder", reflect.TypeOf((*MockOrderRepository)(nil).RemoveOrder), ctx, id)
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
	isgomock struct{}
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// GetList mocks base method.
func (m *MockUserRepository) GetList(ctx context.Context, userID, lastN int, located bool, pattern map[string]string) ([]*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx, userID, lastN, located, pattern)
	ret0, _ := ret[0].([]*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockUserRepositoryMockRecorder) GetList(ctx, userID, lastN, located, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockUserRepository)(nil).GetList), ctx, userID, lastN, located, pattern)
}

// RefundOrders mocks base method.
func (m *MockUserRepository) RefundOrders(ctx context.Context, userID int, orderIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefundOrders", ctx, userID, orderIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// RefundOrders indicates an expected call of RefundOrders.
func (mr *MockUserRepositoryMockRecorder) RefundOrders(ctx, userID, orderIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefundOrders", reflect.TypeOf((*MockUserRepository)(nil).RefundOrders), ctx, userID, orderIDs)
}

// ReturnOrders mocks base method.
func (m *MockUserRepository) ReturnOrders(ctx context.Context, userID int, orderIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReturnOrders", ctx, userID, orderIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReturnOrders indicates an expected call of ReturnOrders.
func (mr *MockUserRepositoryMockRecorder) ReturnOrders(ctx, userID, orderIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReturnOrders", reflect.TypeOf((*MockUserRepository)(nil).ReturnOrders), ctx, userID, orderIDs)
}

// MockPVZRepository is a mock of PVZRepository interface.
type MockPVZRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPVZRepositoryMockRecorder
	isgomock struct{}
}

// MockPVZRepositoryMockRecorder is the mock recorder for MockPVZRepository.
type MockPVZRepositoryMockRecorder struct {
	mock *MockPVZRepository
}

// NewMockPVZRepository creates a new mock instance.
func NewMockPVZRepository(ctrl *gomock.Controller) *MockPVZRepository {
	mock := &MockPVZRepository{ctrl: ctrl}
	mock.recorder = &MockPVZRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPVZRepository) EXPECT() *MockPVZRepositoryMockRecorder {
	return m.recorder
}

// GetHistory mocks base method.
func (m *MockPVZRepository) GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistory", ctx, page, limit, orderBy, pattern)
	ret0, _ := ret[0].([]*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHistory indicates an expected call of GetHistory.
func (mr *MockPVZRepositoryMockRecorder) GetHistory(ctx, page, limit, orderBy, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistory", reflect.TypeOf((*MockPVZRepository)(nil).GetHistory), ctx, page, limit, orderBy, pattern)
}

// GetRefunds mocks base method.
func (m *MockPVZRepository) GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefunds", ctx, page, limit, orderBy, pattern)
	ret0, _ := ret[0].([]*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefunds indicates an expected call of GetRefunds.
func (mr *MockPVZRepositoryMockRecorder) GetRefunds(ctx, page, limit, orderBy, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefunds", reflect.TypeOf((*MockPVZRepository)(nil).GetRefunds), ctx, page, limit, orderBy, pattern)
}

// MockAuthRepository is a mock of AuthRepository interface.
type MockAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepositoryMockRecorder
	isgomock struct{}
}

// MockAuthRepositoryMockRecorder is the mock recorder for MockAuthRepository.
type MockAuthRepositoryMockRecorder struct {
	mock *MockAuthRepository
}

// NewMockAuthRepository creates a new mock instance.
func NewMockAuthRepository(ctrl *gomock.Controller) *MockAuthRepository {
	mock := &MockAuthRepository{ctrl: ctrl}
	mock.recorder = &MockAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthRepository) EXPECT() *MockAuthRepositoryMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockAuthRepository) CreateUser(ctx context.Context, username, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, username, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthRepositoryMockRecorder) CreateUser(ctx, username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthRepository)(nil).CreateUser), ctx, username, password)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
	isgomock struct{}
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockRepository) AddOrder(ctx context.Context, order *entities.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockRepositoryMockRecorder) AddOrder(ctx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockRepository)(nil).AddOrder), ctx, order)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(ctx context.Context, username, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, username, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(ctx, username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), ctx, username, password)
}

// GetHistory mocks base method.
func (m *MockRepository) GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistory", ctx, page, limit, orderBy, pattern)
	ret0, _ := ret[0].([]*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHistory indicates an expected call of GetHistory.
func (mr *MockRepositoryMockRecorder) GetHistory(ctx, page, limit, orderBy, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistory", reflect.TypeOf((*MockRepository)(nil).GetHistory), ctx, page, limit, orderBy, pattern)
}

// GetList mocks base method.
func (m *MockRepository) GetList(ctx context.Context, userID, lastN int, located bool, pattern map[string]string) ([]*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx, userID, lastN, located, pattern)
	ret0, _ := ret[0].([]*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockRepositoryMockRecorder) GetList(ctx, userID, lastN, located, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockRepository)(nil).GetList), ctx, userID, lastN, located, pattern)
}

// GetRefunds mocks base method.
func (m *MockRepository) GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefunds", ctx, page, limit, orderBy, pattern)
	ret0, _ := ret[0].([]*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefunds indicates an expected call of GetRefunds.
func (mr *MockRepositoryMockRecorder) GetRefunds(ctx, page, limit, orderBy, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefunds", reflect.TypeOf((*MockRepository)(nil).GetRefunds), ctx, page, limit, orderBy, pattern)
}

// RefundOrders mocks base method.
func (m *MockRepository) RefundOrders(ctx context.Context, userID int, orderIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefundOrders", ctx, userID, orderIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// RefundOrders indicates an expected call of RefundOrders.
func (mr *MockRepositoryMockRecorder) RefundOrders(ctx, userID, orderIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefundOrders", reflect.TypeOf((*MockRepository)(nil).RefundOrders), ctx, userID, orderIDs)
}

// RemoveOrder mocks base method.
func (m *MockRepository)RemoveOrder(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveOrder", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveOrder indicates an expected call of RemoveOrder.
func (mr *MockRepositoryMockRecorder) RemoveOrder(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOrder", reflect.TypeOf((*MockRepository)(nil).RemoveOrder), ctx, id)
}

// ReturnOrders mocks base method.
func (m *MockRepository) ReturnOrders(ctx context.Context, userID int, orderIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReturnOrders", ctx, userID, orderIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReturnOrders indicates an expected call of ReturnOrders.
func (mr *MockRepositoryMockRecorder) ReturnOrders(ctx, userID, orderIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReturnOrders", reflect.TypeOf((*MockRepository)(nil).ReturnOrders), ctx, userID, orderIDs)
}

// MockOrderUsecases is a mock of OrderUsecases interface.
type MockOrderUsecases struct {
	ctrl     *gomock.Controller
	recorder *MockOrderUsecasesMockRecorder
	isgomock struct{}
}

// MockOrderUsecasesMockRecorder is the mock recorder for MockOrderUsecases.
type MockOrderUsecasesMockRecorder struct {
	mock *MockOrderUsecases
}

// NewMockOrderUsecases creates a new mock instance.
func NewMockOrderUsecases(ctrl *gomock.Controller) *MockOrderUsecases {
	mock := &MockOrderUsecases{ctrl: ctrl}
	mock.recorder = &MockOrderUsecasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderUsecases) EXPECT() *MockOrderUsecasesMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockOrderUsecases) AddOrder(ctx context.Context, order *entities.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockOrderUsecasesMockRecorder) AddOrder(ctx, order any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockOrderUsecases)(nil).AddOrder), ctx, order)
}

// RemoveOrder mocks base method.
func (m *MockOrderUsecases)RemoveOrder(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveOrder", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveOrder indicates an expected call of RemoveOrder.
func (mr *MockOrderUsecasesMockRecorder) RemoveOrder(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOrder", reflect.TypeOf((*MockOrderUsecases)(nil).RemoveOrder), ctx, id)
}

// MockUserUsecases is a mock of UserUsecases interface.
type MockUserUsecases struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecasesMockRecorder
	isgomock struct{}
}

// MockUserUsecasesMockRecorder is the mock recorder for MockUserUsecases.
type MockUserUsecasesMockRecorder struct {
	mock *MockUserUsecases
}

// NewMockUserUsecases creates a new mock instance.
func NewMockUserUsecases(ctrl *gomock.Controller) *MockUserUsecases {
	mock := &MockUserUsecases{ctrl: ctrl}
	mock.recorder = &MockUserUsecasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecases) EXPECT() *MockUserUsecasesMockRecorder {
	return m.recorder
}

// GetList mocks base method.
func (m *MockUserUsecases) GetList(ctx context.Context, userID, lastN int, located bool, pattern map[string]string) ([]*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx, userID, lastN, located, pattern)
	ret0, _ := ret[0].([]*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockUserUsecasesMockRecorder) GetList(ctx, userID, lastN, located, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockUserUsecases)(nil).GetList), ctx, userID, lastN, located, pattern)
}

// RefundOrders mocks base method.
func (m *MockUserUsecases) RefundOrders(ctx context.Context, userID int, orderIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefundOrders", ctx, userID, orderIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// RefundOrders indicates an expected call of RefundOrders.
func (mr *MockUserUsecasesMockRecorder) RefundOrders(ctx, userID, orderIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefundOrders", reflect.TypeOf((*MockUserUsecases)(nil).RefundOrders), ctx, userID, orderIDs)
}

// ReturnOrders mocks base method.
func (m *MockUserUsecases) ReturnOrders(ctx context.Context, userID int, orderIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReturnOrders", ctx, userID, orderIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReturnOrders indicates an expected call of ReturnOrders.
func (mr *MockUserUsecasesMockRecorder) ReturnOrders(ctx, userID, orderIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReturnOrders", reflect.TypeOf((*MockUserUsecases)(nil).ReturnOrders), ctx, userID, orderIDs)
}

// MockPVZUsecases is a mock of PVZUsecases interface.
type MockPVZUsecases struct {
	ctrl     *gomock.Controller
	recorder *MockPVZUsecasesMockRecorder
	isgomock struct{}
}

// MockPVZUsecasesMockRecorder is the mock recorder for MockPVZUsecases.
type MockPVZUsecasesMockRecorder struct {
	mock *MockPVZUsecases
}

// NewMockPVZUsecases creates a new mock instance.
func NewMockPVZUsecases(ctrl *gomock.Controller) *MockPVZUsecases {
	mock := &MockPVZUsecases{ctrl: ctrl}
	mock.recorder = &MockPVZUsecasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPVZUsecases) EXPECT() *MockPVZUsecasesMockRecorder {
	return m.recorder
}

// GetHistory mocks base method.
func (m *MockPVZUsecases) GetHistory(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistory", ctx, page, limit, orderBy, pattern)
	ret0, _ := ret[0].([]*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHistory indicates an expected call of GetHistory.
func (mr *MockPVZUsecasesMockRecorder) GetHistory(ctx, page, limit, orderBy, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistory", reflect.TypeOf((*MockPVZUsecases)(nil).GetHistory), ctx, page, limit, orderBy, pattern)
}

// GetRefunds mocks base method.
func (m *MockPVZUsecases) GetRefunds(ctx context.Context, page, limit int, orderBy string, pattern map[string]string) ([]*entities.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefunds", ctx, page, limit, orderBy, pattern)
	ret0, _ := ret[0].([]*entities.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefunds indicates an expected call of GetRefunds.
func (mr *MockPVZUsecasesMockRecorder) GetRefunds(ctx, page, limit, orderBy, pattern any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefunds", reflect.TypeOf((*MockPVZUsecases)(nil).GetRefunds), ctx, page, limit, orderBy, pattern)
}

// MockAuthUsecases is a mock of AuthUsecases interface.
type MockAuthUsecases struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUsecasesMockRecorder
	isgomock struct{}
}

// MockAuthUsecasesMockRecorder is the mock recorder for MockAuthUsecases.
type MockAuthUsecasesMockRecorder struct {
	mock *MockAuthUsecases
}

// NewMockAuthUsecases creates a new mock instance.
func NewMockAuthUsecases(ctrl *gomock.Controller) *MockAuthUsecases {
	mock := &MockAuthUsecases{ctrl: ctrl}
	mock.recorder = &MockAuthUsecasesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUsecases) EXPECT() *MockAuthUsecasesMockRecorder {
	return m.recorder
}

// SignIn mocks base method.
func (m *MockAuthUsecases) SignIn(ctx context.Context, username, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", ctx, username, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthUsecasesMockRecorder) SignIn(ctx, username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuthUsecases)(nil).SignIn), ctx, username, password)
}
