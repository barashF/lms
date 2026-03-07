package payment

type YooKassaGateway struct {
	ShopID    int64
	SecretKey string
}

// func (y *YooKassaGateway) CreatePayment(ctx context.Context, request *RequestCreate) (*model.Payment, error) {
// 	requestBody := map[string]any{
// 		"amount": map[string]any{
// 			"value": request.value,
// 			"currency": request.currency,
// 		},
// 		"confirmation":map[string]any{
// 			"type": "redirect",
// 		},
// 	}
// }

// curl https://api.yookassa.ru/v3/payments \
//   -X POST \
//   -u <Идентификатор магазина>:<Секретный ключ> \
//   -H 'Idempotence-Key: <Ключ идемпотентности>' \
//   -H 'Content-Type: application/json' \
//   -d '{
//         "amount": {
//           "value": "2.00",
//           "currency": "RUB"
//         },
//         "confirmation": {
//           "type": "embedded"
//         },
//         "capture": true,
//         "description": "Заказ №72"
//       }'
