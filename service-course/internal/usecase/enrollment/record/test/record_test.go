package test

import (
	"context"
	"log"
	"testing"

	"github.com/barashF/lms/service-course/internal/logger"
	"github.com/barashF/lms/service-course/internal/usecase/enrollment/record"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRecord(t *testing.T) {
	enrollmentRepo := MockEnrollmentRepository{}
	courseRepo := MockCourseRepository{}
	paymentGateway := MockPaymentGateway{}
	orderGateway := MockOrderGateway{}
	profileGateway := MockProfileGateway{}
	appLogger, err := logger.NewZapAdapter()
	if err != nil {
		log.Fatalf("failed initialize logger: %v", err)
	}
	defer appLogger.Sync()

	factory := record.NewFactory(&enrollmentRepo, &paymentGateway, &orderGateway, &profileGateway)
	usecase := record.NewUsecase(&enrollmentRepo, &courseRepo, factory, &orderGateway, appLogger)
	id, err := usecase.Record(context.Background(), uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"))
	assert.NoError(t, err)
	t.Log(id)
	assert.Equal(t, uuid.MustParse("c8bfd113-1023-49eb-a077-6820dd7e7c9b"), id)
}
