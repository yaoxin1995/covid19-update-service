package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pmoule/go2hal/hal"
	"gopkg.in/guregu/null.v3"
)

type SubscriptionCollection []Subscription

type Subscription struct {
	CommonModelFields
	Email          null.String     `json:"email"`
	TelegramChatID null.String     `json:"telegramChatId"`
	Topics         TopicCollection `json:"-"`
}

func NewSubscription(email, telegram *string) (Subscription, error) {
	s := Subscription{
		Email:          null.StringFromPtr(email),
		TelegramChatID: null.StringFromPtr(telegram),
		Topics:         TopicCollection{},
	}
	err := s.Store()
	return s, err
}

func (s *Subscription) Store() error {
	return db.Save(&s).Error
}

func GetSubscription(id uint) (*Subscription, error) {
	s := &Subscription{}
	err := db.Preload("Topics").First(s, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func SubscriptionExists(id uint) bool {
	s := &Subscription{}
	err := db.First(s, id).Error
	if err != nil {
		return false
	}
	return true
}

func GetSubscriptions() (SubscriptionCollection, error) {
	var subs SubscriptionCollection
	err := db.Find(&subs).Preload("Topics").Error
	return subs, err
}

func (s *Subscription) Delete() error {
	return db.Unscoped().Delete(s).Error
}

func (s *Subscription) Update(email, telegram *string) error {
	s.Email = null.StringFromPtr(email)
	s.TelegramChatID = null.StringFromPtr(telegram)
	return s.Store()
}

func (s Subscription) ToHAL() hal.Resource {
	href := fmt.Sprintf("/subscriptions/%d", s.ID)

	root := hal.NewResourceObject()

	// Add email and telegram manually, because HAL JSON encoder does not can handle null.String
	data := root.Data()
	data["email"] = s.Email.Ptr()
	data["telegramChatId"] = s.TelegramChatID.Ptr()
	root.AddData(data)

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: href}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	topicsRel, _ := hal.NewLinkRelation("topics")
	topicsLink := &hal.LinkObject{Href: fmt.Sprintf("%s/topics", href)}
	topicsRel.SetLink(topicsLink)
	root.AddLink(topicsRel)

	return root
}

func (sc SubscriptionCollection) ToHAL() hal.Resource {
	href := fmt.Sprintf("/subscriptions")

	root := hal.NewResourceObject()

	selfRel := hal.NewSelfLinkRelation()
	selfLink := &hal.LinkObject{Href: href}
	selfRel.SetLink(selfLink)
	root.AddLink(selfRel)

	var embeddedSubs []hal.Resource

	for _, s := range sc {
		eHref := fmt.Sprintf("/subscriptions/%d", s.ID)
		eSelfLink, _ := hal.NewLinkObject(eHref)

		eSelfRel, _ := hal.NewLinkRelation("self")
		eSelfRel.SetLink(eSelfLink)

		embeddedSub := hal.NewResourceObject()
		embeddedSub.AddLink(eSelfRel)
		data := root.Data()
		data["email"] = s.Email.Ptr()
		data["telegramChatId"] = s.TelegramChatID.Ptr()
		root.AddData(data)
		embeddedSubs = append(embeddedSubs, embeddedSub)
	}

	subs, _ := hal.NewResourceRelation("subscriptions")
	subs.SetResources(embeddedSubs)
	root.AddResource(subs)

	return root
}
