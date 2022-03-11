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

type NotificationType string

const (
	SuccessNotification   NotificationType = "SUCCESS"
	CmdFailedNotification NotificationType = "CMD_FAILED"
	ErrorNotification     NotificationType = "ERROR"
)

func (db *Db) FindNotification(name string) (*types.PipeNotification, error) {
	notif := new(types.Notification)
	result, err := db.FindByName(types.NotificationsCollection, name)
	if err != nil {
		return nil, err
	}
	err = result.Decode(notif)
	if err != nil {
		return nil, err
	}
	// load the template
	templ := new(types.NotificationTemplate)
	result, err = db.FindByName(types.NotificationTemplatesCollection, notif.Template)
	if err != nil {
		return nil, err
	}
	err = result.Decode(templ)
	if err != nil {
		return nil, err
	}
	return &types.PipeNotification{
		Name:      notif.Name,
		Recipient: notif.Recipient,
		Type:      notif.Type,
		Subject:   templ.Subject,
		Content:   templ.Content,
	}, nil
}

func (db *Db) UpsertNotification(notification types.Notification) (error, int) {
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

func (db *Db) FindAllNotifications() ([]types.Notification, error) {
	var notifications []types.Notification
	if err := db.FindMany(types.NotificationsCollection, nil, func(c *mongo.Cursor) error {
		return c.All(context.Background(), &notifications)
	}); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (db *Db) FindAllNotificationTemplates() ([]types.NotificationTemplate, error) {
	var notificationTemplates []types.NotificationTemplate
	if err := db.FindMany(types.NotificationTemplatesCollection, nil, func(c *mongo.Cursor) error {
		return c.All(context.Background(), &notificationTemplates)
	}); err != nil {
		return nil, err
	}
	return notificationTemplates, nil
}
