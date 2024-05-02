package application

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/pkg/logger"
	"github.com/sourcegraph/conc"
)

type CreateChannelParam struct {
	GuildID       string
	SpreadSheetID string
	Range         string
}

func (as *ApplicationService) CraeteChannel(ctx context.Context, param CreateChannelParam) ([]string, error) {
	var (
		wg      conc.WaitGroup
		message []string
	)
	defer wg.Wait()

	values, err := as.gs.Read(param.SpreadSheetID, param.Range)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read sheet")
	}

	for _, value := range values {
		teamName := value[0]

		if strings.Contains(teamName, "欠席") {
			continue
		}

		logger.Debug(ctx, "creating channel ", logger.Field("role", teamName))

		// 雑談 , 会議 , vcのチャンネルを作る
		wg.Go(func() {
			logger.Debug(ctx, "creating category", logger.Field("category", teamName))
			categoryID, err := as.ds.CreateChannelCategory(ctx, param.GuildID, teamName)
			if err != nil {
				logger.Error(ctx, "failed to create Category", logger.Field("err", err))
				message = append(message, err.Error())
			}

			logger.Debug(ctx, "Create category successful", logger.Field("category", teamName))
			logger.Debug(ctx, "creating channel", logger.Field("channel", "雑談"))

			if _, err := as.ds.CreateChannelText(ctx, param.GuildID, categoryID, "雑談"); err != nil {
				logger.Error(ctx, "failed to create channel", logger.Field("err", err))
				message = append(message, err.Error())
			}

			logger.Debug(ctx, "Create channel successful", logger.Field("channel", "雑談"))
			logger.Debug(ctx, "creating channel", logger.Field("channel", "会議"))

			if _, err := as.ds.CreateChannelText(ctx, param.GuildID, categoryID, "会議"); err != nil {
				logger.Error(ctx, "failed to create channel", logger.Field("err", err))
				message = append(message, err.Error())
			}

			logger.Debug(ctx, "Create channel successful", logger.Field("channel", "会議"))
			logger.Debug(ctx, "creating channel", logger.Field("channel", "vc"))

			if _, err := as.ds.CreateChannelVoice(ctx, param.GuildID, categoryID, "vc"); err != nil {
				logger.Error(ctx, "failed to create channel", logger.Field("err", err))
				message = append(message, err.Error())
			}

			logger.Debug(ctx, "Create channel successful", logger.Field("channel", "vc"))
		})
	}

	return message, nil
}

type DeleteChannelParam struct {
	GuildID       string
	SpreadSheetID string
	Range         string
}

func (as *ApplicationService) DeleteChannel(ctx context.Context, param DeleteChannelParam) ([]string, error) {
	values, err := as.gs.Read(config.Config.Spreadsheets.ID, config.Config.Spreadsheets.Range)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read sheet")
	}

	channels, err := as.ds.GetChannel(ctx, param.GuildID)
	if err != nil {
		logger.Error(ctx, "failed to get channels", logger.Field("err", err))
		return nil, err
	}

	roles := make(map[string]struct{})
	for _, value := range values {
		teamName := value[0]
		if teamName == "欠席" {
			continue
		}

		roles[strings.TrimSpace(teamName)] = struct{}{}
	}

	var deleteTargetChannelIDs []string
	currentCategory := make(map[string]string)

	for _, channel := range channels {
		if channel.Type == 4 {
			// rolesに存在しない場合無視
			if _, ok := roles[channel.Name]; !ok {
				continue
			}

			currentCategory[channel.ID] = channel.Name
			deleteTargetChannelIDs = append(deleteTargetChannelIDs, channel.ID)
			continue
		}

		_, ok := currentCategory[channel.ParentID]
		if !ok {
			continue
		}

		deleteTargetChannelIDs = append(deleteTargetChannelIDs, channel.ID)
	}

	var (
		wg      conc.WaitGroup
		message []string
	)
	defer wg.Wait()

	for _, deleteTargetChannelID := range deleteTargetChannelIDs {
		wg.Go(func() {
			time.Sleep(time.Second / 10)
			logger.Debug(ctx, "deleting channel ", logger.Field("channelId", deleteTargetChannelID))

			if err := as.ds.DeleteChannel(ctx, deleteTargetChannelID); err != nil {
				logger.Error(ctx, "failed to delete channel", logger.Field("err", err))
				message = append(message, err.Error())
			}

			logger.Debug(ctx, "channel delete successful", logger.Field("channelId", deleteTargetChannelID))
		})
	}
	return message, nil
}
