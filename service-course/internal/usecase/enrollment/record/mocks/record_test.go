package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/barashF/lms/service-course/internal/logger"
	"github.com/barashF/lms/service-course/internal/model"
	"github.com/barashF/lms/service-course/internal/usecase/enrollment/record"
	"github.com/barashF/lms/service-course/proto/orders"
	payments "github.com/barashF/lms/service-course/proto/payment"
	"github.com/barashF/lms/service-course/proto/profiles"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	ID1 = uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b")
	ID2 = uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c8b")
)

var (
	enroll1  = &model.Enrollment{ID: ID1, CourseID: ID1, UserID: ID1}
	course1  = &model.Course{ID: ID1, Price: 100, Type: model.CoursePaid}
	order1   = &orders.Order{Id: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", CourseId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b"}
	payment1 = &payments.Payment{Id: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", OrderId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", Status: "completed"}
	profile1 = &profiles.Profile{UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", Tier: "premium"}
)

var (
	enroll2  = &model.Enrollment{ID: ID1, CourseID: ID1, UserID: ID1}
	course2  = &model.Course{ID: ID1, Price: 100, Type: model.CoursePaid}
	order2   = &orders.Order{Id: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", CourseId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b"}
	payment2 = &payments.Payment{Id: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", OrderId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", Status: "failed"}
	profile2 = &profiles.Profile{UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", Tier: "premium"}
)

var (
	enroll3  = &model.Enrollment{ID: ID1, CourseID: ID1, UserID: ID1}
	course3  = &model.Course{ID: ID1, Price: 100, Type: model.CourseCorporate}
	order3   = &orders.Order{Id: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", CourseId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b"}
	payment3 = &payments.Payment{Id: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", OrderId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", Status: "failed"}
	profile3 = &profiles.Profile{UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", Tier: "premium"}
)

var course4 = &model.Course{ID: ID1, Price: 100, Type: model.CourseFree}

func TestRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, model.ErrNoEnrollmentFound)
	enrollRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"), nil)
	courseRepo := NewMockcourseRepository(ctrl)
	courseRepo.EXPECT().FetchByID(gomock.Any(), gomock.Any()).Return(course1, nil)
	orderGateway := NewMockorderGateway(ctrl)
	orderGateway.EXPECT().GetOrderByID(gomock.Any(), gomock.Any()).Return(order1, nil)
	appLogger, _ := logger.NewZapAdapter()
	paymentGateway := NewMockpaymentGateway(ctrl)
	paymentGateway.EXPECT().GetPaymentByOrderID(gomock.Any(), gomock.Any()).Return(payment1, nil)
	profileGateway := NewMockprofileGateway(ctrl)
	// profileGateway.EXPECT().GetProfileByUserID(gomock.Any(), gomock.Any()).Return(profile1, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)
	id, err := usecase.Record(context.Background(), uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"))

	assert.NoError(t, err)
	assert.Equal(t, id, uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"))
}

func TestRecordFailedPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, model.ErrNoEnrollmentFound)
	// enrollRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"), nil)
	courseRepo := NewMockcourseRepository(ctrl)
	courseRepo.EXPECT().FetchByID(gomock.Any(), gomock.Any()).Return(course1, nil)
	orderGateway := NewMockorderGateway(ctrl)
	orderGateway.EXPECT().GetOrderByID(gomock.Any(), gomock.Any()).Return(order1, nil)
	appLogger, _ := logger.NewZapAdapter()
	paymentGateway := NewMockpaymentGateway(ctrl)
	paymentGateway.EXPECT().GetPaymentByOrderID(gomock.Any(), gomock.Any()).Return(payment2, nil)
	profileGateway := NewMockprofileGateway(ctrl)
	// profileGateway.EXPECT().GetProfileByUserID(gomock.Any(), gomock.Any()).Return(profile1, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)
	_, err := usecase.Record(context.Background(), uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"))

	assert.Error(t, err)
	assert.ErrorIs(t, err, record.ErrCourseNotBeenPaid)
}

func TestRecordBaseStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, model.ErrNoEnrollmentFound)
	enrollRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"), nil)
	courseRepo := NewMockcourseRepository(ctrl)
	courseRepo.EXPECT().FetchByID(gomock.Any(), gomock.Any()).Return(course3, nil)
	orderGateway := NewMockorderGateway(ctrl)
	orderGateway.EXPECT().GetOrderByID(gomock.Any(), gomock.Any()).Return(order1, nil)
	appLogger, _ := logger.NewZapAdapter()
	paymentGateway := NewMockpaymentGateway(ctrl)
	// paymentGateway.EXPECT().GetPaymentByOrderID(gomock.Any(), gomock.Any()).Return(payment2, nil)
	profileGateway := NewMockprofileGateway(ctrl)
	// profileGateway.EXPECT().GetProfileByUserID(gomock.Any(), gomock.Any()).Return(profile1, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)
	id, err := usecase.Record(context.Background(), uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"))

	assert.NoError(t, err)
	assert.Equal(t, id, uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"))
}

func TestRecordFreeStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, model.ErrNoEnrollmentFound)
	enrollRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"), nil)
	courseRepo := NewMockcourseRepository(ctrl)
	courseRepo.EXPECT().FetchByID(gomock.Any(), gomock.Any()).Return(course4, nil)
	orderGateway := NewMockorderGateway(ctrl)
	orderGateway.EXPECT().GetOrderByID(gomock.Any(), gomock.Any()).Return(order1, nil)
	appLogger, _ := logger.NewZapAdapter()
	paymentGateway := NewMockpaymentGateway(ctrl)
	// paymentGateway.EXPECT().GetPaymentByOrderID(gomock.Any(), gomock.Any()).Return(payment2, nil)
	profileGateway := NewMockprofileGateway(ctrl)
	// profileGateway.EXPECT().GetProfileByUserID(gomock.Any(), gomock.Any()).Return(profile1, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)
	id, err := usecase.Record(context.Background(), uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"))

	assert.NoError(t, err)
	assert.Equal(t, id, uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"))
}

func TestRecord_OrderNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	courseRepo := NewMockcourseRepository(ctrl)
	orderGateway := NewMockorderGateway(ctrl)
	paymentGateway := NewMockpaymentGateway(ctrl)
	profileGateway := NewMockprofileGateway(ctrl)
	appLogger, _ := logger.NewZapAdapter()

	orderGateway.EXPECT().GetOrderByID(gomock.Any(), ID1).Return(nil, fmt.Errorf("order not found"))

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)

	_, err := usecase.Record(context.Background(), ID1)
	assert.Error(t, err)
}

func TestRecord_InvalidCourseID(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	courseRepo := NewMockcourseRepository(ctrl)
	orderGateway := NewMockorderGateway(ctrl)
	paymentGateway := NewMockpaymentGateway(ctrl)
	profileGateway := NewMockprofileGateway(ctrl)
	appLogger, _ := logger.NewZapAdapter()

	orderGateway.EXPECT().GetOrderByID(gomock.Any(), ID1).Return(&orders.Order{
		Id:       "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
		UserId:   "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
		CourseId: "invalid-uuid", // ❌ невалидный UUID
	}, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)

	_, err := usecase.Record(context.Background(), ID1)
	assert.Error(t, err)
}

func TestRecord_CourseNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	courseRepo := NewMockcourseRepository(ctrl)
	orderGateway := NewMockorderGateway(ctrl)
	paymentGateway := NewMockpaymentGateway(ctrl)
	profileGateway := NewMockprofileGateway(ctrl)
	appLogger, _ := logger.NewZapAdapter()

	orderGateway.EXPECT().GetOrderByID(gomock.Any(), ID1).Return(order1, nil)
	courseRepo.EXPECT().FetchByID(gomock.Any(), ID1).Return(nil, model.ErrNoCourseFound)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)

	_, err := usecase.Record(context.Background(), ID1)
	assert.Error(t, err)
}

func TestRecord_InvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	courseRepo := NewMockcourseRepository(ctrl)
	orderGateway := NewMockorderGateway(ctrl)
	paymentGateway := NewMockpaymentGateway(ctrl)
	profileGateway := NewMockprofileGateway(ctrl)
	appLogger, _ := logger.NewZapAdapter()

	orderGateway.EXPECT().GetOrderByID(gomock.Any(), ID1).Return(&orders.Order{
		Id:       "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
		UserId:   "invalid-uuid", // ❌ невалидный UUID
		CourseId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b",
	}, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)

	_, err := usecase.Record(context.Background(), ID1)
	assert.Error(t, err)
}

func TestRecord_EnrollmentAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	courseRepo := NewMockcourseRepository(ctrl)
	orderGateway := NewMockorderGateway(ctrl)
	paymentGateway := NewMockpaymentGateway(ctrl)
	profileGateway := NewMockprofileGateway(ctrl)
	appLogger, _ := logger.NewZapAdapter()

	// ✅ Validator находит существующий энроллмент
	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), ID1, ID1).Return(&model.Enrollment{}, nil)
	// ❌ Create не должен вызываться

	orderGateway.EXPECT().GetOrderByID(gomock.Any(), ID1).Return(order1, nil)
	courseRepo.EXPECT().FetchByID(gomock.Any(), ID1).Return(course1, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)

	_, err := usecase.Record(context.Background(), ID1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, record.ErrSecondRecord)
}

func TestRecord_CreateEnrollmentFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	courseRepo := NewMockcourseRepository(ctrl)
	orderGateway := NewMockorderGateway(ctrl)
	paymentGateway := NewMockpaymentGateway(ctrl)
	profileGateway := NewMockprofileGateway(ctrl)
	appLogger, _ := logger.NewZapAdapter()

	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, model.ErrNoEnrollmentFound)
	enrollRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.Nil, fmt.Errorf("db error"))

	orderGateway.EXPECT().GetOrderByID(gomock.Any(), ID1).Return(order1, nil)
	courseRepo.EXPECT().FetchByID(gomock.Any(), ID1).Return(course1, nil)
	paymentGateway.EXPECT().GetPaymentByOrderID(gomock.Any(), ID1).Return(payment1, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)

	_, err := usecase.Record(context.Background(), ID1)
	assert.Error(t, err)
}

func TestRecord_PremiumStrategy_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	courseRepo := NewMockcourseRepository(ctrl)
	orderGateway := NewMockorderGateway(ctrl)
	paymentGateway := NewMockpaymentGateway(ctrl)
	profileGateway := NewMockprofileGateway(ctrl)
	appLogger, _ := logger.NewZapAdapter()

	coursePremium := &model.Course{ID: ID1, Price: 100, Type: model.CoursePremium}

	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), ID1, ID1).Return(nil, model.ErrNoEnrollmentFound)
	enrollRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(ID1, nil)

	orderGateway.EXPECT().GetOrderByID(gomock.Any(), ID1).Return(order1, nil)
	courseRepo.EXPECT().FetchByID(gomock.Any(), ID1).Return(coursePremium, nil)
	paymentGateway.EXPECT().GetPaymentByOrderID(gomock.Any(), ID1).Return(payment1, nil)
	profileGateway.EXPECT().GetProfileByUserID(gomock.Any(), ID1).Return(profile1, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)

	id, err := usecase.Record(context.Background(), ID1)
	assert.NoError(t, err)
	assert.Equal(t, ID1, id)
}

func TestRecord_PremiumStrategy_WrongTier(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	courseRepo := NewMockcourseRepository(ctrl)
	orderGateway := NewMockorderGateway(ctrl)
	paymentGateway := NewMockpaymentGateway(ctrl)
	profileGateway := NewMockprofileGateway(ctrl)
	appLogger, _ := logger.NewZapAdapter()

	coursePremium := &model.Course{ID: ID1, Price: 100, Type: model.CoursePremium}
	profileBasic := &profiles.Profile{UserId: "c8bfd113-1023-49eb-a077-6820dd7e7c9b", Tier: "basic"}

	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), ID1, ID1).Return(nil, model.ErrNoEnrollmentFound)

	orderGateway.EXPECT().GetOrderByID(gomock.Any(), ID1).Return(order1, nil)
	courseRepo.EXPECT().FetchByID(gomock.Any(), ID1).Return(coursePremium, nil)
	paymentGateway.EXPECT().GetPaymentByOrderID(gomock.Any(), ID1).Return(payment1, nil)
	profileGateway.EXPECT().GetProfileByUserID(gomock.Any(), ID1).Return(profileBasic, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)

	_, err := usecase.Record(context.Background(), ID1)
	assert.Error(t, err)
	assert.ErrorIs(t, err, record.ErrTier)
}

func TestRecord_PremiumStrategy_ProfileGatewayError(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	courseRepo := NewMockcourseRepository(ctrl)
	orderGateway := NewMockorderGateway(ctrl)
	paymentGateway := NewMockpaymentGateway(ctrl)
	profileGateway := NewMockprofileGateway(ctrl)
	appLogger, _ := logger.NewZapAdapter()

	coursePremium := &model.Course{ID: ID1, Price: 100, Type: model.CoursePremium}

	// Моки для успешного прохождения валидации и оплаты
	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), ID1, ID1).Return(nil, model.ErrNoEnrollmentFound)
	orderGateway.EXPECT().GetOrderByID(gomock.Any(), ID1).Return(order1, nil)
	courseRepo.EXPECT().FetchByID(gomock.Any(), ID1).Return(coursePremium, nil)
	paymentGateway.EXPECT().GetPaymentByOrderID(gomock.Any(), ID1).Return(payment1, nil)

	profileGateway.EXPECT().GetProfileByUserID(gomock.Any(), ID1).Return(nil, fmt.Errorf("profile service unavailable"))

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)

	_, err := usecase.Record(context.Background(), ID1)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, record.ErrTier) // ошибка не из-за tier, а из-за гейтвея
}

func TestNOPPaymentChecker_CheckPayment(t *testing.T) {
	checker := &record.NOPPaymentChecker{}
	err := checker.CheckPayment(context.Background(), uuid.New())
	assert.NoError(t, err)
}

func TestCheckerFactory_GetPaymentChecker_ZeroPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	paymentGateway := NewMockpaymentGateway(ctrl)
	factory := record.NewCheckerFactory(paymentGateway)

	// Цена 0 → должен вернуться NOPPaymentChecker
	checker := factory.GetPaymentChecker(0)

	// Проверяем, что это действительно NOP (можно через type assertion или просто вызов)
	_, ok := checker.(*record.NOPPaymentChecker)
	assert.True(t, ok, "Expected NOPPaymentChecker for zero price")

	// И что он работает
	err := checker.CheckPayment(context.Background(), uuid.New())
	assert.NoError(t, err)
}

func TestValidator_Validate_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	enrollRepo := NewMockenrollmentRepository(ctrl)
	enrollRepo.EXPECT().FetchByCourseIDUserID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error database"))
	// enrollRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"), nil)
	courseRepo := NewMockcourseRepository(ctrl)
	courseRepo.EXPECT().FetchByID(gomock.Any(), gomock.Any()).Return(course1, nil)
	orderGateway := NewMockorderGateway(ctrl)
	orderGateway.EXPECT().GetOrderByID(gomock.Any(), gomock.Any()).Return(order1, nil)
	appLogger, _ := logger.NewZapAdapter()
	paymentGateway := NewMockpaymentGateway(ctrl)
	// paymentGateway.EXPECT().GetPaymentByOrderID(gomock.Any(), gomock.Any()).Return(payment1, nil)
	profileGateway := NewMockprofileGateway(ctrl)
	// profileGateway.EXPECT().GetProfileByUserID(gomock.Any(), gomock.Any()).Return(profile1, nil)

	factory := record.NewFactory(enrollRepo, paymentGateway, orderGateway, profileGateway)
	usecase := record.NewUsecase(enrollRepo, courseRepo, factory, orderGateway, appLogger)
	_, err := usecase.Record(context.Background(), uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"))
	t.Log(err)
	assert.Error(t, err)
	assert.EqualError(t, err, "error database")
}
