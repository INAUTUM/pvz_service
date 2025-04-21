package service

import (
	"context"
	"testing"

	"pvz_service/internal/domain"
	"pvz_service/internal/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPVZService_CreatePVZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPVZRepository(ctrl)
	svc := NewPVZService(mockRepo)

	// Генерируем валидный UUID
	testUUID := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

	tests := []struct {
		name    string
		input   domain.PVZ
		mock    func()
		want    *domain.PVZ
		wantErr bool
	}{
		{
			name:  "valid creation",
			input: domain.PVZ{City: "Москва"},
			mock: func() {
				mockRepo.EXPECT().Create(
					gomock.Any(), 
					domain.PVZ{City: "Москва"},
				).Return(&domain.PVZ{
					ID:   testUUID,
					City: "Москва",
				}, nil)
			},
			want: &domain.PVZ{
				ID:   testUUID,
				City: "Москва",
			},
		},
		{
			name:  "invalid city",
			input: domain.PVZ{City: "Новосибирск"},
			mock:  func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := svc.CreatePVZ(context.Background(), tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}