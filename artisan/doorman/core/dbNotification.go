/*
  Onix Config Manager - Artisan's Doorman
  Copyright (c) 2018-Present by www.gatblau.org
  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
  Contributors to this project, hereby assign copyright in this code to the project,
  to be licensed under the same terms as the rest of the code.
*/

package core

import (
	"context"
	"fmt"
	"github.com/gatblau/onix/artisan/doorman/types"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func UpsertNotification(notification types.Notification) (error, int) {
	db := NewDb()
	_, err := db.FindByName(types.NotificationTemplatesCollection, notification.Template)
	if err != nil {
		return fmt.Errorf("cannot find notification template %s for notification %s: %s", notification.Template, notification.Name, err), http.StatusBadRequest
	}
	var resultCode int
	_, err, resultCode = db.UpsertObject(types.NotificationsCollection, notification)
	if err != nil {
		return fmt.Errorf("cannot update notification in database"), resultCode
	}
	return nil, resultCode
}

func FindAllNotifications() ([]types.Notification, error) {
	db := NewDb()
	var notifications []types.Notification
	if err := db.FindMany(types.NotificationsCollection, nil, func(c *mongo.Cursor) error {
		return c.All(context.Background(), &notifications)
	}); err != nil {
		return nil, err
	}
	return notifications, nil
}

func FindAllNotificationTemplates() ([]types.NotificationTemplate, error) {
	db := NewDb()
	var notificationTemplates []types.NotificationTemplate
	if err := db.FindMany(types.NotificationTemplatesCollection, nil, func(c *mongo.Cursor) error {
		return c.All(context.Background(), &notificationTemplates)
	}); err != nil {
		return nil, err
	}
	return notificationTemplates, nil
}
