// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	// "time"
	// "strconv"
	"log"
)

func (a *App) SendCustomAutoResponseToUsers(channel *model.Channel, sender *model.User) (bool, *model.AppError) {

    message := sender.NotifyProps[model.AUTO_RESPONDER_MESSAGE_NOTIFY_PROP]

    if message == "" {
        return false, nil
    }
    log.Print("part 3 under for")

    autoResponderPost := &model.Post{
        ChannelId: channel.Id,
        Message:   message,
        RootId:    "",
        ParentId:  "",
        Type:      model.POST_AUTO_RESPONDER,
        UserId:    sender.Id,
    }

    log.Print("part 4 under for")

    if _, err := a.CreatePost(autoResponderPost, channel, false); err != nil {
        mlog.Error(err.Error())
        return false, err
    }

    log.Print("part 5 under for")

    log.Print("comes under if nec")

	return true, nil
}



func (a *App) SendAutoResponseIfNecessary(channel *model.Channel, sender *model.User) (bool, *model.AppError) {
	if channel.Type != model.CHANNEL_DIRECT {
		return false, nil
	}
    log.Print("comes under if nec")
	receiverId := channel.GetOtherUserIdForDM(sender.Id)

	receiver, err := a.GetUser(receiverId)
	if err != nil {
		return false, err
	}

	return a.SendAutoResponse(channel, receiver)
}

func (a *App) SendAutoResponse(channel *model.Channel, receiver *model.User) (bool, *model.AppError) {
	if receiver == nil || receiver.NotifyProps == nil {
		return false, nil
	}

    message := receiver.NotifyProps[model.AUTO_RESPONDER_MESSAGE_NOTIFY_PROP]
    // needed if auto responder is dependent on active status
	// active := receiver.NotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] == "true"

	// needed if auto responder duration is calculated here
	// duration := receiver.NotifyProps[model.AUTO_RESPONDER_DURATION_NOTIFY_PROP]

	// if !active || message == "" {
	//	return false, nil
	// }

	if message == "" {
		return false, nil
	}

    // needed if auto responder duration is calculated here
    // replyPeriod, err := strconv.Atoi(duration)
    // if err == nil {
    //        time.Sleep(time.Duration(replyPeriod) * time.Second)
    // }


    // log.Print("comes before autoResponderPost")

	autoResponderPost := &model.Post{
		ChannelId: channel.Id,
		Message:   message,
		RootId:    "",
		ParentId:  "",
		Type:      model.POST_AUTO_RESPONDER,
		UserId:    receiver.Id,
	}

	if _, err := a.CreatePost(autoResponderPost, channel, false); err != nil {
	    log.Print("comes under last err")
		mlog.Error(err.Error())
		return false, err
	}

	return true, nil
}

func (a *App) SetAutoResponderStatus(user *model.User, oldNotifyProps model.StringMap) {
	active := user.NotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] == "true"
	oldActive := oldNotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] == "true"

	autoResponderDisabled := oldActive && !active

	if autoResponderDisabled {
		a.SetStatusOnline(user.Id, true)
	}
}

func (a *App) DisableAutoResponder(userId string, asAdmin bool) *model.AppError {
	user, err := a.GetUser(userId)
	if err != nil {
		return err
	}

	active := user.NotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] == "true"

	if active {
		patch := &model.UserPatch{}
		patch.NotifyProps = user.NotifyProps
		patch.NotifyProps[model.AUTO_RESPONDER_ACTIVE_NOTIFY_PROP] = "false"

		_, err := a.PatchUser(userId, patch, asAdmin)
		if err != nil {
			return err
		}
	}

	return nil
}
