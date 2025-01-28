package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go_parser/models"
	"go_parser/services"
	"go_parser/utils"
)

func ProcessMessage(body []byte, collection *mongo.Collection) error {
	var message struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &message); err != nil {
		utils.Logger.Printf("Ошибка разбора сообщения: %v\n", err)
		return err
	}
	utils.Logger.Printf("Обработка URL: %s\n", message.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Проверка наличия записи в MongoDB
	var record models.MongoRecord
	err := collection.FindOne(ctx, bson.M{"url": message.URL}).Decode(&record)
	if err == nil {
		// Запись найдена, проверяем метку ParsedAt
		if !record.ParsedAt.IsZero() {
			utils.Logger.Printf("Задание для URL '%s' уже выполнено.\n", message.URL)
			return nil // Задание уже выполнено, пропускаем
		}
	} else if err != mongo.ErrNoDocuments {
		// Ошибка при поиске записи
		utils.Logger.Printf("Ошибка поиска записи в MongoDB: %v\n", err)
		return err
	}

	// Получение HTML-страницы
	utils.Logger.Println("Получение HTML-страницы...")
	content, err := services.FetchPageWithPlaywright(message.URL)
	if err != nil {
		utils.Logger.Printf("Ошибка получения HTML-страницы: %v\n", err)
		return err
	}
	utils.Logger.Println("HTML-страница успешно получена.")

	// Обновление или вставка записи
	utils.Logger.Println("Обновление записи в MongoDB...")
	record = models.MongoRecord{
		URL:       message.URL,
		Content:   content,
		CreatedAt: time.Now(),
		ParsedAt:  time.Now(),
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"url": record.URL},
		bson.M{
			"$set": bson.M{
				"content":   record.Content,
				"parsed_at": record.ParsedAt,
			},
			"$setOnInsert": bson.M{"created_at": record.CreatedAt},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		utils.Logger.Printf("Ошибка обновления записи: %v\n", err)
		return err
	}

	utils.Logger.Printf("Запись для URL '%s' успешно обновлена.\n", record.URL)
	return nil
}
